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
type EchoHandler struct {
	activeClient sync.Map
	closing      atomic.Boolean
}

func MakeEchoHandler() *EchoHandler {
	return &EchoHandler{}
}

func (c *EchoHandler) Handler(ctx context.Context, conn net.Conn) {
	// 判断连接池是否关闭
	if c.closing.Get() {
		_ = conn.Close() // 退出
	}

	client := &EchoClient{
		Conn: conn,
	}
	c.activeClient.Store(client, struct{}{})

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				logger.Info("connection close")
				c.activeClient.Delete(client)
			} else {
				logger.Warn(err)
			}
			return
		}
		client.Wating.Add(1)
		b := []byte(msg)
		_, _ = client.Conn.Write(b)
		client.Wating.Done()
	}
}

func (c *EchoHandler) Close() error {
	logger.Info("handler shutting down...") // 连接池关闭
	c.closing.Set(true)

	// 关闭活跃的连接
	c.activeClient.Range(func(key, value any) bool {
		client := key.(*EchoClient)
		_ = client.Close()
		return true
	})
	return nil
}

// 客户端连接
type EchoClient struct {
	Conn   net.Conn
	Wating wait.Wait
}

func (c *EchoClient) Close() error {
	c.Wating.WaitWithTimeout(10 * time.Second) // 等待业务处理完，超时直接关闭
	_ = c.Conn.Close()
	return nil
}
