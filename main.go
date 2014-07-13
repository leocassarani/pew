package main

import (
	"fmt"
	"github.com/leocassarani/pew/process"
	"github.com/leocassarani/pew/profile"
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

	usage := &profile.Usage{}
	mem, err := process.NewMemoryProbe(runner, usage)
	if err != nil {
		exit(err)
	}
	go mem.SampleEvery(1 * time.Second)
	defer mem.Close()

	select {
	case err = <-runner.Exit:
		if err == nil {
			mem.Stop()
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

	err = usage.WriteTo(csv)
	if err != nil {
		exit(err)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
	os.Exit(1)
}
