package process

import (
	"os"
	"os/exec"
)

type Runner struct {
	Exit chan error // The channel on which process exit will be communicated.

	name string
	args []string
	cmd  *exec.Cmd
}

func NewRunner(command []string) *Runner {
	return &Runner{
		Exit: make(chan error, 1),
		name: command[0],
		args: command[1:],
	}
}

func (r *Runner) Start() error {
	r.cmd = exec.Command(r.name, r.args...)
	r.cmd.Stdout = os.Stdout
	r.cmd.Stderr = os.Stderr

	return r.cmd.Start()
}

func (r *Runner) Wait() {
	err := r.cmd.Wait()
	r.Exit <- err
}

func (r *Runner) Process() *os.Process {
	return r.cmd.Process
}

// Stop sends a SIGINT signal to the child process.
func (r *Runner) Stop() error {
	return r.Process().Signal(os.Interrupt)
}
