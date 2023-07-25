/*
ApcProvider
Connects to an APC Host and provides slices of the Events and Status topics
on a period cycle declared in the host mopdel as network sampling period.
Slices are returned over a seperate channels to whoever is listening.
*/
package providers

import (
	"context"
	"encoding/binary"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"net"
	"strings"
	"time"
)

const (
	commandStatus = "status"
	commandEvents = "events"
)

type apcProvider struct {
	ctx           context.Context
	host          *entities.ApcHost
	periodTicker  *time.Ticker
	activeSession net.Conn
	events        []string
	status        []string
	tuple         entities.ChannelTuple
}

var (
	_ interfaces.ApcProvider = (*apcProvider)(nil)
	_ interfaces.Provider    = (*apcProvider)(nil)
)

func NewAPCProvider(ctx context.Context, host *entities.ApcHost, tuple entities.ChannelTuple) (interfaces.ApcProvider, error) {
	provider := &apcProvider{
		ctx:    ctx,
		host:   host,
		status: []string{},
		events: []string{},
		tuple:  tuple,
	}
	err := provider.begin()
	if err != nil {
		return nil, err
	} else {
		return provider, nil
	}
}

// connect dials the apc server and establishes a connection
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
		a.host.State = commons.HostStatusUnknown
		commons.DebugLog("connect() dial Error: ", err.Error(), ", host: ", a.host.Name, ", context: ", ctx.Err())
	} else {
		a.activeSession = conn
		//a.host.State = commons.HostStatusOnline
	}
	return err
}

// begin connects to apc server and starts periodic data collection
func (a *apcProvider) begin() error {

	err := a.connect()
	if err != nil {
		commons.DebugLog("begin() connect Error: ", err.Error(), ", host: ", a.host.Name)
	} else {
		a.periodicUpdateStart()
	}
	return err
}

// periodicUpdateStart creates or reset the ticker and go routine
// which issues apc requests according to ticker's period
func (a *apcProvider) periodicUpdateStart() {
	if a.periodTicker != nil {
		a.periodTicker.Reset(a.host.NetworkSamplePeriod * time.Second)
		return
	}
	a.periodTicker = time.NewTicker(a.host.NetworkSamplePeriod * time.Second)

	go func(s *apcProvider) {
	back:
		for {
			select {
			case <-s.ctx.Done():
				commons.DebugLog("periodicUpdateStart(", s.host.Name, ") ending: ", s.ctx.Err().Error())
				break back

			case <-s.periodTicker.C:
				_ = s.request(commandStatus, a.tuple.Status)
				time.Sleep(16 * time.Millisecond)
				_ = s.request(commandEvents, a.tuple.Events)
			}
		}
		commons.DebugLog("periodicUpdateStart(", a.host.Name, ") ended ")
	}(a)
}

// periodicUpdateStop stops the ticker driving apc queries
func (a *apcProvider) periodicUpdateStop() {
	commons.DebugLog("periodicUpdateStop(", a.host.Name, ") called.")
	a.periodTicker.Stop()
}

// Shutdown closes the apc connection and stops go routines
func (a *apcProvider) Shutdown() {
	commons.DebugLog("ApcProvider::Shutdown(", a.host.Name, ") called.")
	if a.activeSession == nil {
		return
	}

	a.periodicUpdateStop()
	err := a.activeSession.Close()
	if err != nil {
		commons.DebugLog("ApcProvider::Shutdown()::Close(", a.host.Name, ") Error: ", err.Error())
	}
	a.activeSession = nil
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
		commons.DebugLog("ApcProvider::sendCommand(", a.host.Name, ") write len error: ", err.Error())
		return err
	}

	_, err = a.activeSession.Write([]byte(command))
	if err != nil {
		commons.DebugLog("ApcProvider::sendCommand(", a.host.Name, ") write command error: ", err.Error())
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
		commons.DebugLog("ApcProvider::receiveMessage(", a.host.Name, ") read len error: ", err.Error())
		return "", err
	}

	msgLen = binary.BigEndian.Uint16(b)
	if (read == 2 && msgLen == 0) || (read == 0 && msgLen == 0) || (msgLen > 1024) {
		return "", nil
	}

	line := make([]byte, msgLen)

	read, err = a.activeSession.Read(line)
	if err != nil {
		commons.DebugLog("ApcProvider::receiveMessage(", a.host.Name, ") read message error: ", err.Error())
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
		if len(msg) > 12 {
			if command == commandEvents {
				msg = commons.ChangeTimeFormat(msg[0:25], time.RFC1123) + msg[26:]
				msg = strings.TrimSpace(msg) + "\n"
				a.addEvent(msg)
			} else {
				trigger := strings.Count(msg, ":")
				if trigger >= 3 {
					msg = msg[0:11] + commons.ChangeTimeFormat(msg[11:], time.RFC1123)
				}
				msg = strings.TrimSpace(msg) + "\n"
				a.addStatus(msg)
			}
		}
		time.Sleep(64 * time.Millisecond) // let apcupsd breath a little
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
