package jsonplaceholder

import "net/http"

// Option represents the client options
type Option func(*Client)

// WithHTTPClient sets a custom http client
func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}
