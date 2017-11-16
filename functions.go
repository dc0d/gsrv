package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pkg/errors"
)

//-----------------------------------------------------------------------------

func process(events <-chan entry, done chan<- entry) {
	go func() {
		for ev := range events {
			ev := ev
			go func() {
				defer func() { done <- ev }()
				dir := ev.dir()

				cmd := exec.Command("ggen")
				cmd.Dir = dir

				stdout, err := cmd.StdoutPipe()
				if err != nil {
					logerr.Println(err)
					return
				}
				defer stdout.Close()

				stderr, err := cmd.StderrPipe()
				if err != nil {
					logerr.Println(err)
					return
				}
				defer stderr.Close()

				go func() { io.Copy(os.Stdout, stdout) }()
				go func() { io.Copy(os.Stderr, stderr) }()

				if err := cmd.Start(); err != nil {
					logerr.Println(err)
					return
				}
				if err := cmd.Wait(); err != nil {
					logerr.Println(err)
					return
				}
			}()
		}
	}()
}

//-----------------------------------------------------------------------------

func throttle(notifications chan fsnotify.Event) (<-chan entry, chan<- entry) {
	var (
		events = make(chan entry, 128)
		done   = make(chan entry, 128)

		reentryAge = time.Millisecond * 600
		purgeAge   = time.Second * 5
	)
	go func() {
		defer func() {
			close(events)
		}()
		entries := make(map[string]eventry)
		for {
			select {
			case ev, ok := <-notifications:
				if !ok {
					return
				}
				ent := eventry{
					at:    time.Now(),
					entry: entry{ev},
				}
				if old, ok := entries[ent.dir()]; ok {
					ent = old
					ent.at = time.Now()
				}
				entries[ent.dir()] = ent
			case ent := <-done:
				delete(entries, ent.dir())
			case <-time.After(time.Millisecond * 300):
				// delete old unhandled entries
				var indexes []string
				for k, v := range entries {
					if v.signaled && v.age() >= purgeAge {
						indexes = append(indexes, k)
					}
				}
				for _, vk := range indexes {
					delete(entries, vk)
				}
				indexes = nil

				// select aged events
				var aged []entry
				for _, v := range entries {
					if v.signaled || v.age() < reentryAge {
						continue
					}
					v.signaled = true
					entries[v.dir()] = v
					aged = append(aged, v.entry)
				}

				go func() {
					for _, v := range aged {
						events <- v
					}
				}()
			}
		}
	}()
	return events, done
}

type entry struct {
	fsnotify.Event
}

func (ev entry) dir() string {
	dir := ev.Name
	if dirExists(dir) == errNotDir {
		dir = filepath.Dir(dir)
	}
	return dir
}

type eventry struct {
	at       time.Time
	signaled bool
	entry
}

func (ev eventry) age() time.Duration { return time.Since(ev.at) }

func (ev eventry) String() string {
	return fmt.Sprintf("<%v | %v | %v>", ev.signaled, ev.entry, ev.at)
}

//-----------------------------------------------------------------------------

func checkSrcDir() (string, error) {
	gopath := os.Getenv("GOPATH")
	if err := dirExists(gopath); err != nil {
		return "", errors.WithMessage(err, fmt.Sprintf("not found, $GOPATH = %v", gopath))
	}

	parts := strings.Split(gopath, string([]rune{filepath.ListSeparator}))
	gopath = parts[0]

	src := filepath.Join(gopath, "src")
	if err := dirExists(src); err != nil {
		return "", errors.WithMessage(err, "src directory not found")
	}

	return src, nil
}

//-----------------------------------------------------------------------------

func dirExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return errNotDir
	}
	return nil
}

//-----------------------------------------------------------------------------

func waitExit() {
	stopSignal := make(chan struct{})
	onSignal(func() { close(stopSignal) })
	<-stopSignal
}

//-----------------------------------------------------------------------------

type sentinelErr string

func (v sentinelErr) Error() string { return string(v) }

func errorf(format string, a ...interface{}) error {
	return sentinelErr(fmt.Sprintf(format, a...))
}

//-----------------------------------------------------------------------------

func onSignal(f func(), sig ...os.Signal) {
	if f == nil {
		return
	}
	sigc := make(chan os.Signal, 1)
	if len(sig) > 0 {
		signal.Notify(sigc, sig...)
	} else {
		signal.Notify(sigc,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT,
			syscall.SIGSTOP,
			syscall.SIGTSTP,
			syscall.SIGKILL)
	}
	go func() {
		<-sigc
		f()
	}()
}

//-----------------------------------------------------------------------------
