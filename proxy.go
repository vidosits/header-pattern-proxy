// Package header_pattern_proxy
package header_pattern_proxy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

// Config the plugin configuration.
type Config struct {
	Header  string            `json:"headers,omitempty"` // target header
	Mapping map[string]string `json:"mapping,omitempty"` // mapping holding (regex, target) pairs
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Header:  "",
		Mapping: make(map[string]string),
	}
}

type SiteProxy struct {
	config *Config
	next   http.Handler
	name   string
}

// New created a new SiteProxy plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.Header) == 0 {
		return nil, fmt.Errorf("header cannot be empty")
	}

	if len(config.Mapping) == 0 {
		return nil, fmt.Errorf("mapping cannot be empty")
	}

	return &SiteProxy{
		config: config,
		next:   next,
		name:   name,
	}, nil
}

func (a *SiteProxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for pattern, destination := range a.config.Mapping {
		matched, _ := regexp.MatchString(pattern, req.Header.Get(a.config.Header))

		if matched {
			destinationUrl, err := url.Parse(destination)

			if err != nil {
				continue
			}

			proxy := httputil.NewSingleHostReverseProxy(destinationUrl)

			req.URL.Host = destinationUrl.Host
			req.URL.Scheme = "http"
			req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
			req.Host = destinationUrl.Host

			proxy.ServeHTTP(rw, req)
		}
	}

	a.next.ServeHTTP(rw, req)
}
