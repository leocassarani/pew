package probe

import (
	"fmt"
	"os"
	"time"
)

type Probe struct {
	stop    chan struct{}
	process *ProcessProbe
	Memory  *MemoryUsage
}

func New(proc *os.Process) (*Probe, error) {
	process, err := NewProcessProbe(proc)
	if err != nil {
		return nil, err
	}

	return &Probe{
		process: process,
		stop:    make(chan struct{}, 1),
		Memory:  &MemoryUsage{},
	}, nil
}

func (p *Probe) Start(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			err := p.sample()
			if err != nil {
				fmt.Fprintf(os.Stderr, "pew: %v\n", err)
			}
		case <-p.stop:
			ticker.Stop()
			break
		}
	}
}

func (p *Probe) sample() error {
	pStat, err := p.process.Sample()
	if err != nil {
		return err
	}

	p.Memory.Record(pStat)
	return nil
}

func (p *Probe) Stop() {
	p.stop <- struct{}{}
}

func (p *Probe) Close() {
	p.process.Close()
}
