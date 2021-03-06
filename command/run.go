package command

import (
	"fmt"
	"github.com/leocassarani/pew/process"
	"github.com/leocassarani/pew/profile"
	"github.com/leocassarani/pew/store"
	"os"
	"os/signal"
	"time"
)

func Run(args []string) error {
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	store := store.FileStore{Root: pwd}
	file, err := store.Create(args[0])
	if err != nil {
		return err
	}
	defer file.Close()

	writer := profile.NewWriter(file)

	// Catch SIGINT signals on the interrupt channel.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	runner := process.NewRunner(args)
	err = runner.Start()
	if err != nil {
		return err
	}
	go runner.Wait()

	probe, err := process.NewMemoryProbe(runner.Process())
	if err != nil {
		return err
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
			break loop
		case <-interrupt:
			// Handle a Ctrl-C by shutting down the child process.
			err = runner.Stop()
			break loop
		}
	}

	probe.Stop()
	return err
}

func log(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
}
