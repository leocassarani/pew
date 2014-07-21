package main

import (
	"fmt"
	"github.com/leocassarani/pew/command"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		printUsage(args)
		os.Exit(1)
	}

	cmd := args[1]
	cmdArgs := args[2:]

	switch cmd {
	case "run":
		err := command.Run(cmdArgs)
		exit(err)
	case "help":
		fallthrough
	case "--help":
		printUsage(args)
		os.Exit(0)
	}

}

func exit(err error) {
	if err != nil {
		log(err)
		os.Exit(1)
	}
	os.Exit(0)
}

func log(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
}

func printUsage(args []string) {
	fmt.Printf("usage: %s <command>\n", args[0])
}
