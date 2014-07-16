package process

import (
	"github.com/leocassarani/pew/process/linux"
	"time"
)

type MemorySample linux.ProcStat

type MemoryProbe struct {
	Samples chan MemorySample // The channel on which process samples are delivered.

	pstat *linux.ProcStatMonitor
	stop  chan struct{}
}

func NewMemoryProbe(runner *Runner) (*MemoryProbe, error) {
	pstat, err := linux.NewProcStatMonitor(runner.process())
	if err != nil {
		return nil, err
	}

	return &MemoryProbe{
		Samples: make(chan MemorySample),
		pstat:   pstat,
		stop:    make(chan struct{}, 1),
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

	m.Samples <- MemorySample{RSS: stat.RSS}
	return nil
}

func (m *MemoryProbe) Stop() {
	m.stop <- struct{}{}
}

func (m *MemoryProbe) Close() {
	m.pstat.Close()
}
