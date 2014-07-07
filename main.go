package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
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
	filepath := path.Join("/proc", strconv.Itoa(process.Pid), "stat")
	file, err := os.Open(filepath)
	if err != nil {
		exit(err)
	}
	defer file.Close()

	for {
		_, err := file.Seek(0, 0)
		if err != nil {
			exit(err)
		}

		text, err := ioutil.ReadAll(file)
		if err != nil {
			exit(err)
		}

		fields := strings.Split(string(text), " ")
		rss, err := strconv.Atoi(fields[23])
		if err != nil {
			exit(err)
		}
		fmt.Println(rss * os.Getpagesize())

		<-time.After(1 * time.Second)
	}
}

func exit(err error) {
	fmt.Fprintf(os.Stderr, "pew: %v\n", err)
	os.Exit(1)
}
