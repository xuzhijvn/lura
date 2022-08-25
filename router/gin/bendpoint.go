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

	return func(configuration *config.EndpointConfig, prxy proxy.Proxy) proxy.Proxy {

		//rp := proxy.DefaultHTTPResponseParserFactory(proxy.HTTPResponseParserConfig{dec, ef})
		//ch := client.GetHTTPStatusHandler(remote)
		//
		//cacheControlHeaderValue := fmt.Sprintf("public, max-age=%d", int(configuration.CacheTTL.Seconds()))
		//isCacheEnabled := configuration.CacheTTL.Seconds() != 0
		//requestGenerator := NewRequest(configuration.HeadersToPass)
		//render := getRender(configuration)
		//logPrefix := "[ENDPOINT: " + configuration.Endpoint + "]"
		return func(ctx context.Context, requestToBakend *proxy.Request) (*proxy.Response, error) {

			resp, err := prxy(ctx, requestToBakend)

			select {
			case <-ctx.Done():
				if err == nil {
					err = server.ErrInternalError
				}
			default:
			}

			if requestToBakend.Body != nil {
				requestToBakend.Body.Close()
			}

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
			if err != nil {
				return nil, err
			}

			//resp, err = ch(ctx, resp)
			//if err != nil {
			//	if t, ok := err.(responseError); ok {
			//		return &proxy.Response{
			//			Data: map[string]interface{}{
			//				fmt.Sprintf("error_%s", t.Name()): t,
			//			},
			//			proxy.Metadata: proxy.Metadata{StatusCode: t.StatusCode()},
			//		}, nil
			//	}
			//	return nil, err
			//}

			return resp, nil
		}
	}
}

// NewRequest gets a request from the current gin context and the received query string
//func NewRequest(headersToSend []string) func(*gin.Context, []string) *proxy.Request {
//	if len(headersToSend) == 0 {
//		headersToSend = server.HeadersToSend
//	}
//
//	return func(c *gin.Context, queryString []string) *proxy.Request {
//		params := make(map[string]string, len(c.Params))
//		for _, param := range c.Params {
//			params[strings.Title(param.Key[:1])+param.Key[1:]] = param.Value
//		}
//
//		headers := make(map[string][]string, 3+len(headersToSend))
//
//		for _, k := range headersToSend {
//			if k == bRequestParamsAsterisk {
//				headers = c.Request.Header
//
//				break
//			}
//
//			if h, ok := c.Request.Header[textproto.CanonicalMIMEHeaderKey(k)]; ok {
//				headers[k] = h
//			}
//		}
//
//		headers["X-Forwarded-For"] = []string{c.ClientIP()}
//		headers["X-Forwarded-Host"] = []string{c.Request.Host}
//		// if User-Agent is not forwarded using headersToSend, we set
//		// the KrakenD router User Agent value
//		if _, ok := headers["User-Agent"]; !ok {
//			headers["User-Agent"] = server.UserAgentHeaderValue
//		} else {
//			headers["X-Forwarded-Via"] = server.UserAgentHeaderValue
//		}
//
//		query := make(map[string][]string, len(queryString))
//		queryValues := c.Request.URL.Query()
//		for i := range queryString {
//			if queryString[i] == bRequestParamsAsterisk {
//				query = c.Request.URL.Query()
//
//				break
//			}
//
//			if v, ok := queryValues[queryString[i]]; ok && len(v) > 0 {
//				query[queryString[i]] = v
//			}
//		}
//
//		return &proxy.Request{
//			Method:  c.Request.Method,
//			Query:   query,
//			Body:    c.Request.Body,
//			Params:  params,
//			Headers: headers,
//		}
//	}
//}
//
//type responseError interface {
//	error
//	StatusCode() int
//}
//
//type multiError interface {
//	error
//	Errors() []error
//}