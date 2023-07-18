package services

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"time"
)

const (
	commandStatus = "status"
	commandEvents = "events"
)

type APCServer interface {
	Name() string
	SetName(newValue string)
	IpAddress() string
	SetIpAddress(newValue string)
	Begin() error
	AddEvent(newValue string)
	AddStatus(newValue string)
	Events() []string
	Status() []string
	PeriodicUpdateStart()
	PeriodicUpdateStop()
	Request(command string, r chan string) error
	SendCommand(command string) error
	ReceiveMessage() (string, error)
	End()
}

type apcServer struct {
	ctx           context.Context
	ipAddress     string
	name          string
	periodTicker  *time.Ticker
	samplePeriod  time.Duration
	activeSession net.Conn
	events        []string
	status        []string
	tickerChan    chan bool
	dialer        net.Dialer
	rcvr          chan string
}

var _ APCServer = (*apcServer)(nil)

func NewServer(ctx context.Context, name, ip string, secondsBetweenSamples time.Duration, receiver chan string) APCServer {
	return &apcServer{
		ctx:          ctx,
		ipAddress:    ip,
		name:         name,
		samplePeriod: secondsBetweenSamples,
		status:       []string{},
		events:       []string{},
		tickerChan:   make(chan bool),
		dialer:       net.Dialer{},
		rcvr:         receiver,
	}
}
func (a *apcServer) Begin() error {

	if a.activeSession != nil {
		return nil
	}
	if a.ctx.Err() != nil {
		return a.ctx.Err()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := a.dialer.DialContext(ctx, "tcp", a.ipAddress)
	if err != nil {
		log.Println("Begin() dial Error: ", err.Error(), ", host: ", a.name, ", context: ", ctx.Err())
	} else {
		a.activeSession = conn
		a.PeriodicUpdateStart()
	}
	return err
}
func (a *apcServer) End() {
	if a.activeSession == nil {
		return
	}

	a.PeriodicUpdateStop()
	err := a.activeSession.Close()
	if err != nil {
		log.Println("Close() Error: ", err.Error())
	}
	a.activeSession = nil

	return
}

func (a *apcServer) Name() string {
	return a.name
}
func (a *apcServer) SetName(newValue string) {
	a.name = newValue
}
func (a *apcServer) IpAddress() string {
	return a.ipAddress
}
func (a *apcServer) SetIpAddress(newValue string) {
	a.ipAddress = newValue
}
func (a *apcServer) AddEvent(newValue string) {
	a.events = append(a.events, newValue)
}
func (a *apcServer) AddStatus(newValue string) {
	a.status = append(a.status, newValue)
}
func (a *apcServer) Events() []string {
	return append([]string{}, a.events...)
}
func (a *apcServer) Status() []string {
	return append([]string{}, a.status...)
}
func (a *apcServer) PeriodicUpdateStart() {
	a.periodTicker = time.NewTicker(a.samplePeriod * time.Second)

	go func(rcvr chan string) {
	back:
		for {
			select {
			case <-a.ctx.Done():
				log.Println("PeriodicUpdateStart() ending: ", a.ctx.Err().Error())
				break back

			case <-a.tickerChan:
				break back

			case <-a.periodTicker.C:
				_ = a.Request(commandStatus, rcvr)
				_ = a.Request(commandEvents, rcvr)
			}
		}
		log.Println("PeriodicUpdateStart() ended ")
	}(a.rcvr)
}
func (a *apcServer) PeriodicUpdateStop() {
	a.periodTicker.Stop()
	a.tickerChan <- true
}

func (a *apcServer) SendCommand(command string) error {
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

func (a *apcServer) ReceiveMessage() (string, error) {
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

func (a *apcServer) Request(command string, r chan string) error {
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
			a.AddEvent(msg)
		} else {
			a.AddStatus(msg)
		}
		r <- command + ": " + msg
	}

	return err
}
