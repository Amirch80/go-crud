package routes

import (
	"crud-task/internal/container"
	"net/http"
)

type RouteGroup struct {
	mux    *http.ServeMux
	prefix string
}

func newRouteGroup(mux *http.ServeMux, prefix string) *RouteGroup {
	return &RouteGroup{
		mux:    mux,
		prefix: prefix,
	}
}

func (group *RouteGroup) handle(method, path string, handler http.HandlerFunc) {
	group.mux.HandleFunc(method+" "+group.prefix+path, handler)
}

func (group *RouteGroup) GET(path string, handler http.HandlerFunc) {
	group.handle(http.MethodGet, path, handler)
}

func (group *RouteGroup) POST(path string, handler http.HandlerFunc) {
	group.handle(http.MethodPost, path, handler)
}

func (group *RouteGroup) PUT(path string, handler http.HandlerFunc) {
	group.handle(http.MethodPut, path, handler)
}

func (group *RouteGroup) DELETE(path string, handler http.HandlerFunc) {
	group.handle(http.MethodDelete, path, handler)
}

func Routes(container *container.Container) *http.ServeMux {
	multiplexer := http.NewServeMux()

	api := newRouteGroup(multiplexer, "/api")
	api.GET("/tasks", container.TaskHandler.All)
	api.GET("/tasks/{id}", container.TaskHandler.Show)
	api.POST("/tasks", container.TaskHandler.Create)
	api.PUT("/tasks/{id}", container.TaskHandler.Update)
	api.DELETE("/tasks/{id}", container.TaskHandler.Delete)

	return multiplexer
}
