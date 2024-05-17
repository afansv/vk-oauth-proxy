package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy takes target host and creates a reverse proxy
func NewProxy(targetHost string, modifier func(resp *http.Response) error) (*httputil.ReverseProxy, error) {
	u, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(u)
	director := proxy.Director

	hostHeader, err := getUpstreamHostHeader(targetHost)
	if err != nil {
		return nil, fmt.Errorf("getUpstreamHostHeader: %w", err)
	}

	proxy.Director = func(request *http.Request) {
		director(request)
		request.Host = hostHeader
		request.Header.Del("Accept-Encoding")
	}
	proxy.ModifyResponse = modifier
	return proxy, nil
}

func getUpstreamHostHeader(upstreamHost string) (string, error) {
	u, err := url.Parse(upstreamHost)
	if err != nil {
		return "", fmt.Errorf("url.Parse: %w", err)
	}
	return u.Host, nil
}
