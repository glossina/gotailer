package gotailer

import (
	"fmt"
	"os"

	seeker "github.com/sirkon/gotailer/seekers"
)

// PollingMonitor is a file monitor implementation based on file pooling
type PollingMonitor struct {
	name         string
	file         *os.File
	seeker       seeker.Seeker
	prevPos      int64
	reopenSeeker seeker.Seeker
}

// NewPollingMonitor creats new pooling file monitor
// name - file name to monitor
// seeker - function to seek over a file
// reopenSeeker - seek on reopen
func NewPollingMonitor(name string, seeker seeker.Seeker, reopenSeeker seeker.Seeker) (*PollingMonitor, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	var pos int64
	pos, err = seeker.Seek(file)
	if err != nil {
		if err = file.Close(); err != nil {
			return nil, err
		}
		return nil, err
	}

	res := &PollingMonitor{
		file:         file,
		name:         name,
		seeker:       seeker,
		prevPos:      pos,
		reopenSeeker: reopenSeeker,
	}
	return res, nil
}

// Test function
func (pm *PollingMonitor) Test() (ok bool, err error) {
	ok = false
	file, err := os.Open(pm.name)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			if saveErr := file.Close(); saveErr != nil {
				err = fmt.Errorf("Unrecoverable error: %s, %s previously", err, saveErr)
				return
			}
		}
	}()

	stat, err := file.Stat()
	if err != nil {
		return
	}

	prevStat, err := pm.file.Stat()
	if err != nil {
		return
	}

	if !os.SameFile(stat, prevStat) || pm.prevPos > stat.Size() {
		// We should just leave currently opened file: the user
		// may need to read it out
		var pos int64
		pos, err = pm.reopenSeeker.Seek(file)
		if err != nil {
			return
		}
		pm.file = file
		pm.prevPos = pos
		ok = true
	} else {
		pm.prevPos = prevStat.Size()
	}
	return
}

// GetFile function
func (pm *PollingMonitor) GetFile() *os.File {
	return pm.file
}
