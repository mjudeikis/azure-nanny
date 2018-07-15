package main

import (
	"github.com/mjudeikis/azure-nanny/pkg/nanny"
)

func main() {

	if err := nanny.Run(); err != nil {
		panic(err)
	}
}
