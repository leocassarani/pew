package profile

import (
	"encoding/csv"
	"github.com/leocassarani/pew/process"
	"io"
	"strconv"
	"time"
)

var (
	CsvHeader = []string{"Timestamp", "RSS"}
)

type Writer struct {
	csv *csv.Writer
}

func NewWriter(wr io.Writer) *Writer {
	csvWriter := csv.NewWriter(wr)
	csvWriter.Write(CsvHeader)
	return &Writer{csv: csvWriter}
}

func (w *Writer) Write(sample process.MemorySample) error {
	row := csvRow(sample)
	w.csv.Write(row)

	// Write to disk immediately.
	w.csv.Flush()

	return w.csv.Error()
}

func csvRow(sample process.MemorySample) []string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	rss := strconv.Itoa(sample.RSS)
	return []string{timestamp, rss}
}
