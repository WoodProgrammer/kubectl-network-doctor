package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func checkHostsFile() {
	if _, err := os.Stat("hosts.txt"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("hosts.txt file is not exist creating .... ")
		fmt.Println("This file contains sample DNS records to track internal and external DNS resolution time period")

		f, err := os.Create("hosts.txt")

		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		_, err2 := f.WriteString("www.youtube.com\nwww.google.com\n")

		if err2 != nil {
			log.Fatal(err2)
		}
	}

}

func main() {

	fmt.Println("kubectl network doctor 0.0.1")
	pluginScript := "plugin/main.sh"
	checkHostsFile()

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
