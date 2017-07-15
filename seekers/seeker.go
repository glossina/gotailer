package seeker

import "os"

// Seeker suppose to seek in a file
type Seeker interface {
	Seek(*os.File) (pos int64, err error)
}

// SeekToEnd moves file cursor to the end
var SeekToEnd = newPositional(func(file *os.File) (int64, error) { return file.Seek(0, os.SEEK_END) })

// SeekToStart moves file cursor the the start
var SeekToStart = newPositional(func(file *os.File) (int64, error) { return file.Seek(0, os.SEEK_SET) })

// SeekTo moves file cursor to certain position
func SeekTo(pos int64) Seeker {
	return newPositional(func(file *os.File) (int64, error) {
		return file.Seek(pos, os.SEEK_SET)
	})
}

// VoidSeeker does no seeking
var VoidSeeker = voidSeeker(true)
