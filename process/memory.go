package process

import (
	"github.com/leocassarani/pew/process/linux"
	"os"
	"time"
)

type MemorySample struct {
	RSS int
}

type MemoryProbe struct {
	Samples chan MemorySample // The channel on which process samples are delivered.
	Errors  chan error        // The channel on which process sampling errors are delivered.

	pstat *linux.ProcStatMonitor
	stop  chan struct{}
}

func NewMemoryProbe(process *os.Process) (*MemoryProbe, error) {
	pstat, err := linux.NewProcStatMonitor(process)
	if err != nil {
		return nil, err
	}

	return &MemoryProbe{
		Samples: make(chan MemorySample),
		Errors:  make(chan error),
		stop:    make(chan struct{}, 1),
		pstat:   pstat,
	}, nil
}

func (m *MemoryProbe) SampleEvery(d time.Duration) {
	// Start immediately.
	err := m.Sample()
	if err != nil {
		m.Errors <- err
	}

	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			err = m.Sample()
			if err != nil {
				m.Errors <- err
			}
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

	m.Samples <- MemorySample{RSS: stat.RSS * stat.Pagesize}
	return nil
}

func (m *MemoryProbe) Stop() {
	m.stop <- struct{}{}
}

func (m *MemoryProbe) Close() {
	m.pstat.Close()
}
