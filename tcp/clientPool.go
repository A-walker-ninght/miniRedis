package tcp

import (
	"bufio"
	"context"
	"io"
	"net"
	"sync"
	"time"

	"github.com/A-walker-ninght/miniRedis/lib/logger"
	"github.com/A-walker-ninght/miniRedis/lib/sync/atomic"
	"github.com/A-walker-ninght/miniRedis/lib/sync/wait"
)

// 客户端连接池
type ClientPool struct {
	activeClient sync.Map
	closing      atomic.Boolean
}

func MakeClientPool() *ClientPool {
	return &ClientPool{}
}

func (c *ClientPool) Handler(ctx context.Context, conn net.Conn, closeCh chan struct{}) {
	// 判断连接池是否关闭
	if c.closing.Get() {
		_ = conn.Close() // 退出
	}

	client := &ClientHandler{
		conn: conn,
	}
	c.activeClient.Store(client, struct{}{})

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				_ = conn.Close()
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Wating.Add(1)
		b := []byte(msg)
		client.conn.Write(b)
		client.Wating.Done()
	}
}

func (c *ClientPool) Close() error {
	logger.Info("handler shutting down...") // 连接池关闭
	c.closing.Set(true)

	// 关闭活跃的连接
	c.activeClient.Range(func(key, value any) bool {
		client := key.(*ClientHandler)
		_ = client.Close()
		return true
	})
	return nil
}

// 客户端连接
type ClientHandler struct {
	conn   net.Conn
	Wating wait.Wait
}

func (c *ClientHandler) Close() error {
	c.Wating.WaitWithTimeout(10 * time.Second) // 等待业务处理完，超时直接关闭
	c.conn.Close()
	return nil
}
