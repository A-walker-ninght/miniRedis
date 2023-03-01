package database

import (
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}
func (e *EchoDatabase) Exec(client resp.Connection, args [][]byte) resp.Reply {
	return reply.MakeMultiBulkReply(args)
}

func (e *EchoDatabase) Close() error {
	//TODO implement me
	panic("implement me")
}

func (e *EchoDatabase) AfterClientClose(c resp.Connection) {
	//TODO implement me
	panic("implement me")
}
