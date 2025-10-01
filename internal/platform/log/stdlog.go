package log

import "log"

type StdLogger struct{}

func (StdLogger) Infof(format string, args ...any) {
	log.Printf("INFO: "+format, args...)
}

func (StdLogger) Errorf(format string, args ...any) {
	log.Printf("ERROR: "+format, args...)
}
