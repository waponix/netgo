package router

import (
	"net/http"
	"strings"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

// ===== STARTOF Router =====
type RouterInterface interface {
	Register(...RouterInterface) *router
	RegisterGroup(string, ...RouterInterface) *router
}

type router struct {
	Routes RoutesMap
}

var routerInstance *router

func Instance() *router {
	if routerInstance == nil {
		routerInstance = &router{
			Routes: make(RoutesMap),
		}
	}

	return routerInstance
}

// register routes
func (r *router) Register(routers ...RouteInterface) *router {
	return r.register(routers)
}

func (r *router) register(routers []RouteInterface) *router {

	for _, rt := range routers {
		_, ok := r.Routes[rt.Path()]

		if ok {
			// should have a mechanism to determine which middlewares and handler to use depending on the method
			// this will ignore for now
			continue
		}

		r.Routes[rt.Path()] = rt
	}

	return r
}

// register a group of routes by defining the group's base path first
func (r *router) RegisterGroup(p string, rts ...RouteInterface) *router {
	for _, rt := range rts {
		p1 := strings.Split(p, "/")
		p2 := strings.Split(rt.Path(), "/")

		// remove empty item
		if p2[0] == "" {
			p2 = p2[1:]
		}

		// join the groups path to the route path
		rt.SetPath(strings.Join(p1, "/") + "/" + strings.Join(p2, "/"))
	}

	return r.register(rts)
}

func (r *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range r.Routes {
		handler := route.Apply()

		mux.Handle(route.Path(), handler)
	}

	return mux
}

// ===== ENDOF Router =====

// ===== STARTOF Route =====
type RouteInterface interface {
	Method() []string
	SetPath(string) RouteInterface
	Path() string
	Middlewares() []MiddlewareFunc
	Apply() http.Handler
}

type route struct {
	method      []string
	path        string
	handler     http.HandlerFunc
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

func (r *route) Middlewares() []MiddlewareFunc {
	return r.middlewares
}

func (r *route) Apply() http.Handler {
	if len(r.middlewares) <= 0 {
		return http.Handler(r.handler)
	}

	handler := http.Handler(r.handler)
	for _, m := range r.middlewares {
		m := middleware(m)
		handler = m(handler)
	}

	return handler
}

func middleware(m MiddlewareFunc) RouteMiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ok := m(r, next)

			if ok {
				next.ServeHTTP(w, r)
			}
		})
	}
}

func Route(m []string, p string, r http.HandlerFunc, ms ...MiddlewareFunc) RouteInterface {

	return &route{
		method:      m,
		path:        p,
		handler:     r,
		middlewares: ms,
	}
}

func Get(p string, r http.HandlerFunc, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{GET},
		path:        p,
		handler:     r,
		middlewares: ms,
	}
}

func Post(p string, r http.HandlerFunc, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{POST},
		path:        p,
		handler:     r,
		middlewares: ms,
	}
}

func Put(p string, r http.HandlerFunc, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{PUT},
		path:        p,
		handler:     r,
		middlewares: ms,
	}
}

func Delete(p string, r http.HandlerFunc, ms ...MiddlewareFunc) RouteInterface {
	return &route{
		method:      []string{DELETE},
		path:        p,
		handler:     r,
		middlewares: ms,
	}
}

// ===== ENDOF Route =====

// ===== TYPES =====

type RouteMiddlewareFunc func(http.Handler) http.Handler
type MiddlewareFunc func(*http.Request, http.Handler) bool
type RoutesMap map[string]RouteInterface
