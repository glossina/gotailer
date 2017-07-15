package session

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	seeker "github.com/DenisCheremisov/gotailer/seekers"
)

// FileSession keeps tailer session information in files
type FileSession struct {
	root string
}

// NewFileSession constructor
func NewFileSession(root string) *FileSession {
	return &FileSession{
		root: root,
	}
}

// SavePosition ...
func (fs *FileSession) SavePosition(id string, pos int64) (err error) {
	return ioutil.WriteFile(filepath.Join(fs.root, id), []byte(fmt.Sprintf("%d", pos)), 0644)
}

// RestorePosition ...
func (fs *FileSession) RestorePosition(id string) (res seeker.Seeker, err error) {
	posData, err := ioutil.ReadFile(filepath.Join(fs.root, id))
	if os.IsNotExist(err) {
		res = seeker.SeekToStart
		return
	}
	if err != nil {
		return
	}
	pos, err := strconv.ParseInt(string(posData), 10, 64)
	if err != nil {
		return
	}
	res = seeker.SeekTo(pos)
	return
}
