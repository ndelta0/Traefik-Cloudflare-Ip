package traefik_cloudflare_ip

import (
	"context"
	"net/http"
)

type Config struct {
}

func CreateConfig() *Config {
	return &Config{}
}

type Plugin struct {
	next http.Handler
	name string
}

func (p Plugin) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	trueIp := request.Header.Get("cf-connecting-ip")
	if len(trueIp) > 0 {
		request.Header.Set("x-real-ip", trueIp)

		var trueForwarded string
		forwarded := request.Header.Get("x-forwarded-for")
		if len(forwarded) > 0 {
			trueForwarded = trueIp + "," + forwarded
		} else {
			trueForwarded = trueIp
		}

		request.Header.Set("x-forwarded-for", trueForwarded)
	}

	p.next.ServeHTTP(writer, request)
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	plugin := &Plugin{
		next: next,
		name: name,
	}

	return plugin, nil
}
