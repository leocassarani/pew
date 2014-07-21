package main

import (
	"fmt"
	"github.com/leocassarani/pew/command"
	"os"
)

const (
	Binary = "pew"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cmd := args[1]
	cmdArgs := args[2:]

	switch cmd {
	case "run":
		err := command.Run(cmdArgs)
		exit(err)
	case "attach":
		err := command.Attach(cmdArgs[0])
		exit(err)
	case "help":
		fallthrough
	case "--help":
		printUsage()
		exit(nil)
	default:
		err := fmt.Errorf("'%s' is not a command. See '%s --help'.", cmd, Binary)
		exit(err)
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
	fmt.Fprintf(os.Stderr, "%s: %v\n", Binary, err)
}

func printUsage() {
	fmt.Printf("usage: %s <command>\n", Binary)
}
