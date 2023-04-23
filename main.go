package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const colorRed = "\033[0;31m"
const colorNone = "\033[0m"
const colorBlue = "\033[1;34m"

func checkHostsFile() {
	if _, err := os.Stat("hosts.txt"); errors.Is(err, os.ErrNotExist) {
		fmt.Println("hosts.txt file is not exist creating .... %s", colorRed)
		fmt.Println("This file contains sample DNS records to track internal and external DNS resolution time period %s", colorRed)

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

	fmt.Printf("%s KUBECTL Network Docker Version number::0.0.1 \n", colorBlue)
	fmt.Println(`
              (
               )
              (
        /\  .-"""-.  /\
       //\\/  ,,,  \//\\
       |/\| ,;;;;;, |/\|
       //\\\;-"""-;///\\
      //  \/   .   \/  \\
     (| ,-_| \ | / |_-, |)
       // __\.-.-./__ \\
      // /.-(() ())-.\ \\
     (\ |)   '---'   (| /)
        (|           |)  
        \)           (/
	`)

	fmt.Println(`
This script basically creates number of pods \n
and run some queries please check WoodProgrammer/kubectl-network-doctor`)

	pluginScript := "/bin/bash plugin/main.sh"
	checkHostsFile()

	cmd := exec.Command(pluginScript)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	/*stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		fmt.Println(err.Error())
	}

	if err = cmd.Run(); err != nil {
		fmt.Println(err.Error())

	}

	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}*/
}
