package session

import seeker "github.com/glossina/gotailer/seekers"

// TailerSession is a tailer position saver.
type TailerSession interface {
	SavePosition(id string, pos int64) (err error)
	RestorePosition(id string) (res seeker.Seeker, err error)
}
