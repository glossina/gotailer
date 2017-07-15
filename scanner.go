package gotailer

import (
	"io"
	"os"
)

// Scanner scans file line by line
type Scanner struct {
	file   *os.File
	reader *BufReader
	buf    []byte
	line   []byte
	err    error

	pos int64 // for telling the position
}

// NewScanner constructs new scanner
// file — previously opened file object
// size — initial buffer size (can extend in the process)
// bufSize — buffered reader buffer size
func NewScanner(file *os.File, size int, bufSize int) *Scanner {
	res := &Scanner{
		file:   file,
		reader: NewBufReaderSize(file, bufSize),
		buf:    make([]byte, 0, size),
		err:    nil,
	}
	res.initPos()
	return res
}

func (s *Scanner) initPos() (err error) {
	s.pos, err = s.file.Seek(0, os.SEEK_CUR)
	return
}

// Scan performs scanning
// finish - true, if file will not grow anymore (use Test result value to get it)
// raises panic on error
func (s *Scanner) Scan(finish bool) bool {
	line, isPrefix, err := s.reader.ReadLine()
	if err != nil && err != io.EOF {
		s.err = err
		return false
	}
	if !isPrefix {
		if len(s.buf) > 0 {
			s.buf = append(s.buf, line...)
			s.line = s.buf
			s.buf = s.buf[:0]
		} else {
			s.line = line
		}
		s.pos += int64(len(s.line)) + 1
		return len(s.line) > 0
	}
	s.buf = append(s.buf, line...)
	for isPrefix {
		if err == io.EOF {
			if finish {
				s.line = s.buf
				s.pos += int64(len(s.line))
				return len(s.line) > 0
			}
			return false
		}
		line, isPrefix, err = s.reader.ReadLine()
		if err != nil && err != io.EOF {
			s.err = err
			return false
		}
		s.buf = append(s.buf, line...)
	}
	s.line = s.buf
	s.pos += int64(len(s.line)) + 1
	return len(s.line) > 0
}

// Switch switches to another file
func (s *Scanner) Switch(file *os.File) error {
	if err := s.file.Close(); err != nil {
		return err
	}
	s.file = file
	if err := s.initPos(); err != nil {
		return err
	}
	s.reader.Switch(file)
	s.buf = s.buf[:0]
	return nil
}

// Bytes returns bytes read from a line
func (s *Scanner) Bytes() []byte {
	return s.line
}

// Err returns error
func (s *Scanner) Err() error {
	return s.err
}

// Tell returns current position
func (s *Scanner) Tell() int64 {
	return s.pos
}
