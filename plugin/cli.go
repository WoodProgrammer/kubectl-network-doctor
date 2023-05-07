package main

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
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
	Use:     "mode",
	Short:   "Mode of the checks in cluster",
	Aliases: []string{"mode"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := []string{"./main"}
		clientset := configHandler()
		if args[0] == "dns" {

			data := generateHostsFile("hosts.txt")

			createConfigMap("dns-test-configmap", "kube-system", data, clientset)
			createDeployment("dns", "emirozbir/dns-func-test:0.0.1", command, "kube-system", clientset)
			fmt.Println("INFO:::Waiting for the results of the logs until the deployment get ready approx:: 50 sec")
			bar := progressbar.Default(50)
			for i := 0; i < 50; i++ {
				bar.Add(1)
				time.Sleep(1 * time.Second)
			}
			gatherLogs("dns-test", "kube-system", clientset)
			deleteDeployment("dns-deployment", "kube-system", clientset)
			deleteConfigMap("dns-test-configmap", "kube-system", clientset)
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
