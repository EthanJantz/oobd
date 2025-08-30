package main

import (
	"fmt"
	"log"

	"github.com/ethanjantz/oobd/rcapi"
	"github.com/ethanjantz/oobd/recurser"
)

// TODO: config file, or just trust there's no spurious 404s from that API endpoint?
// TODO: eventually: remove these users if they're really deactivated
// RC IDs, not system user IDs
var skip = map[uint32]struct{}{
	2186: struct{}{},
	2588: struct{}{},
	1342: struct{}{},
	124: struct{}{},
	4453: struct{}{},
	5127: struct{}{},
	5809: struct{}{},
	6284: struct{}{},
}

func main() {
	recursers, err := recurser.List()
	if err != nil {
		log.Fatalln(err)
	}

	for _, r := range recursers {
		rcId := r.RcId()
		if _, ok := skip[rcId]; ok {
			continue
		}
		fmt.Printf("Recurser: %+v\n", r)
		InBatch, err := rcapi.IsInBatch(rcId)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(InBatch)
	}
}
