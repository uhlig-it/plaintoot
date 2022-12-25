package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "plaintoot",
	Short: "Provides a plain-text representation of a single post",
}
