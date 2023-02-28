package tcp

import (
	"context"
	"net"
)

type Handler interface {
	Handler(ctx context.Context, conn net.Conn, closeCh chan struct{})
	Close() error
}
