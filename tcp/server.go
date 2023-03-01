package tcp

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/A-walker-ninght/miniRedis/interface/tcp"
	"github.com/A-walker-ninght/miniRedis/lib/logger"
)

type Config struct {
	Address string
}

func ListenAndServeWithSignal(cfg *Config, handler tcp.Handler) error {
	closeCh := make(chan struct{})
	sigCh := make(chan os.Signal)

	// 操作系统监听退出
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		switch sig {
		case syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
			closeCh <- struct{}{}
		}
	}()

	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("bind: %s, start listening...", cfg.Address))
	ListenAndServe(listener, handler, closeCh)

	return nil
}

func ListenAndServe(listener net.Listener, handler tcp.Handler, closeCh chan struct{}) {
	// 优雅退出：1.操作系统退出 2.用户退出 3.连接错误

	// 操作系统
	go func() {
		<-closeCh
		logger.Info("shutting down...")
		_ = listener.Close()
		_ = handler.Close()
	}()

	// 延迟关闭
	defer func() {
		_ = listener.Close()
		_ = handler.Close()
	}()

	ctx := context.Background()
	var wg sync.WaitGroup

	for {
		conn, err := listener.Accept()

		if err != nil {
			break
		}
		logger.Info("accept link")

		// 处理业务，限时
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()
<<<<<<< HEAD
			handler.Handler(ctx, conn, closeCh)
=======
			handler.Handler(ctx, conn)
>>>>>>> 70f3717 (resp 2023.3.1)
		}()
	}
	wg.Wait() // 等待业务处理完
}
