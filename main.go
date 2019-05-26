package gonoolite

import (
	"log"
	"time"

	"go.bug.st/serial.v1"
)

type GoNoolite struct {
	port     serial.Port
	portName string
	baudrate int
	parity   serial.Parity
	stopBits serial.StopBits

	response Response

	OnState  func(channel int, state bool)
	OnBinded func(channel int)

	states [64]bool
}

func ListSerialPorts() ([]string, error) {
	return serial.GetPortsList()
}

type optionFunc func(this *GoNoolite)

func WithPort(p string) optionFunc {
	return func(this *GoNoolite) {
		this.portName = p
	}
}

func (t *GoNoolite) statesChecker() {
	for {
		time.Sleep(time.Millisecond * 1000)
		t.ReadState(1)
	}
}
func (t *GoNoolite) reader() {
	for {
		err := t.response.Receive(t.port)
		if err != nil {
			log.Printf("RCV ERROR: %+v", err)
			continue
		}
		if t.response.Mode() == FTX && t.response.CTR() == CommandDone && t.response.Command() == Send_State {
			if t.response.D2() == 1 {
				t.states[t.response.Channel()] = true
				if t.OnState != nil {
					t.OnState(t.response.Channel(), true)
				}
			} else if t.response.D2() == 0 {
				t.states[t.response.Channel()] = false
				if t.OnState != nil {
					t.OnState(t.response.Channel(), false)
				}
			} else {
				log.Printf("Bad status")
			}
			continue
		}
		if t.response.Mode() == FTX && t.response.CTR() == BindingDone {
			if t.OnBinded != nil {
				t.OnBinded(t.response.Channel())
			}
			continue
		}

		log.Printf("Unknown response")

	}
}
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
	if err == nil {
		port.ResetInputBuffer()
		port.ResetOutputBuffer()

		r := Request{}
		r.Mode(SVC)
		r.Send(port)

		go new.reader()
		go new.statesChecker()
	}

	return new, err
}

func (t *GoNoolite) SwitchOn(channel int) error {
	r := Request{}
	r.Mode(FTX)
	r.Control(SendCmdToEndDevice, 0)
	r.Channel(channel)
	r.CommandToSend(On)
	return r.Send(t.port)
}
func (t *GoNoolite) SwitchOff(channel int) error {
	r := Request{}
	r.Mode(FTX)
	r.Control(SendCmdToEndDevice, 0)
	r.Channel(channel)
	r.CommandToSend(Off)
	return r.Send(t.port)
}

func (t *GoNoolite) Bind(channel int) error {
	r := Request{}
	r.Mode(FTX)
	r.Control(SendCmdToEndDevice, 0)
	r.Channel(channel)
	r.CommandToSend(Bind)
	return r.Send(t.port)
}

func (t *GoNoolite) ReadState(channel int) error {
	r := Request{}
	r.Mode(FTX)
	r.Control(SendBroadcastCmd, 0)
	r.Channel(channel)
	r.CommandToSend(Read_State)
	return r.Send(t.port)
}

// func (t *GoNoolite) Unbind(channel byte) error {
// 	buff := [14]byte{}
// 	buff[3] = channel
// 	buff[4] = 9
// 	_, err := t.sendCmd(buff)
// 	return err
// }
