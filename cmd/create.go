package cmd

import (
	// "github.com/opennetworktools/netmap/internal"

	"github.com/opennetworktools/netmap/internal"
	"github.com/spf13/cobra"
)

var Hostname, Username, Password string

func init() {
	rootCmd.AddCommand(createCmd)
	// createCmd.Flags().StringVarP(&Hostname, "hostname", "n", "", "hostname to connect")
	// createCmd.Flags().StringVarP(&Username, "username", "u", "", "username to connect to the host")
	// createCmd.Flags().StringVarP(&Password, "password", "p", "", "password to connect to the host")

	createCmd.AddCommand(lldpCmd)
	createCmd.AddCommand(ospfCmd)

	lldpCmd.Flags().StringVarP(&Hostname, "hostname", "n", "", "hostname to connect")
	lldpCmd.Flags().StringVarP(&Username, "username", "u", "", "username to connect to the host")
	lldpCmd.Flags().StringVarP(&Password, "password", "p", "", "password to connect to the host")

	ospfCmd.Flags().StringVarP(&Hostname, "hostname", "n", "", "hostname to connect")
	ospfCmd.Flags().StringVarP(&Username, "username", "u", "", "username to connect to the host")
	ospfCmd.Flags().StringVarP(&Password, "password", "p", "", "password to connect to the host")
}

var createCmd = &cobra.Command{
	Use:                   "create",
	Short:                 "create command is used to create a topology diagram",
	Long:                  `create command is used to create a topology diagram by passing flags`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var lldpCmd = &cobra.Command{
	Use:                   "lldp",
	Short:                 "create graph using LLDP",
	Long:                  `create graph using LLDP information by passing flags`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("hostname")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		internal.Traverse(host, username, password)
	},
}

var ospfCmd = &cobra.Command{
	Use:                   "ospf",
	Short:                 "create graph using OSPF",
	Long:                  `create graph using OSPF LSDB information by passing flags`,
	DisableFlagsInUseLine: true,
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("hostname")
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		// internal.Traverse(host, username, password)
		internal.Ospf(host, username, password)
	},
}
