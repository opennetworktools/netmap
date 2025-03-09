package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "netmap",
	Short: "Netmap stands for Network Mapper, a visualizer for your inventory of network devices.",
	Long: `Netmap uses LLDP information to map devices. Built with love by Roopesh and friends in Go.
	Complete documentation is available at https://github.com/opennetworktools/netmap`,
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
