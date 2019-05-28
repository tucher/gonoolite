package gonoolite

import (
	"log"
	"sync"
	"time"

	"go.bug.st/serial.v1"
)

type GoNoolite struct {
	port     serial.Port
	portName string

	response Response

	OnState  func(devID uint32, state bool)
	OnBinded func(channel int, devID uint32)

	sendChannel chan []byte
	rcvFlagChan chan bool
	checking    bool
	mtx         sync.Mutex
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

func (t *GoNoolite) reader() {
	t.port.ResetInputBuffer()
	for {
		err := t.response.Receive(t.port)
		// log.Println("Read")
		select {
		case t.rcvFlagChan <- true:
		default:
		}
		if err != nil {
			log.Printf("RCV ERROR: %+v", err)
			continue
		}
		if t.response.Mode() == FTX && t.response.CTR() == CommandDone && t.response.Command() == Send_State {

			if t.response.D2() == 1 {
				if t.OnState != nil {
					t.OnState(t.response.DevID(), true)
				}
			} else if t.response.D2() == 0 {
				if t.OnState != nil {
					t.OnState(t.response.DevID(), false)
				}
			} else {
				log.Printf("Bad status")
			}
			continue
		}
		if t.response.Mode() == FTX && t.response.CTR() == BindingDone {
			if t.OnBinded != nil {
				t.OnBinded(t.response.Channel(), t.response.DevID())
			}
			continue
		}

		log.Printf("Unknown response: %+v", t.response)

	}
}

func (t *GoNoolite) sender() {
	t.port.ResetOutputBuffer()
	time.Sleep(time.Second * 10)
	// start := time.Now()
	for {
		data := []byte{}
		select {
		case data = <-t.sendChannel:
		default:
			if t.IsPolling() {
				r := Request{}
				r.Mode(FTX).Control(SendBroadcastCmd, 0).Channel(0).CommandToSend(Read_State)
				data = r.Serialize()
			}
		}

		if len(data) > 0 {
			t.port.Write(data)
			// log.Printf("Since Write: %+v", time.Since(start))
			// start = time.Now()
		}
		tmr := time.After(time.Millisecond * 700)
	W:
		for {
			select {
			case <-tmr:
				break W
			case <-t.rcvFlagChan:
				tmr = time.After(time.Millisecond * 700)
			}
		}
	}
}
func New(options ...optionFunc) (*GoNoolite, error) {
	new := &GoNoolite{
		portName:    "/dev/ttyAMA0",
		sendChannel: make(chan []byte),
		rcvFlagChan: make(chan bool),
	}
	for _, o := range options {
		o(new)
	}
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
		DataBits: 8,
	}

	port, err := serial.Open(new.portName, mode)
	new.port = port
	if err == nil {

		go new.reader()
		go new.sender()

		// r := Request{}
		// r.Mode(SVC)
		// new.sendChannel <- r.Serialize()

	}

	return new, err
}
func (t *GoNoolite) SetPolling(st bool) {

	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.checking = st
}

func (t *GoNoolite) IsPolling() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.checking
}

func (t *GoNoolite) SetState(channel int, st bool) {
	r := Request{}
	r.Mode(FTX)
	r.Control(SendBroadcastCmd, 0)
	r.Channel(channel)
	if st {
		r.CommandToSend(On)
	} else {
		r.CommandToSend(Off)
	}
	t.sendChannel <- r.Serialize()
}

func (t *GoNoolite) Bind(channel int) {
	r := Request{}
	r.Mode(FTX)
	r.Control(SendCmdToEndDevice, 0)
	r.Channel(channel)
	r.CommandToSend(Bind)
	t.sendChannel <- r.Serialize()
}

func (t *GoNoolite) StartBinding(channel int) {
	r := Request{}
	r.Mode(FTX)
	r.Control(SendCmdToEndDevice, 0)
	r.Channel(channel)
	r.CommandToSend(Service)
	r.Data(1, 0, 0, 0)
	t.sendChannel <- r.Serialize()
}

// func (t *GoNoolite) Unbind(channel byte) error {
// 	buff := [14]byte{}
// 	buff[3] = channel
// 	buff[4] = 9
// 	_, err := t.sendCmd(buff)
// 	return err
// }
