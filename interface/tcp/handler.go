package tcp

import (
	"context"
	"net"
)

type Handler interface {
<<<<<<< HEAD
	Handler(ctx context.Context, conn net.Conn, closeCh chan struct{})
=======
	Handler(ctx context.Context, conn net.Conn)
>>>>>>> 70f3717 (resp 2023.3.1)
	Close() error
}
