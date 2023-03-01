package database

import "github.com/A-walker-ninght/miniRedis/interface/resp"

type CmdLine = [][]byte

type DataBase interface {
	Exec(client resp.Connection, args [][]byte) resp.Reply
	Close() error
	AfterClientClose(c resp.Connection) // 关闭后的善后
}

type DataEntry struct {
	Data interface{}
}
