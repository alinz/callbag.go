package callbag

func Subject() Source {
	sinks := make([]Source, 0)

	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			sink := v.Source()
			sinks = append(sinks, sink)

			sink(NewGreets(func(p Payload) {
				if _, ok := p.(Terminate); ok {
					idx := -1
					for i, s := range sinks {
						// Need to compare pointers than actual value of func type
						if &s == &sink {
							idx = i
							break
						}
					}

					if idx != -1 {
						sinks[idx] = sinks[len(sinks)-1]
						sinks[len(sinks)-1] = nil
						sinks = sinks[:len(sinks)-1]
					}
				}
			}))
		default:
			for _, sink := range sinks {
				sink(v)
			}
		}
	}
}
