package gotailer

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	seeker "github.com/DenisCheremisov/gotailer/seekers"
	"github.com/stretchr/testify/assert"
)

const (
	trueRequired  = "Must be `true` here"
	falseRequired = "Must be `false` here"
)

type TestComparator struct {
	msg string
}

func (tc *TestComparator) Equal(expected string, actual []byte) bool {
	if expected != string(actual) {
		tc.msg = fmt.Sprintf("%s != %s", expected, string(actual))
		return false
	}
	return true
}

func TestScanner(t *testing.T) {
	tc := &TestComparator{}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up
	defer func() {
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fname := tmpfile.Name()
	file, err := os.Open(fname)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	firstTwo := "fa\nthis is a second line"
	_, err = tmpfile.WriteString(firstTwo)
	if err != nil {
		t.Fatal(err)
	}

	scanner := NewScanner(file, 1, 5)

	// Read first line, stuck at the second non-finished one
	if !scanner.Scan(false) {
		t.Error(trueRequired)
		return
	}
	data := scanner.Bytes()
	if !tc.Equal("fa", data) {
		t.Error(tc.msg)
		return
	}
	if !assert.Equal(t, int64(3), scanner.Tell()) {
		return
	}
	if scanner.Scan(false) {
		t.Error(falseRequired)
		return
	}
	if !assert.Equal(t, int64(3), scanner.Tell()) {
		return
	}

	// Read to the new line now
	_, err = tmpfile.WriteString("\n")
	if err != nil {
		t.Fatal(err)
	}
	if !scanner.Scan(false) {
		t.Error(trueRequired)
		return
	}
	data = scanner.Bytes()
	if !tc.Equal("this is a second line", data) {
		t.Error(tc.msg)
		return
	}
	if !assert.Equal(t, int64(len(firstTwo))+1, scanner.Tell()) {
		return
	}

	// File will not grow anymore, finish it
	_, err = tmpfile.WriteString("finish")
	if err != nil {
		t.Fatal(err)
	}
	if !scanner.Scan(true) {
		t.Error(trueRequired)
		return
	}
	data = scanner.Bytes()
	if !tc.Equal("finish", data) {
		t.Error(tc.msg)
		return
	}
	if !assert.Equal(t, int64(len(firstTwo))+1+int64(len("finish")), scanner.Tell()) {
		return
	}
}

func TestScannerWatcherCombined(t *testing.T) {
	tc := &TestComparator{}

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(tmpfile.Name()) // clean up
	defer func() {
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fname := tmpfile.Name()
	watcher, err := NewPollingMonitor(fname, seeker.SeekToEnd, seeker.SeekToStart)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := watcher.GetFile().Close(); err != nil {
			t.Fatal(err)
		}
	}()
	scanner := NewScanner(watcher.GetFile(), 100, 100)

	// Write five lines
	for i := 0; i < 5; i++ {
		if _, err := tmpfile.WriteString("bugagashechki\n"); err != nil {
			t.Fatal(err)
		}
	}

	// Watcher should stay
	finished, err := watcher.Test()
	if err != nil {
		t.Fatal(err)
	}
	if finished {
		t.Fatal(falseRequired)
	}
	i := 0

	// Should read 5 lines
	for scanner.Scan(finished) {
		i++
		if !tc.Equal("bugagashechki", scanner.Bytes()) {
			t.Fatal(tc.msg)
		}
	}
	if i != 5 {
		t.Fatalf("%d != 5", i)
	}
	if scanner.Err() != io.EOF && scanner.Err() != nil {
		t.Fatal(err)
	}

	// Truncate and write from the start
	if err := tmpfile.Truncate(0); err != nil {
		t.Fatal(err)
	}
	if _, err := tmpfile.Seek(0, os.SEEK_SET); err != nil {
		t.Fatal(err)
	}
	// Should reopen now
	finished, err = watcher.Test()
	if err != nil {
		t.Fatal(err)
	}
	if !finished {
		t.Fatal(trueRequired)
	}
	if _, err := tmpfile.WriteString("bugagashechki\n"); err != nil {
		t.Fatal(err)
	}

	// Old file is empty. Nothing should be read, switch to the next file
	for scanner.Scan(finished) {
		t.Fatal("Must not be here")
	}
	if err := scanner.Switch(watcher.GetFile()); err != nil {
		t.Fatal(err)
	}

	i = 0
	for scanner.Scan(false) {
		i++
		if !tc.Equal("bugagashechki", scanner.Bytes()) {
			t.Fatal(tc.msg)
		}
	}
	if i != 1 {
		t.Fatalf("%d != 1", i)
	}
}
