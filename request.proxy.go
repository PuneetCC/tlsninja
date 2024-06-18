package tlsninja

import (
	"context"
)

type ProxyRequest struct {
	Proxy string
}

func NewProxyRequest(proxy string) ProxyRequest {
	return ProxyRequest{Proxy: proxy}
}

func (p *ProxyRequest) Do(config IRequestConfig) (*IRequestResponse, error) {
	proxyProvider := func(url string) string { return p.Proxy }
	client := NewHTTPClient(context.Background(), config.JA3Fingerprint, &proxyProvider)

	// default timeout - 10s
	if config.Timeout == 0 {
		config.Timeout = 10
	}

	resp, err := client.Request(RequestConfig{
		Method:      config.Method,
		URL:         config.URL,
		QueryParams: config.QueryParams,
		Payload:     config.Payload,
		Headers:     config.Headers,
		Timeout:     config.Timeout,
	})

	if err != nil {
		return nil, err
	}

	return &IRequestResponse{
		StatusCode: resp.StatusCode,
		Body:       resp.Body,
		Headers:    resp.Headers,
	}, nil
}
