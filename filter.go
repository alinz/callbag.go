package callbag

func Filter(cond func(val interface{}) bool) Transform {
	return func(source Source) Source {
		return func(p Payload) {
			var talkback Source

			switch v := p.(type) {
			case Greets:
				sink := v.Source()
				source(NewGreets(func(p Payload) {
					switch v := p.(type) {
					case Greets:
						talkback = v.Source()
						sink(v)
					case Data:
						if cond(v.Value()) {
							sink(v)
						} else {
							talkback(NewData(nil))
						}
					case Terminate:
						sink(v)
					}
				}))
			default:
				return
			}
		}
	}
}
