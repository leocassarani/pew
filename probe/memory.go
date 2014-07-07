package probe

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
)

const RssIndex = 23

type Memory struct {
	Readings []int
	stat     *os.File
	pagesize int
}

func NewMemory(process *os.Process) (*Memory, error) {
	pid := strconv.Itoa(process.Pid)
	filepath := path.Join("/proc", pid, "stat")

	stat, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	return &Memory{
		stat:     stat,
		pagesize: os.Getpagesize(),
	}, nil
}

func (m *Memory) Probe() error {
	_, err := m.stat.Seek(0, 0)
	if err != nil {
		return err
	}

	text, err := ioutil.ReadAll(m.stat)
	if err != nil {
		return err
	}

	stats := strings.Split(string(text), " ")
	rss, err := strconv.Atoi(stats[RssIndex])
	if err != nil {
		return err
	}

	m.Readings = append(m.Readings, rss*m.pagesize)
	return nil
}

func (m *Memory) Close() {
	m.stat.Close()
}
