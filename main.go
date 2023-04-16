package main

import (
	"fmt"
	"os/exec"
)

func main() {

	fmt.Println("kubectl network doctor 0.0.1")
	pluginScript := "plugin/main.sh"

	cmd := exec.Command(pluginScript)
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		fmt.Println(err.Error())
	}

	if err = cmd.Start(); err != nil {
		fmt.Println(err.Error())

	}

	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
}
