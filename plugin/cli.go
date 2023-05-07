package main

import (
	"fmt"
	"os"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl nd",
	Short: "kubectl nd - a simple Kubernetes plugin to gather cluster network dump",
	Long:  `kubectl nd rocking yeaaa :)) `,
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
		trCommand := []string{"./main.sh"}
		clientset := configHandler()
		data := generateHostsFile("hosts.txt")

		if args[0] == "dns" {
			getDeployment("coredns", "kube-system", clientset)

			fmt.Println("INFO:: Gathering CoreDNS Logs")
			gatherLogs("k8s-app", "kube-dns", "kube-system", clientset) //coredns Logs

			createConfigMap("dns-test-configmap", "kube-system", data, clientset)
			createDeployment("dns", "emirozbir/dns-func-test:0.0.1", command, "kube-system", "/app", clientset)

			fmt.Println("INFO:::Waiting for the results of the logs until the DnsTester deployment get ready approx:: 50 sec")
			bar := progressbar.Default(50)
			for i := 0; i < 50; i++ {
				bar.Add(1)
				time.Sleep(1 * time.Second)
			}

			gatherLogs("app", "dns-test", "kube-system", clientset) // dns-test-logs

			deleteDeployment("dns-deployment", "kube-system", "dns-stack", clientset)
			deleteConfigMap("dns-test-configmap", "kube-system", clientset)
		} else if args[0] == "traceroute" {

			createConfigMap("traceroute-test-configmap", "kube-system", data, clientset)
			createDeployment("traceroute", "emirozbir/traceroute-test:0.0.1", trCommand, "kube-system", "/opt/traceroute", clientset)
			fmt.Println("INFO:::Waiting for the results of the logs until the TraceRoute deployment get ready approx:: 50 sec")
			bar := progressbar.Default(50)
			for i := 0; i < 50; i++ {
				bar.Add(1)
				time.Sleep(1 * time.Second)
			}
			gatherLogs("app", "traceroute-test", "kube-system", clientset)
			deleteDeployment("traceroute-deployment", "kube-system", "traceroute-stack", clientset)
			deleteConfigMap("traceroute-test-configmap", "kube-system", clientset)

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
