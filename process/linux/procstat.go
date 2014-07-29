package linux

import (
	"bufio"
	"os"
	"path"
	"strconv"
)

const (
	CommIndex = 1
	RssIndex  = 23
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

func (p *ProcStatMonitor) Sample() (ProcStat, error) {
	stat := ProcStat{Pagesize: p.pagesize}

	_, err := p.file.Seek(0, 0)
	if err != nil {
		return stat, err
	}

	scanner := bufio.NewScanner(p.file)
	scanner.Split(bufio.ScanWords)

	for i := 0; ; i++ {
		ok := scanner.Scan()
		if !ok {
			break
		}
		stat.setValueAtIndex(scanner.Text(), i)
	}

	if err = scanner.Err(); err != nil {
		return stat, err
	}

	return stat, err
}

func (p *ProcStatMonitor) Close() {
	p.file.Close()
}

type ProcStat struct {
	// The filename of the process's executable.
	Command string

	// Resident Set Size of a process.
	RSS int

	// Size of a memory page on the host.
	Pagesize int
}

func (s *ProcStat) setValueAtIndex(value string, idx int) error {
	switch idx {
	case CommIndex:
		// The command will be wrapped in parentheses, so we remove
		// the first and last character to extract the bare filename.
		s.Command = value[1:len(value)-1]
	case RssIndex:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		s.RSS = intValue
	}
	return nil
}
