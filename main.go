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

	go poll(cmd.Process)

	err = cmd.Wait()
	if err != nil {
		exit(err)
	}
}

func poll(process *os.Process) {
	mem, err := probe.NewMemory(process)
	if err != nil {
		exit(err)
	}
	defer mem.Close()

	for {
		err = mem.Probe()
		if err != nil {
			exit(err)
		}

		// Poll at 1-second intervals.
		<-time.After(1 * time.Second)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
	os.Exit(1)
}
