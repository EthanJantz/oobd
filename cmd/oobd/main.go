package main

import (
	"expvar"
	"flag"
	"fmt"
	"github.com/ethanjantz/oobd/pkg/rcapi"
	"github.com/ethanjantz/oobd/pkg/recurser"
	"log"
	"log/slog"
	"os"
	"strings"
	"time"
)

// TODO: config file, or just trust there's no spurious 404s from that API endpoint?
// TODO: eventually: remove these users if they're really deactivated
// RC IDs, not system user IDs
var skip = map[uint32]struct{}{
	2186: struct{}{},
	2588: struct{}{},
	1342: struct{}{},
	124:  struct{}{},
	4453: struct{}{},
	5127: struct{}{},
	5809: struct{}{},
	6284: struct{}{},
}

type config struct {
	port int
	env  string
	cors struct {
		trustedOrigins []string
	}

	skip map[string]string
}

type application struct {
	config config
	logger *slog.Logger
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

	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}
