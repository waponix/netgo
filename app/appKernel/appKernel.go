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

func (k Kernel) Init() {
	http.ListenAndServe(":8080",
		router.Instance().
			Register(
				router.Get("/home", func(w http.ResponseWriter, r *http.Request) {
					// Write the response content
					fmt.Fprintf(w, "<h1>Hello, World!</h1>")

					// You can also set a specific HTTP status code if needed (e.g., 200 OK)
					w.WriteHeader(http.StatusOK)
				}),
			).
			RegisterGroup(
				"/api",
				router.Get("/report/session", func(w http.ResponseWriter, r *http.Request) {
					fmt.Fprintf(w, "<h1>Session API</h1>")

					w.WriteHeader(http.StatusOK)
				}),
				router.Post("/report/session/id", nil),
			).
			Mux())
}
