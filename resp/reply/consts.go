package reply

type PongReply struct {
}

var pongBytes = []byte("+Pong\r\n")

func (r PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}
