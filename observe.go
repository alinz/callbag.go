package callbag

func Observe(op func(val interface{})) Sink {
	return func(source Source) {

		source(NewGreets(func(p Payload) {
			switch v := p.(type) {
			case Data:
				op(v.Value())
			}
		}))

		select {}
	}
}
