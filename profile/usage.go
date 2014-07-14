package profile

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

var (
	CsvHeader = []string{"Timestamp", "Metric", "Value"}
)

type Usage struct {
	readings []Reading
}

func (u *Usage) Record(metric string, value int) {
	reading := Reading{
		Time:   time.Now(),
		Metric: metric,
		Value:  value,
	}
	u.readings = append(u.readings, reading)
}

func (u *Usage) WriteTo(w io.Writer) error {
	out := csv.NewWriter(w)
	out.Write(CsvHeader)
	for _, reading := range u.readings {
		out.Write(reading.CSVRow())
	}

	out.Flush()
	if err := out.Error(); err != nil {
		return err
	}

	return nil
}

type Reading struct {
	Time   time.Time
	Metric string
	Value  int
}

func (r Reading) CSVRow() []string {
	timestamp := fmt.Sprintf("%v", r.Time.Unix())
	value := strconv.Itoa(r.Value)
	return []string{timestamp, r.Metric, value}
}
