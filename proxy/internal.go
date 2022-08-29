package proxy

import (
	"context"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/encoding"
	"github.com/luraproject/lura/v2/transport/http/client"
)

func NewInternalProxy(remote *config.Backend, decode encoding.Decoder) Proxy {
	return NewInternalProxyWithHTTPExecutor(remote, decode)
}

func NewInternalProxyWithHTTPExecutor(remote *config.Backend, dec encoding.Decoder) Proxy {
	if remote.Encoding == encoding.NOOP {
		return NewInternalProxyDetailed(remote, client.NoOpHTTPStatusHandler, NoOpHTTPResponseParser)
	}

	ef := NewEntityFormatter(remote)
	rp := DefaultHTTPResponseParserFactory(HTTPResponseParserConfig{dec, ef})
	return NewInternalProxyDetailed(remote, client.GetHTTPStatusHandler(remote), rp)
}

func NewInternalProxyDetailed(backend *config.Backend, ch client.HTTPStatusHandler, rp HTTPResponseParser) Proxy {
	return func(ctx context.Context, request *Request) (*Response, error) {

		return RouteTable[backend.URLPattern](ctx, request)
	}
}
