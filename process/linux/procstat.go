package linux

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

type ProcStatMonitor struct {
	file     *os.File
	pagesize int
}

func NewProcStatMonitor(proc *os.Process) (*ProcStatMonitor, error) {
	pid := strconv.Itoa(proc.Pid)
	fpath := path.Join("/", "proc", pid, "stat")

	file, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}

	return &ProcStatMonitor{
		file:     file,
		pagesize: os.Getpagesize(),
	}, nil
}

type ProcStat struct {
	// Resident Set Size of a process.
	RSS int
}

func (p *ProcStatMonitor) Sample() (stat ProcStat, err error) {
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
	stat.RSS = rss * p.pagesize

	return stat, err
}

func (p *ProcStatMonitor) Close() {
	p.file.Close()
}
