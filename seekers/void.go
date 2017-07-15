package seeker

import "os"

type voidSeeker bool

// Seek implementation
func (vs voidSeeker) Seek(*os.File) (int64, error) { return 0, nil }
