package cmd

import (
	"opennetworktools/netmap/internal"

	"github.com/spf13/cobra"
)

var Hostname, Username, Password string

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.Flags().StringVarP(&Hostname, "hostname", "n", "", "hostname to connect")
	createCmd.Flags().StringVarP(&Username, "username", "u", "", "username to connect to the host")
	createCmd.Flags().StringVarP(&Password, "password", "p", "", "password to connect to the host")
}

var createCmd = &cobra.Command{
	Use:                   "create",
	Short:                 "create command is used to create a topology diagram",
	Long:                  `create command is used to create a topology diagram by passing flags`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("hostname")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		internal.Traverse(host, username, password)
	},
}
