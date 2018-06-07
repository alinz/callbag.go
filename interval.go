package callbag

import (
	"time"
)

func setInterval(fn func(), timeout time.Duration) func() {
	ticker := time.NewTicker(timeout)
	clear := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				fn()
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	return func() {
		close(clear)
	}
}

// Interval is a callbag listenable source that sends incremental numbers every x period.
//
func Interval(period time.Duration) Source {
	return func(p Payload) {
		var sink Source
		var i int

		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			clear := setInterval(func() {
				sink(NewData(i))
				i++
			}, period)

			sink(NewGreets(func(p Payload) {
				if _, ok := p.(Terminate); ok {
					clear()
				}
			}))

		default:
			return
		}
	}
}

// PausableInterval is a callbag listenable source that sends incremental numbers every x period
// but can be paused (and resumed) when it is pulled by a sink.
//
// NOTE: Don't use forEach directly as the sink for this source, because forEach pulls every time it receives data.
// You can use this source as the argument for sample, though.
//
func PausableInterval(period time.Duration) Source {
	return func(p Payload) {
		var sink Source
		var i int
		var clear func()

		resume := func() {
			clear = setInterval(func() {
				sink(NewData(i))
				i++
			}, period)
		}

		pause := func() {
			if !isNil(clear) {
				clear()
				clear = nil
			}
		}

		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if isNil(clear) {
						resume()
					} else {
						pause()
						clear = nil
					}
				case Terminate:
					if isNil(clear) {
						pause()
					}
				}
			}))

			resume()

		default:
			return
		}
	}
}
