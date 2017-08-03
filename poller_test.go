package gotailer

import (
	"io/ioutil"
	"os"
	"testing"

	seeker "github.com/sirkon/gotailer/seekers"
)

func TestPooler(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		if err := tmpfile.Close(); err != nil {
			t.Fatal(err)
		}
	}()

	fname := tmpfile.Name()
	monitor, err := NewPollingMonitor(fname, seeker.SeekToEnd, seeker.SeekToEnd)
	if err != nil {
		t.Fatal(err)
	}

	// Nothing was appended to file, should stay
	ok, err := monitor.Test()
	if err != nil {
		t.Error(err)
		return
	}
	if ok {
		t.Error(falseRequired)
		return
	}

	// Should stay again
	ok, err = monitor.Test()
	if err != nil {
		t.Error(err)
		return
	}
	if ok {
		t.Error(falseRequired)
		return
	}

	// Write to file. Still should stay.
	_, err = tmpfile.WriteString("1\n")
	if err != nil {
		t.Fatal(err)
	}
	ok, err = monitor.Test()
	if err != nil {
		t.Error(err)
		return
	}
	if ok {
		t.Error(falseRequired)
		return
	}

	// Remove file. Create new one empty. Should switch now.
	os.Remove(fname)
	file, err := os.OpenFile(fname, os.O_CREATE, os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	ok, err = monitor.Test()
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error(trueRequired)
		return
	}

	// Should switch again. Because the file is different.
	os.Remove(fname)
	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}
	file, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatal(err)
		}
		os.Remove(fname)
	}()
	ok, err = monitor.Test()
	if err != nil {
		t.Error(err)
		return
	}
	if !ok {
		t.Error(trueRequired)
		return
	}

	// Write to file. Should not switch
	_, err = file.WriteString("1\n2\n")
	if err != nil {
		t.Fatal(err)
	}
	ok, err = monitor.Test()
	if err != nil {
		t.Error(err)
		return
	}
	if ok {
		t.Error(falseRequired)
		return
	}
}
