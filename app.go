package main

import (
	"github.com/waponix/netgo/app/appKernel"
)

func main() {
	kernel := appKernel.New()

	kernel.Init()
}
