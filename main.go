package main

import (
	"fmt"
	"os/exec"
)

func main() {
	fmt.Println("kubectl network doctor 0.0.1")
	pluginScript := "plugin/main.sh"

	cmd := exec.Command(pluginScript)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(string(stdout))
}
