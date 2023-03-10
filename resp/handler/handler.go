package handler

import (
	"context"
	database2 "github.com/A-walker-ninght/miniRedis/database"
	"github.com/A-walker-ninght/miniRedis/interface/database"
	"github.com/A-walker-ninght/miniRedis/lib/logger"
	"github.com/A-walker-ninght/miniRedis/lib/sync/atomic"
	"github.com/A-walker-ninght/miniRedis/resp/connetion"
	"github.com/A-walker-ninght/miniRedis/resp/parser"
	"github.com/A-walker-ninght/miniRedis/resp/reply"
	"github.com/hdt3213/godis/redis/connection"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db         database.DataBase
	closing    atomic.Boolean
}

func MakeHandler() *RespHandler {
	var db database.DataBase
	db = database2.NewStandaloneDatabase()
	return &RespHandler{
		db: db,
	}
}

func (r *RespHandler) Handler(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		_ = conn.Close()
	}
	client := connection.NewConn(conn)
	r.activeConn.Store(client, struct{}{})

	ch := parser.ParseStream(conn)

	for payload := range ch {
		if payload.Err != nil {
			if payload.Err == io.EOF || payload.Err == io.ErrUnexpectedEOF ||
				strings.Contains(payload.Err.Error(), "use of closed network connection") {
				r.CloseClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}

			// protocol error
			errReply := reply.MakeStandardErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.CloseClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}

		if payload.Data == nil {
			continue
		}
		repl, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(client, repl.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(unknownErrReplyBytes)
		}
	}
}

func (r *RespHandler) CloseClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down")
	r.closing.Set(true)
	r.activeConn.Range(func(key interface{}, value interface{}) bool {
		client := key.(*connetion.Connection)
		_ = client.Close()
		return true
	})
	_ = r.db.Close()
	return nil
}
