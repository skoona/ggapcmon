package providers

import (
	"context"
	"encoding/binary"
	"github.com/skoona/ggapcmon/internal/entities"
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
	host          entities.ApcHost
	periodTicker  *time.Ticker
	activeSession net.Conn
	events        []string
	status        []string
	rcvr          chan []string
	log           *log.Logger
}

var (
	_ interfaces.ApcProvider = (*apcProvider)(nil)
	_ interfaces.Provider    = (*apcProvider)(nil)
)

func NewAPCProvider(ctx context.Context, host entities.ApcHost, receiver chan []string, log *log.Logger) (interfaces.ApcProvider, error) {
	provider := &apcProvider{
		ctx:    ctx,
		host:   host,
		status: []string{},
		events: []string{},
		rcvr:   receiver,
		log:    log,
	}
	err := provider.begin()
	if err != nil {
		return nil, err
	} else {
		return provider, nil
	}
}

// connect dials th apc server and establishes a connection
func (a *apcProvider) connect() error {

	if a.activeSession != nil {
		return nil
	}
	if a.ctx.Err() != nil {
		return a.ctx.Err()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", a.host.IpAddress)
	if err != nil {
		a.activeSession = nil
		a.log.Println("connect() dial Error: ", err.Error(), ", host: ", a.host.Name, ", context: ", ctx.Err())
	} else {
		a.activeSession = conn
	}
	return err
}

// begin connects to apc server and starts periodic data collection
func (a *apcProvider) begin() error {

	err := a.connect()
	if err != nil {
		a.log.Println("begin() connect Error: ", err.Error(), ", host: ", a.host.Name)
	} else {
		a.periodicUpdateStart()
	}
	return err
}

// periodicUpdateStart creates or reset the ticker and go routine
// which issues apc requests according to ticker's period
func (a *apcProvider) periodicUpdateStart() {
	if a.periodTicker != nil {
		a.periodTicker.Reset(a.host.SecondsPerSample * time.Second)
		return
	}
	a.periodTicker = time.NewTicker(a.host.SecondsPerSample * time.Second)

	go func(s *apcProvider) {
	back:
		for {
			select {
			case <-s.ctx.Done():
				a.log.Println("periodicUpdateStart(", s.Name(), ") ending: ", s.ctx.Err().Error())
				break back

			case <-s.periodTicker.C:
				_ = s.request(commandStatus, a.rcvr)
				_ = s.request(commandEvents, a.rcvr)
			}
		}
		a.log.Println("periodicUpdateStart(", a.Name(), ") ended ")
	}(a)
}

// periodicUpdateStop stops the ticker driving apc queries
func (a *apcProvider) periodicUpdateStop() {
	a.log.Println("periodicUpdateStop(", a.Name(), ") called.")
	a.periodTicker.Stop()
}

// End closes the apc connection and stops go routines
func (a *apcProvider) Shutdown() {
	a.log.Println("ApcProvider::Shutdown(", a.Name(), ") called.")
	if a.activeSession == nil {
		return
	}

	a.periodicUpdateStop()
	err := a.activeSession.Close()
	if err != nil {
		a.log.Println("ApcProvider::Shutdown()::Close(", a.Name(), ") Error: ", err.Error())
	}
	a.activeSession = nil
}

func (a *apcProvider) Name() string {
	return a.host.Name
}
func (a *apcProvider) IpAddress() string {
	return a.host.IpAddress
}

func (a *apcProvider) addEvent(newValue string) {
	a.events = append(a.events, newValue)
}
func (a *apcProvider) addStatus(newValue string) {
	a.status = append(a.status, newValue)
}
func (a *apcProvider) eventsSafeCopy() []string {
	return append([]string{}, a.events...)
}
func (a *apcProvider) statusSafeCopy() []string {
	return append([]string{}, a.status...)
}

func (a *apcProvider) sendCommand(command string) error {
	var msgLen = uint16(len(command))
	b := make([]byte, 2)

	binary.BigEndian.PutUint16(b, msgLen)
	_, err := a.activeSession.Write(b)
	if err != nil {
		a.log.Println("sendCommand() write len error: ", err.Error())
		return err
	}

	_, err = a.activeSession.Write([]byte(command))
	if err != nil {
		a.log.Println("sendCommand() write command error: ", err.Error())
		return err
	}

	return nil
}
func (a *apcProvider) receiveMessage() (string, error) {
	var msgLen uint16
	message := []byte{}
	b := make([]byte, 2)

	read, err := a.activeSession.Read(b)
	if err != nil {
		a.log.Println("receiveMessage() read len error: ", err.Error())
		return "", err
	}

	msgLen = binary.BigEndian.Uint16(b)
	if (read == 2 && msgLen == 0) || (read == 0 && msgLen == 0) || (msgLen > 1024) {
		return "", nil
	}

	line := make([]byte, msgLen)

	read, err = a.activeSession.Read(line)
	if err != nil {
		a.log.Println("receiveMessage() read message error: ", err.Error())
		return string(message), err
	}
	if read > 2 {
		message = append(message, line[0:read]...)
	}

	return string(message), err
}

// request gathers a list of responses for each command
// returns a string slice over the channel
func (a *apcProvider) request(command string, r chan []string) error {
	if a.activeSession == nil {
		err := a.connect()
		if err != nil {
			return err
		}
	}
	err := a.sendCommand(command)
	if err != nil {
		a.log.Println("request::sendCommand() send command error: ", err.Error())
		return err
	}
	if command == commandEvents {
		a.events = a.events[:0]
	} else {
		a.status = a.status[:0]
	}

transact:
	for err == nil {
		msg, err := a.receiveMessage()
		if err != nil {
			break transact
		}
		if len(msg) == 0 {
			break transact
		}
		if command == commandEvents {
			a.addEvent(command + ": " + msg)
		} else {
			a.addStatus(command + ": " + msg)
		}
	}
	a.activeSession.Close()
	a.activeSession = nil

	if command == commandEvents {
		r <- a.eventsSafeCopy()
	} else {
		r <- a.statusSafeCopy()
	}

	return err
}
