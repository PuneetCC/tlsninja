package tlsninja

import (
	"context"
	"math/rand"
	"net"
	"net/url"
	"time"

	"github.com/Danny-Dasilva/CycleTLS/cycletls"
)

type RequestConfig struct {
	Method        string
	URL           string
	QueryParams   map[string]string
	Payload       []byte
	Headers       map[string]string
	Timeout       int
	MaxRetries    int
	RetryDelay    time.Duration
	RetryableFn   func(*cycletls.Response, error) bool
	SkipRedirects bool `json:"skipRedirects"`
}

type HTTPClientResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       []byte `json:"body"`
	Headers    map[string]string
}

type ContextHttpClientKey struct{}

type HTTPClient struct {
	Context           context.Context
	Client            cycletls.CycleTLS
	AdditionalHeaders map[string]string
	JA3Fingerprint    string
	ProxyProvider     *func(url string) string
}

func NewHTTPClient(ctx context.Context, ja3 string, proxyProvider *func(url string) string) *HTTPClient {
	return &HTTPClient{
		Context:        ctx,
		Client:         cycletls.CycleTLS{},
		JA3Fingerprint: ja3,
		ProxyProvider:  proxyProvider,
	}
}

// retries only in case of network errors.
func defaultRetryableFn(resp *cycletls.Response, err error) bool {
	if err != nil {
		// Retry if there was a network error.
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			return true
		}
	}
	return false
}

func (h *HTTPClient) Request(config RequestConfig) (*HTTPClientResponse, error) {

	// Prepare headers
	headers := map[string]string{}
	for key, value := range config.Headers {
		headers[key] = value
	}
	for key, value := range h.AdditionalHeaders {
		headers[key] = value
	}

	baseUrl, _ := url.Parse(config.URL)

	// Prepare query params
	params := url.Values{}
	for key, value := range config.QueryParams {
		params.Add(key, value)
	}

	// Set the raw query of the URL to the encoded query parameters.
	baseUrl.RawQuery = params.Encode()

	finalURL := baseUrl.String()

	// Prepare options for CycleTLS
	options := cycletls.Options{
		Body:    string(config.Payload),
		Ja3:     h.JA3Fingerprint,
		Headers: headers,
		Timeout: config.Timeout,
	}

	if config.SkipRedirects {
		options.DisableRedirect = config.SkipRedirects
	}

	if h.ProxyProvider != nil {
		provider := *h.ProxyProvider
		options.Proxy = provider(finalURL)
	}

	if userAgent, ok := headers["User-Agent"]; ok {
		options.UserAgent = userAgent
	}

	// Retry logic
	var resp cycletls.Response
	var err error
	retryCount := 0

	if config.MaxRetries == 0 {
		resp, err = h.Client.Do(finalURL, options, config.Method)
	} else {
		for {
			resp, err = h.Client.Do(finalURL, options, config.Method)
			if config.RetryableFn == nil {
				config.RetryableFn = defaultRetryableFn
			}
			if err == nil && !config.RetryableFn(&resp, err) {
				break
			}
			if retryCount >= config.MaxRetries {
				break
			}
			retryCount++
			if config.RetryDelay == 0 {
				config.RetryDelay = time.Duration(rand.Intn(800)+200) * time.Millisecond
			}
			time.Sleep(config.RetryDelay)
		}
	}

	if err != nil {
		return nil, err
	}

	// Convert headers to the desired format
	responseHeaders := map[string]string{}
	for key, values := range resp.Headers {
		responseHeaders[key] = values
	}

	return &HTTPClientResponse{
		StatusCode: resp.Status,
		Body:       []byte(resp.Body),
		Headers:    responseHeaders,
	}, nil
}

func (h *HTTPClient) Get(url string, query map[string]string, headers map[string]string) (*HTTPClientResponse, error) {
	return h.Request(RequestConfig{
		Method:      "GET",
		URL:         url,
		QueryParams: query,
		Headers:     headers,
	})
}

func (h *HTTPClient) Post(url string, payload []byte, query map[string]string, headers map[string]string) (*HTTPClientResponse, error) {
	return h.Request(RequestConfig{
		Method:      "POST",
		URL:         url,
		QueryParams: query,
		Payload:     payload,
		Headers:     headers,
	})
}

func (h *HTTPClient) Put(url string, payload []byte, query map[string]string, headers map[string]string) (*HTTPClientResponse, error) {
	return h.Request(RequestConfig{
		Method:      "PUT",
		URL:         url,
		QueryParams: query,
		Payload:     payload,
		Headers:     headers,
	})
}

func (h *HTTPClient) Patch(url string, payload []byte, query map[string]string, headers map[string]string) (*HTTPClientResponse, error) {
	return h.Request(RequestConfig{
		Method:      "PATCH",
		URL:         url,
		QueryParams: query,
		Payload:     payload,
		Headers:     headers,
	})
}

func (h *HTTPClient) Delete(url string, payload []byte, query map[string]string, headers map[string]string) (*HTTPClientResponse, error) {
	return h.Request(RequestConfig{
		Method:      "DELETE",
		URL:         url,
		QueryParams: query,
		Payload:     payload,
		Headers:     headers,
	})
}
