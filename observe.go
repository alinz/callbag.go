package callbag

import "sync"

// Observe is a Callbag listener sink that gets data from any listenable source.
// Think of it as "observable.subscribe".
//
// NOTE: It will wait until source stops sending data.
//
func Observe(op func(val interface{})) Sink {
	return func(source Source) {
		done := make(chan struct{}, 1)
		once := sync.Once{}

		source(NewGreets(func(p Payload) {
			switch v := p.(type) {
			case Data:
				op(v.Value())
			case Terminate:
				once.Do(func() {
					close(done)
				})
			}
		}))

		<-done
	}
}
