package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "knd",
	Short: "knd - a simple Kubernetes plugin to gather cluster network dump",
	Long:  `knd rocking yeaaa :)) `,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args[0])
	},
}

var mode = &cobra.Command{
	Use:     "type",
	Short:   "Type of the checks in cluster",
	Aliases: []string{"mode"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//command := []string{"./main"}
		clientset := configHandler()
		if args[0] == "dns" {
			//createDeployment("dns", "emirozbir/dns-func-test:0.0.1", command, "kube-system")
			gatherLogs("dns-test", "kube-system", clientset)
		}
	},
}

func Execute() {
	rootCmd.AddCommand(mode)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
