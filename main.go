package main

import (
	"fmt"
	"github.com/leocassarani/pew/probe"
	"github.com/leocassarani/pew/process"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Catch SIGINT signals on the interrupt channel.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	command := os.Args[1:]

	runner := process.NewRunner(command)
	err := runner.Run()
	if err != nil {
		exit(err)
	}

	probe, err := probe.New(runner.Process())
	if err != nil {
		exit(err)
	}
	go probe.Start(1 * time.Second)
	defer probe.Close()

	select {
	case err = <-runner.Exit:
		if err == nil {
			probe.Stop()
		} else {
			exit(err)
		}
	case <-interrupt:
		// Handle a Ctrl-C by shutting down the child process.
		err = runner.Stop()
		if err != nil {
			exit(err)
		}
	}

	if err = os.MkdirAll(".pew", os.ModePerm); err != nil {
		exit(err)
	}

	csv, err := os.Create(".pew/memory.csv")
	if err != nil {
		exit(err)
	}

	err = probe.Memory.WriteTo(csv)
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
	os.Exit(1)
}
