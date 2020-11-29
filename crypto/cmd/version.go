package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const Major = "0"
const Minor = "1"
const Fix = "0"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s.%s.%s", Major, Minor, Fix)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
