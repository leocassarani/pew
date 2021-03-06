package command

import (
	"github.com/leocassarani/pew/process"
	"github.com/leocassarani/pew/process/linux"
	"github.com/leocassarani/pew/profile"
	"github.com/leocassarani/pew/store"
	"os"
	"os/signal"
	"strconv"
	"time"
)

func Attach(pidStr string) error {
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	monitor, err := linux.NewProcStatMonitor(proc)
	if err != nil {
		return err
	}

	// Fetch a sample of the /proc/pid/stat file so we can
	// extract the filename of the process's executable and
	// use it to name the output file for its memory profile.
	stat, err := monitor.Sample()
	if err != nil {
		return err
	}

	store := store.FileStore{Root: pwd}
	file, err := store.Create(stat.Command)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := profile.NewWriter(file)

	// Catch SIGINT signals on the interrupt channel.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	probe, err := process.NewMemoryProbe(proc)
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
				break loop
			}
		case <-probe.Errors:
			// If we fail to read a sample, assume the process has terminated.
			break loop
		case <-interrupt:
			break loop
		}
	}

	probe.Stop()
	return err
}
