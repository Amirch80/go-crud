package routes

import (
	"crud-task/internal/container"
	"crud-task/internal/middleware"
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

func Routes(container *container.Container) http.Handler {
	multiplexer := http.NewServeMux()

	api := newRouteGroup(multiplexer, "/api")
	api.GET("/tasks", container.TaskHandler.All)
	api.GET("/tasks/{id}", container.TaskHandler.Show)
	api.POST("/tasks", container.TaskHandler.Create)
	api.PUT("/tasks/{id}", container.TaskHandler.Update)
	api.DELETE("/tasks/{id}", container.TaskHandler.SoftDelete)
	api.DELETE("/tasks/{id}/force", container.TaskHandler.Delete)
	api.GET("/tasks/statuses", container.TaskHandler.StatusOptions)

	return middleware.CORS(multiplexer)
}
