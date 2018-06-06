package callbag

func Observe(op func(val interface{})) Sink {
	return func(source Source) {
		done := make(chan struct{}, 1)

		source(NewGreets(func(p Payload) {
			switch v := p.(type) {
			case Data:
				op(v.Value())
			case Terminate:
				close(done)
			}
		}))

		<-done
	}
}
