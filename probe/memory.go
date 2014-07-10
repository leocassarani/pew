package probe

import (
	"encoding/csv"
	"fmt"
	"io"
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
	stop     chan struct{}
	pagesize int
}

type MemoryReading struct {
	Time time.Time
	RSS  int
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
		stop:     make(chan struct{}, 1),
		pagesize: os.Getpagesize(),
	}, nil
}

func (m *Memory) Probe(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			err := m.takeReading()
			if err != nil {
				fmt.Fprintf(os.Stderr, "pew: %v\n", err)
			}
		case <-m.stop:
			ticker.Stop()
			break
		}
	}
}

func (m *Memory) takeReading() error {
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
		Time: time.Now(),
		RSS:  rss * m.pagesize,
	}
	m.Readings = append(m.Readings, reading)

	return nil
}

func (m *Memory) Stop() {
	m.stop <- struct{}{}
}

func (m *Memory) Close() {
	m.stat.Close()
}

func (m *Memory) WriteTo(w io.Writer) error {
	out := csv.NewWriter(w)
	for _, reading := range m.Readings {
		out.Write(reading.CSVRow())
	}

	out.Flush()
	if err := out.Error(); err != nil {
		return err
	}

	return nil
}

func (r MemoryReading) CSVRow() []string {
	timestamp := fmt.Sprintf("%v", r.Time.Unix())
	rss := strconv.Itoa(r.RSS)
	return []string{timestamp, rss}
}
