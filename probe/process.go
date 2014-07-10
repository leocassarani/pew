package probe

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const (
	RssIndex = 23
)

type ProcessProbe struct {
	file *os.File
}

type ProcessStatus struct {
	// Resident Set Size of the process's memory.
	RSS int
}

func NewProcessProbe(proc *os.Process) (*ProcessProbe, error) {
	pid := strconv.Itoa(proc.Pid)
	fpath := path.Join("/", "proc", pid, "stat")

	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	return &ProcessProbe{
		file: file,
	}, nil
}

func (p *ProcessProbe) Sample() (stat ProcessStatus, err error) {
	_, err = p.file.Seek(0, 0)
	if err != nil {
		return stat, err
	}

	text, err := ioutil.ReadAll(p.file)
	if err != nil {
		return stat, err
	}

	fields := strings.Split(string(text), " ")
	rss, err := strconv.Atoi(fields[RssIndex])
	if err != nil {
		return stat, err
	}
	stat.RSS = rss

	return stat, err
}

func (p *ProcessProbe) Close() {
	p.file.Close()
}
