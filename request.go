package gonoolite

type Request struct {
	content [14]byte
}

func (t *Request) Mode(v EnumMode) *Request {
	t.content[0] = byte(v)
	return t
}

type RequestEnumCTR int

const (
	SendCmdToEndDevice       RequestEnumCTR = 0
	SendBroadcastCmd         RequestEnumCTR = 1
	ReadRcvBuf               RequestEnumCTR = 2
	EnableBinding            RequestEnumCTR = 3
	DisableBinding           RequestEnumCTR = 4
	ClearMemoredChannel      RequestEnumCTR = 5
	ClearAllMemored          RequestEnumCTR = 6
	UnbindAddressFromChannel RequestEnumCTR = 7
	SendCmdToGivenNLFAddress RequestEnumCTR = 8
)

func (t *Request) Control(ctr RequestEnumCTR, numRetries int) *Request {
	if numRetries > 3 {
		numRetries = 3
	}
	if numRetries < 0 {
		numRetries = 0
	}
	t.content[1] = byte(ctr) | (byte(numRetries) << 6)
	return t
}

func (t *Request) Channel(channel int) *Request {
	if channel > 63 {
		channel = 63
	}
	if channel < 0 {
		channel = 0
	}
	t.content[3] = byte(channel)
	return t
}

func (t *Request) Address(addr uint32) *Request {
	for i := 3; i >= 0; i-- {
		t.content[13-i] = byte(addr >> uint(8*i))
	}
	return t
}

func (t *Request) D2(value byte) *Request {
	t.content[8] = value
	return t
}

func (t *Request) Data(b0 byte, b1 byte, b2 byte, b3 byte) *Request {
	t.content[6] = b0
	t.content[7] = b1
	t.content[8] = b2
	t.content[9] = b3
	return t
}

func (t *Request) CommandToSend(cmd EnumCMD) *Request {
	t.content[4] = byte(cmd)
	return t
}

func (t *Request) Serialize() []byte {
	buff := make([]byte, 17)
	buff[0] = 171
	buff[16] = 172
	sum := 171
	for i := 0; i < 14; i++ {
		buff[i+1] = t.content[i]
		sum += int(t.content[i])
	}
	buff[15] = byte(sum % 256)
	return buff
}
