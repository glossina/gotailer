# gotailer
Go synchronous file tailing API

### Preamble
For our task we didn't tail thousands or hundreds or even just tens of files. These were just 6-8 files that were regularly written (hundreds lines per second),
so the most effective way of dealing this was just wait for a second (or whatever) and read all new data written. Hence the API developed
which is aimed for dealing with small amount of frequently updated files.

Get it
```bash
go get github.com/sirkon/gotailer
```
and use as shown in

### Usage Example
```go
package main

import (
	"log"
	"time"

	"github.com/sirkon/gotailer"
	seeker "github.com/sirkon/gotailer/seekers"
)

const fileName = "/tmp/file_to_watch"

func main() {
	// Setting up watcher object
	watcher, err := gotailer.NewPollingMonitor(fileName, seeker.SeekToEnd, seeker.SeekToEnd)
	if err != nil {
		log.Fatalf("Unrecoverable error: %s", err)
	}
	// Watcher object clearance
	defer func() {
		err := watcher.GetFile().Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Setting up scanner object. No clearance needed as it runs over
	// watcher's provided objects
	scanner := gotailer.NewScanner(watcher.GetFile(), 1024, 512*1024)
	if err != nil {
		log.Fatalf("Unrecoverable error: %s", err)
	}

	for {
		// Check if new file is ready
		finished, err := watcher.Test()
		if err != nil {
			log.Fatalf("Unrecoverable error: %s", err)
		}

		// Read out the current file. Finished set to true force reader to treat the rest of file which
		// is not finished with LF as a line. Otherwise, it will wait for the rest of line.
		for scanner.Scan(finished) {
			log.Println(string(scanner.Bytes()))
			log.Printf("Now at %d position", scanner.Tell())
		}

		// Switch sources if a new file is ready for reading since the current one has just been read out
		if finished {
			if err := scanner.Switch(watcher.GetFile()); err != nil {
				log.Fatal(err)
			}
			continue
		}

		time.Sleep(time.Second)
	}
}
```
