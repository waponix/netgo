package appKernel

import (
	"fmt"

	"github.com/waponix/netgo/router"
)

type Kernel struct {
}

func New() *Kernel {
	return &Kernel{}
}

func (k Kernel) Init() {
	r := router.Router{}
	r.
		Register(
			router.Get("/get"),
			router.Post("/post"),
			router.Put("/put"),
			router.Delete("/delete"),
		).
		RegisterGroup(
			"/api",
			router.Get("/get"),
			router.Post("/post"),
			router.Put("/put"),
			router.Delete("/delete"),
		)

	for _, rt := range r.Routes {
		fmt.Println(rt)
	}
}
