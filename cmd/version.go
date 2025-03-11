package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Netmap",
	Long:  `All software has versions. This is Netmap's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("netmap v0.1.3 -- HEAD")
	},
}
