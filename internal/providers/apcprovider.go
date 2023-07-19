package providers

import (
	"context"
	"encoding/binary"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"log"
	"net"
	"time"
)

const (
	commandStatus = "status"
	commandEvents = "events"
)

type apcProvider struct {
	ctx           context.Context
	ipAddress     string
	name          string
	periodTicker  *time.Ticker
	samplePeriod  time.Duration
	activeSession net.Conn
	events        []string
	status        []string
	rcvr          chan []string
}

var (
	_ interfaces.ApcProvider = (*apcProvider)(nil)
	_ interfaces.Provider    = (*apcProvider)(nil)
)

func NewAPCProvider(ctx context.Context, name, ip string,
	secondsBetweenSamples time.Duration, receiver chan []string) (interfaces.ApcProvider, error) {
	provider := &apcProvider{
		ctx:          ctx,
		ipAddress:    ip,
		name:         name,
		samplePeriod: secondsBetweenSamples,
		status:       []string{},
		events:       []string{},
		rcvr:         receiver,
	}
	err := provider.Begin()
	if err != nil {
		return nil, err
	} else {
		return provider, nil
	}
}

// Connect dials th apc server and establishes a connection
func (a *apcProvider) Connect() error {

	if a.activeSession != nil {
		return nil
	}
	if a.ctx.Err() != nil {
		return a.ctx.Err()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", a.ipAddress)
	if err != nil {
		a.activeSession = nil
		log.Println("Connect() dial Error: ", err.Error(), ", host: ", a.name, ", context: ", ctx.Err())
	} else {
		a.activeSession = conn
	}
	return err
}

// Begin connects to apc server and starts periodic data collection
func (a *apcProvider) Begin() error {

	err := a.Connect()
	if err != nil {
		log.Println("Begin() Connect Error: ", err.Error(), ", host: ", a.name)
	} else {
		a.PeriodicUpdateStart()
	}
	return err
}

// PeriodicUpdateStart creates or reset the ticker and go routine
// which issues apc requests according to ticker's period
func (a *apcProvider) PeriodicUpdateStart() {
	if a.periodTicker != nil {
		a.periodTicker.Reset(a.samplePeriod * time.Second)
		return
	}
	a.periodTicker = time.NewTicker(a.samplePeriod * time.Second)

	go func(rcvr chan []string) {
	back:
		for {
			select {
			case <-a.ctx.Done():
				log.Println("PeriodicUpdateStart(", a.Name(), ") ending: ", a.ctx.Err().Error())
				break back

			case <-a.periodTicker.C:
				_ = a.Request(commandStatus, rcvr)
				_ = a.Request(commandEvents, rcvr)
			}
		}
		log.Println("PeriodicUpdateStart(", a.Name(), ") ended ")
	}(a.rcvr)
}

// PeriodicUpdateStop stops the ticker driving apc queries
func (a *apcProvider) PeriodicUpdateStop() {
	log.Println("PeriodicUpdateStop(", a.Name(), ") called.")
	a.periodTicker.Stop()
}

// End closes the apc connection and stops go routines
func (a *apcProvider) Shutdown() {
	log.Println("Shutdown(", a.Name(), ") called.")
	if a.activeSession == nil {
		return
	}

	a.PeriodicUpdateStop()
	err := a.activeSession.Close()
	if err != nil {
		log.Println("Shutdown()::Close() Error: ", err.Error())
	}
	a.activeSession = nil
}

func (a *apcProvider) Name() string {
	return a.name
}
func (a *apcProvider) SetName(newValue string) {
	a.name = newValue
}
func (a *apcProvider) IpAddress() string {
	return a.ipAddress
}
func (a *apcProvider) SetIpAddress(newValue string) {
	a.ipAddress = newValue
}

func (a *apcProvider) AddEvent(newValue string) {
	a.events = append(a.events, newValue)
}
func (a *apcProvider) AddStatus(newValue string) {
	a.status = append(a.status, newValue)
}
func (a *apcProvider) Events() []string {
	return append([]string{}, a.events...)
}
func (a *apcProvider) Status() []string {
	return append([]string{}, a.status...)
}

func (a *apcProvider) SendCommand(command string) error {
	var msgLen = uint16(len(command))
	b := make([]byte, 2)

	binary.BigEndian.PutUint16(b, msgLen)
	_, err := a.activeSession.Write(b)
	if err != nil {
		log.Println("SendCommand() write len error: ", err.Error())
		return err
	}

	_, err = a.activeSession.Write([]byte(command))
	if err != nil {
		log.Println("SendCommand() write command error: ", err.Error())
		return err
	}

	return nil
}
func (a *apcProvider) ReceiveMessage() (string, error) {
	var msgLen uint16
	message := []byte{}
	b := make([]byte, 2)

	read, err := a.activeSession.Read(b)
	if err != nil {
		log.Println("ReceiveMessage() read len error: ", err.Error())
		return "", err
	}

	msgLen = binary.BigEndian.Uint16(b)
	if (read == 2 && msgLen == 0) || (read == 0 && msgLen == 0) || (msgLen > 1024) {
		return "", nil
	}

	line := make([]byte, msgLen)

	read, err = a.activeSession.Read(line)
	if err != nil {
		log.Println("ReceiveMessage() read message error: ", err.Error())
		return string(message), err
	}
	if read > 2 {
		message = append(message, line[0:read]...)
	}

	return string(message), err
}

// Request gathers a list of responses for each command
// returns a string slice over the channel
func (a *apcProvider) Request(command string, r chan []string) error {
	if a.activeSession == nil {
		err := a.Connect()
		if err != nil {
			return err
		}
	}
	err := a.SendCommand(command)
	if err != nil {
		log.Println("Request::SendCommand() send command error: ", err.Error())
		return err
	}
	if command == commandEvents {
		a.events = a.events[:0]
	} else {
		a.status = a.status[:0]
	}

transact:
	for err == nil {
		msg, err := a.ReceiveMessage()
		if err != nil {
			break transact
		}
		if len(msg) == 0 {
			break transact
		}
		if command == commandEvents {
			a.AddEvent(command + ": " + msg)
		} else {
			a.AddStatus(command + ": " + msg)
		}
	}
	a.activeSession.Close()
	a.activeSession = nil

	if command == commandEvents {
		r <- a.Events()
	} else {
		r <- a.Status()
	}

	return err
}
