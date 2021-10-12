package view

import (
	"fmt"
	"time"

	"github.com/diamondburned/gotk4/pkg/core/glib"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

type stopwatchState uint

const (
	running stopwatchState = iota
	stopped
)

type Stopwatch struct {
	*gtk.Label

	seconds time.Time
	handle  chan stopwatchState
	state   stopwatchState
	ticker  *time.Ticker
}

func NewStopwatch(label *gtk.Label) *Stopwatch {
	return &Stopwatch{
		Label:   label,
		seconds: time.Unix(0, 0),
		state:   stopped,
		handle:  make(chan stopwatchState),
	}
}

func (s *Stopwatch) Start() {
	if s.state == stopped {
		s.state = running
		s.ticker = time.NewTicker(time.Second)

		go func() {
			state := running

			for {
				select {
				case state = <-s.handle:
					switch state {
					case stopped:
						return
					}
				case <-s.ticker.C:
					glib.IdleAdd(func() {
						s.SetText(fmt.Sprint(s))
					})

					s.seconds = s.seconds.Add(time.Second)
				}
			}
		}()
	}
}

func (s *Stopwatch) String() string {
	time := s.seconds.UTC()
	return fmt.Sprintf("%02d:%02d:%02d", time.Hour(), time.Minute(), time.Second())
}

func (s *Stopwatch) Reset() {
	s.state = stopped
	s.seconds = time.Unix(0, 0)
	s.SetText(fmt.Sprint(s))
}

func (s *Stopwatch) Stop() {
	if s.state == running {
		s.state = stopped
		s.handle <- stopped
		s.ticker.Stop()
	}
}

func (s *Stopwatch) Restart() {
	s.Stop()
	s.Reset()
	s.Start()
}
