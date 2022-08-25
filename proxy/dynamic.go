package proxy

import (
	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/transport/http/client"
)

var dynamicProxy = CustomDynamicProxyFactory(client.NewHTTPClient)

func CustomDynamicProxyFactory(cf client.HTTPClientFactory) BackendFactory {
	return func(backend *config.Backend) Proxy {
		if backend.Type == "internal" {
			return NewInternalProxy(backend, backend.Decoder)
		} else {
			return NewHTTPProxy(backend, cf, backend.Decoder)
		}
	}
}
