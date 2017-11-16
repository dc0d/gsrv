package main

import (
	stdlog "log"
	"os"
)

//-----------------------------------------------------------------------------

// errors
var (
	errNotDir = errorf("NOT A DIR")
)

//-----------------------------------------------------------------------------

var (
	conf struct{}

	logerr = stdlog.New(os.Stderr, "err: ", 0)
	loginf = stdlog.New(os.Stdout, "inf: ", 0)
	logwrn = stdlog.New(os.Stdout, "wrn: ", 0)
)

//-----------------------------------------------------------------------------

func init() {
	const dgb = false
	if dgb {
		logerr = stdlog.New(os.Stderr, "err: ", stdlog.Ltime|stdlog.Lshortfile)
		loginf = stdlog.New(os.Stdout, "inf: ", stdlog.Ltime|stdlog.Lshortfile)
		logwrn = stdlog.New(os.Stdout, "wrn: ", stdlog.Ltime|stdlog.Lshortfile)
	}
}

//-----------------------------------------------------------------------------
