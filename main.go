package gonoolite

import (
	"fmt"
	"time"

	"go.bug.st/serial.v1"
)

type GoNoolite struct {
	port     serial.Port
	portName string
	baudrate int
	parity   serial.Parity
	stopBits serial.StopBits
}

func ListSerialPorts() ([]string, error) {
	return serial.GetPortsList()
}

type optionFunc func(this *GoNoolite)

func New(options ...optionFunc) (*GoNoolite, error) {
	new := &GoNoolite{portName: "/dev/ttyAMA0", baudrate: 9600, parity: serial.NoParity, stopBits: serial.OneStopBit}
	for _, o := range options {
		o(new)
	}
	mode := &serial.Mode{
		BaudRate: new.baudrate,
		Parity:   new.parity,
		StopBits: new.stopBits,
		DataBits: 8,
	}

	port, err := serial.Open(new.portName, mode)
	new.port = port
	return new, err
}

func (this *GoNoolite) sendCmd(content [14]byte) (err error) {
	buff := make([]byte, 17)
	buff[0] = 171
	buff[16] = 172
	sum := 171
	for i := 0; i < 14; i++ {
		buff[i+1] = content[i]
		sum += int(content[i])
	}
	buff[15] = byte(sum % 256)
	n := 0
	n, err = this.port.Write(buff)
	if err == nil && n != 17 {
		err = fmt.Errorf("cannot send data to serial port")
	}
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 300) //TODO
	n, err = this.port.Read(buff)
	if err == nil && n != 17 {
		err = fmt.Errorf("cannot receive answer from module. rcv count=%v", n)
	}
	if err != nil {
		return err
	}
	if buff[2] != 0 {
		return fmt.Errorf("noolite module error")
	}
	return nil
}

func (this *GoNoolite) SwitchOn(channel byte) error {
	buff := [14]byte{}
	buff[3] = channel
	buff[4] = 2
	return this.sendCmd(buff)
}

func (this *GoNoolite) SwitchOff(channel byte) error {
	buff := [14]byte{}
	buff[3] = channel
	buff[4] = 0
	return this.sendCmd(buff)
}
