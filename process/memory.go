package process

import (
	"github.com/leocassarani/pew/process/linux"
	"github.com/leocassarani/pew/profile"
	"time"
)

type MemoryProbe struct {
	pstat *linux.ProcStatMonitor
	usage *profile.Usage
	stop  chan struct{}
}

func NewMemoryProbe(runner *Runner, usage *profile.Usage) (*MemoryProbe, error) {
	pstat, err := linux.NewProcStatMonitor(runner.process())
	if err != nil {
		return nil, err
	}

	return &MemoryProbe{
		pstat: pstat,
		usage: usage,
		stop: make(chan struct{}, 1),
	}, nil
}

func (m *MemoryProbe) SampleEvery(d time.Duration) {
	// Start immediately.
	m.Sample()

	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			m.Sample()
		case <-m.stop:
			ticker.Stop()
			break
		}
	}
}

func (m *MemoryProbe) Sample() error {
	stat, err := m.pstat.Sample()
	if err != nil {
		return err
	}

	m.usage.Record("RSS", stat.RSS)
	return nil
}

func (m *MemoryProbe) Stop() {
	m.stop <- struct{}{}
}

func (m *MemoryProbe) Close() {
	m.pstat.Close()
}
