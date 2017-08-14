package api

import (
	"context"
	"net/http"
	"net/url"
)

type contextKey int

const (
	ctxKeyEndpointContext contextKey = iota
	ctxKeyFS
	ctxKeyGraphVariables
)

type endpointContext struct {
	Sort   string
	Host   string
	Scheme string
	Query  url.Values
	FS     http.FileSystem
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
	return context.WithValue(parent, ctxKeyEndpointContext, epCtx)
}

func getEndpointContext(ctx context.Context) (epCtx *endpointContext) {
	epCtx, _ = ctx.Value(ctxKeyEndpointContext).(*endpointContext)
	return
}

func withFilesystem(parent context.Context, fs http.FileSystem) context.Context {
	return context.WithValue(parent, ctxKeyFS, fs)
}

func getFilesystem(ctx context.Context) (fs http.FileSystem) {
	fs, _ = ctx.Value(ctxKeyFS).(http.FileSystem)
	return
}

func withGraphArgs(parent context.Context, args map[string]interface{}) context.Context {
	return context.WithValue(parent, ctxKeyGraphVariables, args)
}

func getGraphArgs(ctx context.Context) (args map[string]interface{}) {
	args, _ = ctx.Value(ctxKeyGraphVariables).(map[string]interface{})
	return
}
