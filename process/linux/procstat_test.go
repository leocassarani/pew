package linux

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	statFixture = "24984 (godzilla) R 24800 24984 24800 34820 24984 4202496 52072 " +
		"0 0 0 95 46 0 0 20 0 37 0 4605873 7268102144 51579 18446744073709551615 " +
		"4194304 9169328 140733323275456 140733323271144 4446131 0 0 0 2143420159 " +
		"18446744073709551615 0 0 17 1 0 0 0 0 0"
)

func TestProcStatParsing(t *testing.T) {
	statFile := makeStatFile(t, []byte(statFixture))
	defer rmStatFile(statFile)

	monitor := ProcStatMonitor{file: statFile, pagesize: 4096}
	stat, err := monitor.Sample()
	if err != nil {
		t.Fatal(err)
	}

	expected := ProcStat{
		Pid:      24984,
		Command:  "godzilla",
		State:    StateRunning,
		RSS:      51579,
		Pagesize: 4096,
	}

	if stat != expected {
		t.Errorf("unexpected ProcStat: %v", stat)
	}
}

func makeStatFile(t *testing.T, fixture []byte) *os.File {
	tmpFile, err := ioutil.TempFile("", "procstat")
	if err != nil {
		t.Fatal(err)
	}

	_, err = tmpFile.Write(fixture)
	if err != nil {
		t.Fatal(err)
	}

	return tmpFile
}

func rmStatFile(file *os.File) {
	os.Remove(file.Name())
}
