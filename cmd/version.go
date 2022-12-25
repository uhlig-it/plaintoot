package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/uhlig-it/plaintoot/plaintoot"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(plaintoot.VersionString())
	},
}
