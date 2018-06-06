package callbag

import (
	"sync"
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

			start := func() {
				i := -1
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
			}

			var once sync.Once
			sink(NewGreets(func(p Payload) {
				once.Do(start)
				if _, ok := p.(Terminate); ok {
					close(clear)
				}
			}))

		default:
			return
		}
	}
}
