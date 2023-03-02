package database

import (
	"github.com/A-walker-ninght/miniRedis/interface/resp"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
)

func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.MakePongReply()
}

func init() {
	RegisterCommand("ping", Ping, 1)
}
