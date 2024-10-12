package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
	"github.com/yzimhao/trading_engine/v2/internal/di"
)

func main() {
	_ = godotenv.Load()

	migrateCmd := &cli.Command{
		Name:     "migrate",
		Usage:    "migrate database",
		Commands: []*cli.Command{},
	}

	cmd := &cli.Command{
		Name: "trader",
		Action: func(_ context.Context, cmd *cli.Command) error {
			app := di.App()
			app.Run()
			return nil
		},
		Commands: []*cli.Command{
			migrateCmd,
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
