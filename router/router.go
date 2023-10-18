package router

import (
	"net/http"
	"strings"

	"github.com/waponix/netgo/utils/sliceUtil"
)

const (
	GET     = http.MethodGet
	POST    = http.MethodPost
	PUT     = http.MethodPut
	PATCH   = http.MethodPatch
	HEAD    = http.MethodHead
	OPTIONS = http.MethodOptions
	DELETE  = http.MethodDelete
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
func (_router *router) Register(routers ...RouteInterface) *router {
	return _router.register(routers)
}

// register a group of routes by defining the group's base path first
func (_router *router) RegisterGroup(path string, rts ...RouteInterface) *router {
	for _, rt := range rts {
		p1 := strings.Split(path, "/")
		p2 := strings.Split(rt.Path(), "/")

		// remove empty item
		if p2[0] == "" {
			p2 = p2[1:]
		}

		// join the groups path to the route path
		rt.SetPath(strings.Join(p1, "/") + "/" + strings.Join(p2, "/"))
	}

	return _router.register(rts)
}

func (_router *router) register(routers []RouteInterface) *router {
	for _, rt := range routers {
		ert, ok := _router.Routes[rt.Path()]
		if ok {
			// merge all the middlewares
			ert.SetMiddlewares(append(ert.Middlewares(), rt.Middlewares()...))

			ert.SetMethods(append(ert.Methods(), rt.Methods()...))

			for _, method := range rt.Methods() {
				ert.SetHandler(method, rt.Handler())
			}

			_router.Routes[rt.Path()] = ert
		} else {
			_router.Routes[rt.Path()] = rt
		}
	}

	return _router
}

func (_router *router) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	for _, route := range _router.Routes {
		handler := route.Apply()

		mux.Handle(route.Path(), handler)
	}

	return mux
}

// ===== ENDOF Router =====

// ===== STARTOF Route =====
type RouteInterface interface {
	Methods() []string
	SetMethods([]string) RouteInterface
	SetPath(string) RouteInterface
	Path() string
	Middlewares() []middleware
	SetMiddlewares([]middleware) RouteInterface
	Handler() http.HandlerFunc
	Handlers() HandlerMap
	SetHandler(string, http.HandlerFunc) RouteInterface
	Apply() http.Handler
}

type route struct {
	methods     []string
	path        string
	handler     http.HandlerFunc
	handlers    HandlerMap
	middlewares []middleware
}

type middleware struct {
	Methods  []string
	Function MiddlewareFunc
}

func (_route *route) SetMiddlewares(m []middleware) RouteInterface {
	_route.middlewares = m
	return _route
}

func (_route *route) Middlewares() []middleware {
	return _route.middlewares
}

func (_route *route) SetMethods(methods []string) RouteInterface {
	_route.methods = methods
	return _route
}

func (_route *route) Methods() []string {
	return _route.methods
}

func (_route *route) SetPath(p string) RouteInterface {
	_route.path = p
	return _route
}

func (_route *route) Path() string {
	return _route.path
}

func (_route *route) SetHandler(k string, h http.HandlerFunc) RouteInterface {
	if _route.handlers == nil {
		_route.handlers = make(HandlerMap)
	}
	_route.handlers[k] = h
	return _route
}

func (_route *route) Handlers() HandlerMap {
	return _route.handlers
}

func (_route *route) Handler() http.HandlerFunc {
	return _route.handler
}

func (_route *route) Apply() http.Handler {

	// wrap the handler function
	mainHandler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for {
			if isMethodAllowed(req.Method, _route.Methods()) {
				handler, methodOk := _route.Handlers()[req.Method]

				if methodOk {
					handler.ServeHTTP(w, req)
					break
				} else {
					_route.Handler().ServeHTTP(w, req)
					break
				}
			}

			var handler http.Handler

			if _route.Methods() == nil || len(_route.Methods()) <= 0 {
				handler = _route.handler
			}

			if handler != nil {
				handler.ServeHTTP(w, req)
				break
			}

			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			break
		}
	})

	handler := http.Handler(mainHandler)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for {
			stop := false
			for _, _middleware := range _route.middlewares {
				// when the middleware is allowed for the request method, execute it
				if isMethodAllowed(r.Method, _middleware.Methods) && !_middleware.Function(w, r) {
					stop = true // set stop flag to true when a middleware returns false
					break
				}
			}

			if stop {
				// when flag is true do not execute next middlewares and the main handler
				break
			}

			handler.ServeHTTP(w, r)
			break
		}
	})
}

func Route(methods []string, path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  methods,
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     methods,
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{},
		middlewares: middlewares,
	}
}

func Get(path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  []string{GET},
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     []string{GET},
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{GET: handler},
		middlewares: middlewares,
	}
}

func Post(path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	methods := []string{POST}
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  methods,
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     methods,
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{POST: handler},
		middlewares: middlewares,
	}
}

func Put(path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	methods := []string{PUT}
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  methods,
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     methods,
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{PUT: handler},
		middlewares: middlewares,
	}
}

func Patch(path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	methods := []string{PATCH}
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  methods,
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     methods,
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{PATCH: handler},
		middlewares: middlewares,
	}
}

func Head(path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	methods := []string{HEAD}
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  methods,
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     methods,
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{HEAD: handler},
		middlewares: middlewares,
	}
}

func Delete(path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	methods := []string{DELETE}
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  methods,
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     methods,
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{DELETE: handler},
		middlewares: middlewares,
	}
}

func Options(path string, handler http.HandlerFunc, middlewareFuncs ...MiddlewareFunc) RouteInterface {
	methods := []string{OPTIONS}
	middlewares := make([]middleware, len(middlewareFuncs))
	for _, middlewareFunc := range middlewareFuncs {
		middlewares = append(middlewares, middleware{
			Methods:  methods,
			Function: middlewareFunc,
		})
	}

	return &route{
		methods:     []string{OPTIONS},
		path:        path,
		handler:     handler,
		handlers:    HandlerMap{OPTIONS: handler},
		middlewares: middlewares,
	}
}

func toMiddleware(m MiddlewareFunc, router RouteInterface) RouteMiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isMethodAllowed(r.Method, router.Methods()) {
				m(w, r)
			}
		})
	}
}

func isMethodAllowed(requestMethod string, allowedMethods []string) bool {
	return sliceUtil.Use(allowedMethods).InItems(requestMethod)
}

// ===== ENDOF Route =====

// ===== TYPES =====

type RouteMiddlewareFunc func(http.Handler) http.Handler
type MiddlewareFunc func(http.ResponseWriter, *http.Request) bool
type RoutesMap map[string]RouteInterface
type HandlerMap map[string]http.HandlerFunc
