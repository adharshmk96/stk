package gsk

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Params interface {
	ByName(string) string
}

type Router interface {
	ServeFiles(string, http.FileSystem)
	ServeHTTP(http.ResponseWriter, *http.Request)
	HandlerFunc(method string, path string, handler http.HandlerFunc)
	ParamsFromContext(context.Context) Params

	Router() *httprouter.Router
}

type gskRouter struct {
	router *httprouter.Router
}

func (gr *gskRouter) ServeFiles(path string, fs http.FileSystem) {
	gr.router.ServeFiles(path, fs)
}

func (gr *gskRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gr.router.ServeHTTP(w, r)
}

func (gr *gskRouter) HandlerFunc(method string, path string, handler http.HandlerFunc) {
	gr.router.HandlerFunc(method, path, handler)
}

func (gr *gskRouter) ParamsFromContext(ctx context.Context) Params {
	return httprouter.ParamsFromContext(ctx)
}

func (gr *gskRouter) Router() *httprouter.Router {
	return gr.router
}

func newGskRouter() Router {
	router := httprouter.New()
	return &gskRouter{
		router: router,
	}
}
