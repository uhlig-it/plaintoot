package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/uhlig-it/plaintoot/plaintoot"
	"github.com/uhlig-it/plaintoot/server"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves a plain-text representation of a single Mastoton post via HTTP",
	RunE: func(command *cobra.Command, args []string) error {
		port, found := os.LookupEnv("PORT")

		if !found {
			port = "8080"
		}

		repo, err := plaintoot.NewRepository(command.Context())

		if err != nil {
			return err
		}

		server := server.NewServer(repo).WithBlurb(command.Short)

		maxUpTimeStr, found := os.LookupEnv("MAX_UPTIME")

		if found {
			m, err := time.ParseDuration(maxUpTimeStr)

			if err != nil {
				return err
			}

			server = server.WithMaxUptime(m)
		}

		startupDelayStr, found := os.LookupEnv("STARTUP_DELAY")

		if found {
			m, err := time.ParseDuration(startupDelayStr)

			if err != nil {
				return err
			}

			server = server.WithStartupDelay(m)
		}

		return server.Start(":" + port)
	},
}
