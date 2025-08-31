package main

import (
	"fmt"
	"log"

	"github.com/ethanjantz/oobd/pkg/scan"
)

func main() {
	r, err := scan.Dir("/home/rafl", scan.DefaultOpts())
	if err != nil {
		log.Fatalln(err)
	}

	for p, m := range r {
		fmt.Printf("%s %d\n", p, m.Size)
	}
}
