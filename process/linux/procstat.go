package linux

import (
	"bufio"
	"os"
	"path"
	"strconv"
)

const (
	PidIndex   = 0
	CommIndex  = 1
	StateIndex = 2
	RssIndex   = 23
)

type ProcessState uint8

const (
	StateUnknown ProcessState = iota
	StateRunning
	StateSleeping
	StateDiskSleeping
	StateZombie
	StateTraced
	StatePaging
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
	// The process identifier.
	Pid int

	// The filename of the process's executable.
	Command string

	// The current state of the process.
	State ProcessState

	// Resident Set Size of a process.
	RSS int

	// Size of a memory page on the host.
	Pagesize int
}

func (s *ProcStat) setValueAtIndex(value string, idx int) error {
	switch idx {
	case PidIndex:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		s.Pid = intValue
	case CommIndex:
		// The command will be wrapped in parentheses, so we remove
		// the first and last character to extract the bare filename.
		s.Command = value[1 : len(value)-1]
	case StateIndex:
		s.State = parseProcessState(value)
	case RssIndex:
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		s.RSS = intValue
	}
	return nil
}

func parseProcessState(str string) ProcessState {
	switch str[0] {
	case 'R':
		return StateRunning
	case 'S':
		return StateSleeping
	case 'D':
		return StateDiskSleeping
	case 'Z':
		return StateZombie
	case 'T':
		return StateTraced
	case 'W':
		return StatePaging
	default:
		return StateUnknown
	}
}
