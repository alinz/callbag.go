package callbag

// Take is a Callbag operator that limits the amount of data sent by a source.
// Works on either pullable and listenable sources.
//
func Take(max int) Transform {
	return func(source Source) Source {
		return func(p Payload) {
			switch v := p.(type) {
			case Greets:
				var sourceTalkback Source
				taken := 0
				sink := v.Source()

				talkback := func(p Payload) {
					if taken < max {
						sourceTalkback(p)
					}
				}

				source(NewGreets(func(p Payload) {
					switch v := p.(type) {
					case Greets:
						sourceTalkback = v.Source()
						sink(NewGreets(talkback))
					case Data:
						if taken < max {
							taken++
							sink(v)
							if taken == max {
								sink(NewTerminate(nil))
								sourceTalkback(NewTerminate(nil))
							}
						}
					default:
						sink(v)
					}
				}))
			}
		}
	}
}
