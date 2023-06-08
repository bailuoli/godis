package tcp

import (
	"context"
	"net"
)

type Handler interface {
	Handle(c context.Context, conn net.Conn)
	Close() error
}
