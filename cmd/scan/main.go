package main

import (
	"github.com/ethanjantz/oobd/pkg/scan"
)

func main() {
	scan.Dir("/home/", scan.DefaultOpts())
}
