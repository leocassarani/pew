package main

import (
	"fmt"
	"github.com/leocassarani/pew/probe"
	"os"
	"os/exec"
	"time"
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

	mem, err := probe.NewMemory(cmd.Process)
	if err != nil {
		exit(err)
	}
	defer mem.Close()

	go mem.Probe(1 * time.Second)

	err = cmd.Wait()
	if err != nil {
		exit(err)
	}

	if err = os.MkdirAll(".pew", os.ModePerm); err != nil {
		exit(err)
	}

	csv, err := os.Create(".pew/memory.csv")
	if err != nil {
		exit(err)
	}

	err = mem.WriteTo(csv)
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
	os.Exit(1)
}
