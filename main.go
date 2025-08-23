package main

import "log"
import "fmt"
import "github.com/ethanjantz/oobd/recurser"
import "github.com/ethanjantz/oobd/rcapi"

func main() {
	recurser.Test()
	recursers, err := recurser.List()
	if err != nil {
		log.Fatalln(err)
	}

	InBatch, err := rcapi.IsInBatch(100)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("%+v", InBatch)
}
