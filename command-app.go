package main

import (
	"github.com/dc0d/dirwatch"
	"github.com/fsnotify/fsnotify"
	"github.com/urfave/cli"
)

func cmdApp(*cli.Context) error {
	var (
		err error
		src string
	)
	if src, err = checkSrcDir(); err != nil {
		return err
	}
	loginf.Printf("watching $GOPATH/src= %s", src)

	notifications := make(chan fsnotify.Event, 128)

	wcth, err := dirwatch.New(src, func(e fsnotify.Event) {
		notifications <- e
	})
	if err != nil {
		return err
	}
	defer wcth.Stop()

	events, done := throttle(notifications)
	process(events, done)

	waitExit()
	return nil
}
