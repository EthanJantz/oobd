package main

import "log"
import "fmt"
import "github.com/ethanjantz/oobd/recurser"

func main() {
	fmt.Println("Hello world!")
	recurser.Test()
	recursers, err := recurser.List()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(recursers)
}
