package callbag

// Merge is a Callbag factory that merges data from multiple callbag sources.
// Works well with listenable sources, and while it may work for some pullable sources,
// it is only designed for listenable sources.
//
func Merge(sources ...Source) Source {
	return func(p Payload) {
		var sink Source
		switch v := p.(type) {
		case Greets:
			startCount := 0
			endCount := 0
			n := len(sources)
			sourceTalkbacks := make([]Source, n)
			sink = v.Source()

			talkback := func(p Payload) {
				if _, ok := p.(Greets); ok {
					return
				}
				for _, sourceTalkback := range sourceTalkbacks {
					if !isNil(sourceTalkback) {
						sourceTalkback(p)
					}
				}
			}

			for i, source := range sources {
				source(NewGreets(func(p Payload) {
					switch v := p.(type) {
					case Greets:
						sourceTalkbacks[i] = v.Source()
						if startCount == 1 {
							sink(NewGreets(talkback))
						}
						startCount++
					case Data:
						sink(v)
					case Terminate:
						sourceTalkbacks[i] = nil
						if endCount == n {
							sink(NewTerminate(nil))
						}
						endCount++
					}
				}))
			}

		default:
			return
		}
	}
}
