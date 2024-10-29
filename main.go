package main

import (
	"github.com/urfave/cli/v2"
	app2 "github.com/vadymlab/slot-game/app"
	"github.com/vadymlab/slot-game/internal/config"
	"github.com/vadymlab/slot-game/internal/database"
	"github.com/vadymlab/slot-game/internal/redis"
	"github.com/vadymlab/slot-game/internal/server"
	"github.com/vadymlab/slot-game/internal/utils"
	"log"
	"os"
)

// main is the entry point for the application. It configures and starts the CLI application.
// It sets up flags for configuration and starts the server using app2.RunServer.
func main() {
	// Initialize the CLI application with flags merged from config, database, and server packages.
	app := &cli.App{
		Flags:  utils.MergeSlices(config.LogFlags, database.DatabaseFlags, server.ApiFlags, config.SlotFlags, redis.Flags),
		Action: app2.RunServer,
	}

	// Run the CLI application and handle any errors encountered during execution.
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
