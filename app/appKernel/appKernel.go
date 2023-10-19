package appKernel

import (
	"net/http"

	"github.com/waponix/netgo/router"
	"github.com/waponix/netgo/src/product"
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
		RegisterGroup(
			"/api",
			router.Get("/product/{productId}", product.GetProductHandler),
		)

	http.ListenAndServe(":8080", router.Instance().Mux())
}
