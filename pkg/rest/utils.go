package rest

import "net/http"

func GetScheme(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if v := r.Header.Get("X-Forwarded-Proto"); v != "" {
		scheme = v
	}
	return scheme
}

func GetHost(r *http.Request) string {
	host := r.Host
	if v := r.Header.Get("X-Forwarded-Host"); v != "" {
		host = v
	}
	return host
}

func GetPath(r *http.Request) string {
	path := r.URL.Path
	if v := r.Header.Get("X-Forwarded-Prefix"); v != "" {
		path = v + path
	}
	return path
}
