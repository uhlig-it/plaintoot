package cmd

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/uhlig-it/plaintoot/plaintoot"
)

var PrintCmd = &cobra.Command{
	Use:   "print URL",
	Short: "Prints a plain-text representation of a single Mastodon post",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a post URL")
		}

		_, err := url.Parse(args[0])

		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true // no need to print usage; we'll handle all errors

		uri, err := url.Parse(args[0])

		if err != nil {
			return err
		}

		repo, err := plaintoot.NewRepository(cmd.Context())

		if err != nil {
			return err
		}

		post, err := repo.Lookup(uri)

		if err != nil {
			return err
		}

		fmt.Printf("%s\n", post)

		return nil
	},
}
