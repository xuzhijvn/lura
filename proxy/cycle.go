package proxy

import (
	"context"
)

func Handle(uri string, ctx context.Context, request *Request) (*Response, error) {

	return Proxys[uri](ctx, request)

}
