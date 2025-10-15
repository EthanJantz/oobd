package main

import (
	"context"
	"expvar"
	"flag"
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
var skip = map[uint32]int8{
	2186: 0,
	2588: 0,
	1342: 0,
	124:  0,
	4453: 0,
	5127: 0,
	5809: 0,
	6284: 0,
}

type config struct {
	port int
	mode string
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
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	flag.StringVar(&cfg.mode, "mode", "manager", "Mode to run in: manager or worker")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger}

	app.logger.Info("starting", "mode", cfg.mode)
	if cfg.mode == "manager" {
		usersNotAtRC, err := scanRecursers()
		if err != nil {
			log.Fatalln(err)
		}
		logger.Log(context.Background(), slog.LevelInfo, "users currently not at rc", "map", usersNotAtRC)
	}

	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func scanRecursers() (map[uint32]recurser.Recurser, error) {
	unixUsers := make(map[uint32]recurser.Recurser)
	recursers, err := recurser.List()
	if err != nil {
		return nil, err
	}

	for _, r := range recursers {
		rcId := r.RcId()
		if _, ok := skip[rcId]; ok {
			continue
		}
		unixUsers[rcId] = r
	}

	rcUsers, err := rcapi.GetRecursers()
	if err != nil {
		return nil, err
	}
	for _, u := range rcUsers {
		if u.CurrentlyAtRc {
			delete(unixUsers, uint32(u.Id))
		}
	}

	return unixUsers, nil
}
