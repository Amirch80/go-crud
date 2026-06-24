package routes

import (
	"crud-task/internal/container"
	"net/http"
)

type RouteGroup struct {
	prefix      string
	multiplexer *http.ServeMux
}

func newRouteGroup(multiplexer *http.ServeMux, prefix string) *RouteGroup {
	return &RouteGroup{
		prefix,
		multiplexer,
	}
}

func (routeGroup *RouteGroup) GET(path string, handler http.HandlerFunc) {
	fullPath := routeGroup.prefix + path
	routeGroup.multiplexer.HandleFunc("GET "+fullPath, handler)
}
func (routeGroup *RouteGroup) POST(path string, handler http.HandlerFunc) {
	fullPath := routeGroup.prefix + path
	routeGroup.multiplexer.HandleFunc("POST "+fullPath, handler)
}
func (routeGroup *RouteGroup) PUT(path string, handler http.HandlerFunc) {
	fullPath := routeGroup.prefix + path
	routeGroup.multiplexer.HandleFunc("PUT "+fullPath, handler)
}
func (routeGroup *RouteGroup) DELETE(path string, handler http.HandlerFunc) {
	fullPath := routeGroup.prefix + path
	routeGroup.multiplexer.HandleFunc("DELETE "+fullPath, handler)
}

func Routes(container *container.Container) *http.ServeMux {
	multiplexer := http.NewServeMux()

	apiRouteGroup := newRouteGroup(multiplexer, "/api")
	apiRouteGroup.GET("/all-tasks", container.TaskHandler.All)
	apiRouteGroup.GET("/show-task/{id}", container.TaskHandler.Show)
	apiRouteGroup.POST("/create-task", container.TaskHandler.Create)
	apiRouteGroup.PUT("/update-task/{id}", container.TaskHandler.Update)
	apiRouteGroup.DELETE("/delete-task/{id}", container.TaskHandler.Delete)

	return multiplexer
}
