package callbag

import (
	"time"
)

func Interval(period time.Duration) Source {
	return func(p Payload) {
		var sink Source

		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			ticker := time.NewTicker(period)
			clear := make(chan bool)

			func() {
				i := 0
				go func() {
					for {
						select {
						case <-ticker.C:
							sink(NewData(i))
							i++
						case <-clear:
							ticker.Stop()
							return
						}
					}
				}()
			}()

			sink(NewGreets(func(p Payload) {
				if _, ok := p.(Terminate); ok {
					close(clear)
				}
			}))

		default:
			return
		}
	}
}
