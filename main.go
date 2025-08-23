package main

import (
	"fmt"
	"log"

	"github.com/ethanjantz/oobd/rcapi"
	"github.com/ethanjantz/oobd/recurser"
)

func main() {
	recurser.Test()
	recursers, err := recurser.List()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("%s", recursers)

	for r := range recursers {
		fmt.Printf("Recurser: %+v", r)
		rcID := 100
		InBatch, err := rcapi.IsInBatch(uint32(rcID))
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(InBatch)
	}
}
