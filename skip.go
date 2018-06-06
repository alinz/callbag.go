package callbag

// Skip is a Callbag operator that skips the first N data points of a source.
// Works on either pullable and listenable sources.
//
func Skip(max int) Transform {
	return func(source Source) Source {
		return func(p Payload) {
			switch v := p.(type) {
			case Greets:
				sink := v.Source()
				skipped := 0
				var talkback Source

				source(NewGreets(func(p Payload) {
					switch v := p.(type) {
					case Greets:
						talkback = v.Source()
						sink(v)
					case Data:
						if skipped < max {
							skipped++
							talkback(NewData(nil))
						} else {
							sink(v)
						}
					default:
						sink(v)
					}
				}))

			default:
				return
			}
		}
	}
}
