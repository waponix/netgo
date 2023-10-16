package router

import "strings"

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Router struct {
	Routes []RouteInterface
}

func (r *Router) Register(routers ...RouteInterface) *Router {
	r.Routes = routers
	return r
}

// register a group of routes by defining the group's base path first
func (r *Router) RegisterGroup(p string, rts ...RouteInterface) *Router {
	for _, rt := range rts {
		p1 := strings.Split(p, "/")
		p2 := strings.Split(rt.Path(), "/")

		// remove empty item
		if p2[0] == "" {
			p2 = p2[1:]
		}

		// join the groups path to the route path
		rt.SetPath(strings.Join(p1, "/") + "/" + strings.Join(p2, "/"))
		r.Routes = append(r.Routes, rt)
	}

	return r
}

type RouteInterface interface {
	Method() []string
	SetPath(string) RouteInterface
	Path() string
	Responder() ResponderFunc
	Middlewares() []MiddlewareFunc
}

type route struct {
	method      []string
	path        string
	responder   ResponderFunc
	middlewares []MiddlewareFunc
}

func (r *route) Method() []string {
	return r.method
}

func (r *route) SetPath(p string) RouteInterface {
	r.path = p
	return r
}

func (r *route) Path() string {
	return r.path
}

func (r *route) Responder() ResponderFunc {
	return r.responder
}

func (r *route) Middlewares() []MiddlewareFunc {
	return r.middlewares
}

func Route(m []string, p string, r ResponderFunc, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      m,
		path:        p,
		responder:   r,
		middlewares: ms,
	}
}

func Get(p string, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{GET},
		path:        p,
		middlewares: ms,
	}
}

func Post(p string, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{POST},
		path:        p,
		middlewares: ms,
	}
}

func Put(p string, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{PUT},
		path:        p,
		middlewares: ms,
	}
}

func Delete(p string, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{DELETE},
		path:        p,
		middlewares: ms,
	}
}

type MiddlewareFunc func() bool
type ResponderFunc func()
