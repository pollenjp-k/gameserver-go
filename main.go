package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/pollenjp/gameserver-go/api"
	"github.com/pollenjp/gameserver-go/api/config"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d : %v", cfg.Port, err)
	}

	// debug print

	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with %s", url)

	// run server

	mux, cleanup, err := api.NewMux(ctx, cfg)
	if err != nil {
		return err
	}
	defer cleanup()

	s := api.NewServer(l, mux)
	return s.Run(ctx)
}
