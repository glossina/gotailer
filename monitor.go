package gotailer

import "os"

// Monitor generalizes file change monitoring
type Monitor interface {
	// Test if we need to switch or reopen to another file
	Test() (bool, error)

	// Get file handler
	GetFile() *os.File
}
