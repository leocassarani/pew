package probe

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const RssIndex = 23

type Memory struct {
	Readings []MemoryReading
	stat     *os.File
	pagesize int
}

type MemoryReading struct {
	time time.Time
	rss  int
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

	reading := MemoryReading{
		time: time.Now(),
		rss:  rss * m.pagesize,
	}
	m.Readings = append(m.Readings, reading)

	return nil
}

func (m *Memory) Close() {
	m.stat.Close()
}
