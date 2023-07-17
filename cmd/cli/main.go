package main

import (
	"context"
	"encoding/binary"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort   = 3551
	HostVserv     = "10.100.1.3:3551"
	HostPve       = "10.100.1.4:3551"
	commandStatus = "status"
	commandEvents = "events"
)

func SendCommand(conn net.Conn, command string) error {
	startTime := time.Now()
	var msgLen = uint16(len(command))
	b := make([]byte, 2)

	binary.BigEndian.PutUint16(b, msgLen)
	_, err := conn.Write(b)
	if err != nil {
		log.Println("SendCommand(a) Duration mcs: ", time.Until(startTime).Microseconds())
		log.Println("Write len Error: ", err.Error())
		return err
	}

	_, err = conn.Write([]byte(command))
	if err != nil {
		log.Println("SendCommand(b) Duration mcs: ", time.Until(startTime).Microseconds())
		log.Println("Write command Error: ", err.Error())
		return err
	}

	log.Println("SendCommand(c) Duration mcs: ", time.Until(startTime).Microseconds())
	return nil
}

func ReceiveMessage(conn net.Conn) (string, error) {
	var msgLen uint16
	message := []byte{}
	b := make([]byte, 2)
	startTime := time.Now()

	read, err := conn.Read(b)

	msgLen = binary.BigEndian.Uint16(b)
	if (read == 2 && msgLen == 0) || (read == 0 && msgLen == 0) {
		log.Println("ReceiveMessage(a) Duration ms: ", time.Until(startTime).Milliseconds(), ", read: ", read, ", msgLen: ", msgLen)
		return string(message), nil

	} else if err != nil {
		log.Println("ReceiveMessage(b) Duration mcs: ", time.Until(startTime).Microseconds(), ", read: ", read, ", msgLen: ", msgLen)
		log.Println("Read len Error: ", err.Error())
		return "", err
	}

	line := make([]byte, msgLen)

	read, err = conn.Read(line)
	if err != nil {
		log.Println("ReceiveMessage(c) Duration mcs: ", time.Until(startTime).Microseconds())
		log.Println("Read message Error: ", err.Error())
		return string(message), err
	}
	if read > 2 {
		message = append(message, line...)
	}
	if msgLen != uint16(read) && read != 0 {
		secondLine := make([]byte, msgLen)
		for x := read; x <= int(msgLen); {
			read, err = conn.Read(secondLine)
			if read > 2 {
				message = append(message, secondLine[0:read]...)
			} else {
				break
			}
			if err != nil {
				log.Println("ReceiveMessage(x) Duration mcs: ", time.Until(startTime).Microseconds(), ", read: ", read, ", msgLen: ", msgLen)
				log.Println("Read x Error: ", err.Error())
				break
			}
			x = x + read
		}
	}

	log.Println("ReceiveMessage(d) Duration mcs: ", time.Until(startTime).Microseconds(), ", read: ", read, ", msgLen: ", msgLen)
	return string(message), err
}

func ApcRequest(conn net.Conn, command string) ([]string, error) {
	var data []string

	err := SendCommand(conn, command)

transact:
	for err == nil {
		msg, err := ReceiveMessage(conn)
		if err != nil {
			break transact
		}
		if len(msg) != 0 {
			data = append(data, msg)
		} else {
			break transact
		}
	}

	return data, err
}

func main() {
	var d net.Dialer

	systemSignalChannel := make(chan os.Signal, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	logger := log.New(os.Stdout, "[DEBUG] ", log.Lmicroseconds|log.Lshortfile)

	go func(stopFlag chan os.Signal) {
		signal.Notify(stopFlag, syscall.SIGINT, syscall.SIGTERM)
		sig := <-stopFlag // wait on ctrl-c
		logger.Println("Signal Received: ", sig.String())
		cancel()
	}(systemSignalChannel)

	startTime := time.Now()
	conn, err := d.DialContext(ctx, "tcp", HostPve)
	if err != nil {
		log.Panic("Dial Error: ", err.Error())
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Close() Error: ", err.Error())
		}
	}(conn)
	log.Println("Dial Duration ms: ", time.Until(startTime).Microseconds())

	ApcInfo := map[string][]string{}

	result, err := ApcRequest(conn, commandStatus)
	if err != nil {

	} else {
		ApcInfo[commandStatus] = result
	}

	result, err = ApcRequest(conn, commandEvents)
	if err != nil {

	} else {
		ApcInfo[commandEvents] = result
	}

	for key, value := range ApcInfo {
		for _, msg := range value {
			log.Print(key, " ==> ", msg)
		}
	}

}
