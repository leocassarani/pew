package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	command := os.Args[1:]
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		exit(err)
	}
	err = cmd.Wait()
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
	os.Exit(1)
}
