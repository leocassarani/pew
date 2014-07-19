package main

import (
	"fmt"
	"github.com/leocassarani/pew/process"
	"github.com/leocassarani/pew/profile"
	"github.com/leocassarani/pew/store"
	"os"
	"os/signal"
	"time"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		printUsage(args)
		os.Exit(1)
	}

	if args[1] == "--help" {
		printUsage(args)
		os.Exit(0)
	}

	// Catch SIGINT signals on the interrupt channel.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	file, err := store.Create()
	if err != nil {
		exit(err)
	}
	defer file.Close()

	writer := profile.NewWriter(file)

	command := args[1:]
	runner := process.NewRunner(command)
	err = runner.Start()
	if err != nil {
		exit(err)
	}
	go runner.Wait()

	probe, err := process.NewMemoryProbe(runner)
	if err != nil {
		exit(err)
	}
	defer probe.Close()
	go probe.SampleEvery(1 * time.Second)

loop:
	for {
		select {
		case sample := <-probe.Samples:
			err = writer.Write(sample)
			if err != nil {
				log(err)
			}
		case err = <-runner.Exit:
			if err != nil {
				exit(err)
			}
			break loop
		case <-interrupt:
			// Handle a Ctrl-C by shutting down the child process.
			err = runner.Stop()
			if err != nil {
				exit(err)
			}
			break loop
		}
	}

	probe.Stop()
}

func exit(err error) {
	log(err)
	os.Exit(1)
}

func log(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
}

func printUsage(args []string) {
	fmt.Printf("usage: %s <command>\n", args[0])
}
