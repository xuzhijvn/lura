// SPDX-License-Identifier: Apache-2.0

package gin

import (
	"context"
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/logging"
	"github.com/luraproject/lura/v2/proxy"
	"github.com/luraproject/lura/v2/transport/http/server"
)

// BHandlerFactory creates a handler function that adapts the gin router with the injected proxy
type BHandlerFactory func(*config.EndpointConfig, proxy.Proxy) proxy.Proxy

// BEndpointHandler implements the HandleFactory interface using the default ToHTTPError function
var BEndpointHandler = CustomErrorBEndpointHandler(logging.NoOp, server.DefaultToHTTPError)

// CustomErrorBEndpointHandler returns a HandleFactory using the injected ToHTTPError function and logger
func CustomErrorBEndpointHandler(logger logging.Logger, errF server.ToHTTPError) BHandlerFactory {

	return func(configuration *config.EndpointConfig, proxyStack proxy.Proxy) proxy.Proxy {

		return func(ctx context.Context, request *proxy.Request) (*proxy.Response, error) {

			resp, err := proxyStack(ctx, request)

			select {
			case <-ctx.Done():
				if err == nil {
					err = server.ErrInternalError
				}
			default:
			}

			if request.Body != nil {
				request.Body.Close()
			}

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			if err != nil {
				return nil, err
			}

			return resp, nil
		}
	}
}
