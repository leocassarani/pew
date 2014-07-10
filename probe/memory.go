package probe

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

type MemoryUsage struct {
	readings []MemoryReading
}

func (m *MemoryUsage) Record(stat ProcessStatus) {
	reading := MemoryReading{
		Time: time.Now(),
		RSS:  stat.RSS,
	}
	m.readings = append(m.readings, reading)
}

func (m *MemoryUsage) WriteTo(w io.Writer) error {
	out := csv.NewWriter(w)
	for _, reading := range m.readings {
		out.Write(reading.CSVRow())
	}

	out.Flush()
	if err := out.Error(); err != nil {
		return err
	}

	return nil
}

type MemoryReading struct {
	Time time.Time
	RSS  int
}

func (r MemoryReading) CSVRow() []string {
	timestamp := fmt.Sprintf("%v", r.Time.Unix())
	rss := strconv.Itoa(r.RSS)
	return []string{timestamp, rss}
}
