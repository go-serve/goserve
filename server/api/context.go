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
	ctxKeyGraphContext
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

type graphContext struct {
	Args   map[string]interface{}
	Source interface{}
}

func withGraphContext(parent context.Context, graphCtx *graphContext) context.Context {
	return context.WithValue(parent, ctxKeyGraphContext, graphCtx)
}

func getGraphContext(ctx context.Context) (graphCtx *graphContext) {
	graphCtx, _ = ctx.Value(ctxKeyGraphContext).(*graphContext)
	return
}
