package oauth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ContentType string

const (
	ContentTypeJson    ContentType = "application/json;charset=utf-8"
	ContentTypeWWWForm ContentType = "application/x-www-form-urlencoded"
)

type Request struct {
	// URL is the target URL for the request.
	URL string

	// Method is the HTTP method (e.g., GET, POST) for the request.
	Method string

	// ProxyURL is an optional proxy URL to be used during the request.
	ProxyURL string

	// ContentType is the type of content being sent in the request body.
	ContentType ContentType

	// Timeout is the maximum duration for the request to complete.
	Timeout time.Duration

	// Header contains any additional headers to be included in the request.
	Header http.Header

	// Data contains the form values for the request body.
	Data url.Values
}

type ROption func(*Request)

// WithContentType sets the ContentType option for the Request.
func WithContentType(contentType ContentType) ROption {
	return func(request *Request) {
		request.ContentType = contentType
	}
}

// WithTimeout sets the Timeout option for the Request.
func WithTimeout(timeout time.Duration) ROption {
	return func(request *Request) {
		request.Timeout = timeout
	}
}

// WithHeader sets the Header option for the Request.
func WithHeader(header http.Header) ROption {
	return func(request *Request) {
		request.Header = header
	}
}

// WithData sets the Data option for the Request.
func WithData(data url.Values) ROption {
	return func(request *Request) {
		request.Data = data
	}
}

// formatParams converts the request data to the appropriate format based on the content type.
func (req *Request) formatParams() io.Reader {
	if len(req.Data) > 0 {
		if req.ContentType == "" {
			if v, ok := req.Header["Content-Type"]; ok {
				req.ContentType = ContentType(v[0])
			}
		}

		switch req.ContentType {
		case ContentTypeJson:
			value, err := json.Marshal(req.Data)
			if err != nil {
				panic(err)
			}
			return bytes.NewReader(value)
		case ContentTypeWWWForm:
			return strings.NewReader(req.Data.Encode())
		}
	}
	return nil
}

// newRequest creates a new http.Request based on the Request parameters.
func (req *Request) newRequest() *http.Request {
	request, err := http.NewRequest(req.Method, req.URL, req.formatParams())
	if err != nil {
		panic(err)
	}
	request.Header = req.Header
	return request
}

// setProxy configures the proxy for the request, if a ProxyURL is provided.
func (req *Request) setProxy() *http.Transport {
	var proxy *http.Transport = nil
	if req.ProxyURL != "" {
		u, err := url.Parse(req.ProxyURL)
		if err != nil {
			panic(err)
		}
		proxy = &http.Transport{
			Proxy: http.ProxyURL(u),
		}
	}
	return proxy
}

// httpClient returns an HTTP client based on the given request configuration.
func (req *Request) httpClient() http.Client {
	client := http.Client{}

	// If no timeout value is specified in the request, set the timeout value for the client to default to 5 seconds.
	if 0 >= req.Timeout {
		req.Timeout = 5 * time.Second
	}
	client.Timeout = req.Timeout

	// Check for proxy configuration and set the client's transport accordingly.
	if proxy := req.setProxy(); proxy != nil {
		client.Transport = proxy
	}

	return client
}

// Do Execute an HTTP request and return the response
func (req *Request) Do() (*http.Response, error) {
	client := req.httpClient()
	return client.Do(req.newRequest())
}

// Post Execute an POST request and return the response.
func (req *Request) Post() (*http.Response, error) {
	req.Method = http.MethodPost
	return req.Do()
}

// Get Execute an GET request and return the response.
func (req *Request) Get() (*http.Response, error) {
	req.Method = http.MethodGet
	return req.Do()
}

// New creates a new Request with the specified URL, method, proxy, and options.
func New(url, method, proxy string, options ...ROption) *Request {
	request := &Request{
		URL:      url,
		Method:   method,
		ProxyURL: proxy,
	}
	for _, option := range options {
		option(request)
	}
	return request
}
