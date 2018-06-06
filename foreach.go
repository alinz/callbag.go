package callbag

// ForEach is a Callbag sink that iterate over pullable data
//
// NOTE: If you want to iterate over listenable sources, use
// Observe instead
//
func ForEach(op func(interface{})) Sink {
	return func(source Source) {
		var talkback Source

		source(NewGreets(func(p Payload) {
			switch v := p.(type) {
			case Greets:
				talkback = v.Source()
				talkback(NewData(nil))
			case Data:
				op(v.Value())
				talkback(NewData(nil))
			}
		}))
	}
}
