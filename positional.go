package gotailer

import "os"

type positional struct {
	seeker func(file *os.File) (int64, error)
}

func newPositional(seeker func(*os.File) (int64, error)) *positional {
	return &positional{
		seeker: seeker,
	}
}

func (s *positional) Seek(file *os.File) (pos int64, err error) {
	return s.seeker(file)
}
