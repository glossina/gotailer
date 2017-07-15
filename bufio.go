package gotailer

import (
	"bytes"
	"io"
	"os"
)

// BufReader struct
type BufReader struct {
	file *os.File
	buf  []byte
	cur  []byte
}

// NewBufReaderSize creates new BufReader
func NewBufReaderSize(file *os.File, size int) *BufReader {
	return &BufReader{
		file: file,
		buf:  make([]byte, size),
		cur:  nil,
	}
}

// ReadLine reads line in a manner similar to bufio.ReadLine with the only
// exception it returns isPrefix = true when it reach file end without the final LF
func (br *BufReader) ReadLine() (res []byte, isPrefix bool, err error) {
	if br.cur == nil || len(br.cur) == 0 {
		n, err := br.file.Read(br.buf)
		if err != nil && err != io.EOF {
			return nil, false, err
		}
		if n == 0 && err == io.EOF {
			return br.buf[:0], true, io.EOF
		}
		br.cur = br.buf[:n]
	}
	pos := bytes.IndexByte(br.cur, '\n')
	if pos < 0 {
		res = br.cur
		br.cur = br.cur[:0]
		return res, true, err
	}
	res = br.cur[:pos]
	br.cur = br.cur[pos+1:]
	return res, false, err
}

// Switch switches the file beneath
func (br *BufReader) Switch(file *os.File) {
	br.file = file
}
