package log

type Log struct {
	Hash    uint64
	Message string
	Time    string
}

func (l *Log) foo() {
}
