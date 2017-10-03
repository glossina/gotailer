package gotailer

// TailerSession is a tailer position saver.
type TailerSession interface {
	SavePosition(id string, pos int64) (err error)
	RestorePosition(id string) (res Seeker, err error)
}
