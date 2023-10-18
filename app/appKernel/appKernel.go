package appKernel

import (
	"fmt"
	"net/http"

	"github.com/waponix/netgo/router"
)

type Kernel struct {
}

func New() *Kernel {
	return &Kernel{}
}

func TestResponder() {

}

func TestMiddleware() bool {
	return true
}

func (_kernel Kernel) Init() {
	router.Instance().
		Register(
			router.Route([]string{router.GET, router.POST}, "/home", func(w http.ResponseWriter, r *http.Request) {
				// Write the response content
				fmt.Println("Main handler was called")
				fmt.Fprintf(w, "<h1>GET and POST Handler</h1>")

				w.WriteHeader(http.StatusOK)
			}, func(w http.ResponseWriter, r *http.Request) bool {
				fmt.Println("middleware 1")
				return true
			}, func(w http.ResponseWriter, r *http.Request) bool {
				fmt.Println("middleware 2")
				fmt.Fprintf(w, "<h1>Stop right there!</h1>")
				return false
			}, func(w http.ResponseWriter, r *http.Request) bool {
				fmt.Println("middleware 3")
				return true
			}),
			router.Put("/home", func(w http.ResponseWriter, r *http.Request) {
				// Write the response content
				fmt.Println("Main handler was called")
				fmt.Fprintf(w, "<h1>PUT handler</h1>")

				w.WriteHeader(http.StatusOK)
			}),
			router.Delete("/home", func(w http.ResponseWriter, r *http.Request) {
				// Write the response content
				fmt.Println("Main handler was called")
				fmt.Fprintf(w, "<h1>DELETE handler</h1>")

				w.WriteHeader(http.StatusOK)
			}),
		)

	http.ListenAndServe(":8080", router.Instance().Mux())
}
