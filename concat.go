package callbag

// Concat is a Callbag factory that concatenates the data from multiple (2 or more)
// callbag sources. It starts each source at a time: waits for the previous
// source to end before starting the next source. Works with both pullable
// and listenable sources.
//
func Concat(sources ...Source) Source {
	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			n := len(sources)
			sink := v.Source()

			if n == 0 {
				sink(NewGreets(func(p Payload) {}))
				sink(NewTerminate(nil))
				return
			}

			i := 0
			var sourceTalkback Source
			talkback := func(p Payload) {
				switch p.(type) {
				case Data:
					sourceTalkback(p)
				case Terminate:
					sourceTalkback(p)
				default:
					return
				}
			}

			var next func()
			next = func() {
				if i == n {
					sink(NewTerminate(nil))
					return
				}

				source := sources[i]
				source(NewGreets(func(p Payload) {
					switch v := p.(type) {
					case Greets:
						sourceTalkback = v.Source()
						if i == 0 {
							sink(NewGreets(talkback))
						} else {
							sourceTalkback(NewData(nil))
						}
					case Data:
						sink(v)
					case Terminate:
						i++
						next()
					}
				}))
			}

			next()

		default:
			return
		}
	}
}
