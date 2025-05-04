package tlsninja

type IRequestConfig struct {
	Method             string            `json:"method"`
	URL                string            `json:"url"`
	QueryParams        map[string]string `json:"queryParams"`
	Payload            []byte            `json:"payload"`
	Headers            map[string]string `json:"headers"`
	Timeout            int               `json:"timeout"`
	JA3Fingerprint     string            `json:"ja3"`
	SkipRedirects      bool              `json:"skipRedirects"`
	HexEncodedResponse bool              `json:"hexEncodedResponse"`
}

type IRequestResponse struct {
	StatusCode int               `json:"statusCode"`
	Body       []byte            `json:"body"`
	Headers    map[string]string `json:"headers"`
}

type IRequest interface {
	Do(config IRequestConfig) (*IRequestResponse, error)
}
