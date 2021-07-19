package atuhdsw52ed

import (
	"context"
	"net"
	"time"

	"github.com/byuoitav/connpool"
	"go.uber.org/zap"
)

const (
	asciiCarriageReturn = 0x0d
	asciiLineFeed       = 0x0a
)

type AtlonaVideoSwitcher5x1 struct {
	pool *connpool.Pool
	log  *zap.Logger
}

func NewAtlonaVideoSwitcher5x1(addr string, opts ...Option) *AtlonaVideoSwitcher5x1 {
	options := &options{
		ttl:   30 * time.Second,
		delay: 500 * time.Millisecond,
		log:   zap.NewNop(),
	}

	for _, o := range opts {
		o.apply(options)
	}

	return &AtlonaVideoSwitcher5x1{
		log: options.log,
		pool: &connpool.Pool{
			TTL:   options.ttl,
			Delay: options.delay,
			NewConnection: func(ctx context.Context) (net.Conn, error) {
				dial := net.Dialer{}
				return dial.DialContext(ctx, "tcp", addr+":23")
			},
			Logger: options.log.Sugar(),
		},
	}
}
