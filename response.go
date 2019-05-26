package gonoolite

import (
	"fmt"
	"io"
)

type Response struct {
	content [14]byte
}

func (r *Response) Mode() EnumMode {
	return EnumMode(r.content[0])
}

func (r *Response) D2() byte {
	return r.content[8]
}

type ResponseEnumCTR int

const (
	CommandDone ResponseEnumCTR = 0
	NoAnswer    ResponseEnumCTR = 1
	ExecErr     ResponseEnumCTR = 2
	BindingDone ResponseEnumCTR = 3
)

func (r *Response) CTR() ResponseEnumCTR {
	return ResponseEnumCTR(r.content[1])
}

func (r *Response) Channel() int {
	return int(r.content[3])
}

func (r *Response) Command() EnumCMD {
	return EnumCMD(r.content[4])
}

func (t *Response) parseAnswer(content []byte) error {

	if len(content) != 17 {
		return fmt.Errorf("cannot receive answer from module")
	}
	if content[0] != 173 || content[16] != 174 {
		return fmt.Errorf("Wrong pkg structure")
	}

	sum := 0
	for i := 0; i < 15; i++ {
		sum += int(content[i])
	}
	if byte(sum%256) != content[15] {
		return fmt.Errorf("Wrong pkg checksum")
	}
	copy(t.content[:], content[1:15])
	return nil
}

func (t *Response) Receive(reader io.Reader) error {
	buff := make([]byte, 1)
	resp := []byte{}
	for true {
		_, err := reader.Read(buff)
		if err != nil {
			return err
		}
		if buff[0] == 173 {
			resp = append(resp, 173)
			break
		}
	}

	for true {
		_, err := reader.Read(buff)
		if err != nil {
			return err
		}
		resp = append(resp, buff[0])
		if buff[0] == 174 {
			break
		}
	}

	return t.parseAnswer(resp)
}
