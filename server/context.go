package server

import (
	"context"
	"log"
	"net/http"
	"net/url"
)

type contextKey int

const (
	ctxKeyEndpointContext contextKey = iota
)

type endpointContext struct {
	Sort   string
	Host   string
	Scheme string
	Query  url.Values
}

func withEndpointContext(parent context.Context, r *http.Request) context.Context {
	scheme := r.URL.Scheme
	if scheme == "" {
		scheme = "http"
	}
	epCtx := &endpointContext{
		Sort:   r.URL.Query().Get("sort"),
		Host:   r.Host,
		Scheme: scheme,
		Query:  r.URL.Query(),
	}
	log.Printf("r: %#v", r.URL.Hostname())
	return context.WithValue(parent, ctxKeyEndpointContext, epCtx)
}

func getEndpointContext(ctx context.Context) (epCtx *endpointContext) {
	epCtx, _ = ctx.Value(ctxKeyEndpointContext).(*endpointContext)
	return
}
