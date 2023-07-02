package netkit

import "net/http"

func NewTransport(options ...func(*http.Transport)) http.RoundTripper {
	t := http.DefaultTransport.(*http.Transport).Clone()
	for _, option := range options {
		option(t)
	}

	return t
}

func WithIdleConnsPerHost(maxIdleConsPerHost int) func(*http.Transport) {
	return func(t *http.Transport) {
		t.MaxIdleConnsPerHost = maxIdleConsPerHost
	}
}
