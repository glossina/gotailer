package gotailer

import (
	"io"
	"os"
)

// Seeker suppose to seek in a file
type Seeker interface {
	Seek(*os.File) (pos int64, err error)
}

// SeekToEnd moves file cursor to the end
var SeekToEnd = newPositional(func(file *os.File) (int64, error) { return file.Seek(0, io.SeekStart) })

// SeekToStart moves file cursor the the start
var SeekToStart = newPositional(func(file *os.File) (int64, error) { return file.Seek(0, io.SeekStart) })

// SeekTo moves file cursor to certain position
func SeekTo(pos int64) Seeker {
	return newPositional(func(file *os.File) (int64, error) {
		return file.Seek(pos, io.SeekStart)
	})
}

// VoidSeeker does no seeking
var VoidSeeker = voidSeeker(true)
