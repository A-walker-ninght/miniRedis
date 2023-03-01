package reply

type UnknowErrReply struct{}

var unknowErrBytes = []byte("-Err unknown\r\n")

func (u *UnknowErrReply) Error() string {
	return "Err unknown"
}

func (u *UnknowErrReply) ToBytes() []byte {
	return unknowErrBytes
}

/* cmd Error */
type ArgNumErrReply struct {
	Cmd string
}

func (a *ArgNumErrReply) Error() string {
	return "ERR wrong number of arguments for '" + a.Cmd + "'command"
}

func (a *ArgNumErrReply) ToBytes() []byte {
	return []byte("-ERR wrong number of arguments for '" + a.Cmd + "'command\r\n")
}

func MakeArgNumErrReply(cmd string) *ArgNumErrReply {
	return &ArgNumErrReply{}
}

/*Syntax Error*/
type SyntaxErrReply struct{}

var syntaxErrBytes = []byte("-Err syntax error\r\n")

func (s *SyntaxErrReply) Error() string {
	return "Err syntax error"
}

func (s *SyntaxErrReply) ToBytes() []byte {
	return syntaxErrBytes
}

func MakeSyntaxErrReply() *SyntaxErrReply {
	return &SyntaxErrReply{}
}

/*Wrong Type Error*/

type WrongTypeErrReply struct{}

var wrongTypeErrBytes = []byte("-WRONGTYPE Operation against a key hoding the kind of value\r\n")

func (w *WrongTypeErrReply) Error() string {
	return "WRONGTYPE Operation against a key hoding the kind of value"
}

func (w *WrongTypeErrReply) ToBytes() []byte {
	return wrongTypeErrBytes
}

/*protocal Error*/

type ProtocalErrReply struct {
	Msg string
}

func (p *ProtocalErrReply) Error() string {
	return "ERR Protocal error: '" + p.Msg
}

func (p *ProtocalErrReply) ToBytes() []byte {
	return []byte("-ERR Protocal error: '" + p.Msg + "'\r\n")
}
