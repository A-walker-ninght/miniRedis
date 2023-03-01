package reply

<<<<<<< HEAD
=======
/* pong reply */
>>>>>>> 70f3717 (resp 2023.3.1)
type PongReply struct {
}

var pongBytes = []byte("+Pong\r\n")

func (r PongReply) ToBytes() []byte {
	return pongBytes
}

func MakePongReply() *PongReply {
	return &PongReply{}
}
<<<<<<< HEAD
=======

/* ok reply */

type OKReply struct{}

var okbytes = []byte("+OK\r\n")

func (o *OKReply) ToBytes() []byte {
	return okbytes
}

func MakeOKReply() *OKReply {
	return &OKReply{}
}

/* null string reply */
type NullBulkReply struct{}

var nullBulkBytes = []byte("$-1\r\n")

func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

func MakeNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}

/* Empty list reply */
type EmptyMultiBulkReply struct{}

var emptyMultiBulkbytes = []byte("*0\r\n")

func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkbytes
}

func MakeEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}

/* no reply*/
type NoReply struct{}

var noBytes = []byte("")

func (n *NoReply) ToBytes() []byte {
	return noBytes
}

func MakeNoReply() *NoReply {
	return &NoReply{}
}
>>>>>>> 70f3717 (resp 2023.3.1)
