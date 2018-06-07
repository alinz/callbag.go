package callbag

// Map Callbag operator that applies a transformation on data passing through it.
//
func Map(op func(val Value) Value) Transform {
	return func(source Source) Source {
		return func(p Payload) {
			switch v := p.(type) {
			case Greets:
				sink := v.Source()
				source(NewGreets(func(p Payload) {
					switch v := p.(type) {
					case Greets:
						sink(v)
					case Data:
						sink(NewData(op(v.Value())))
					case Terminate:
						sink(NewTerminate(v.Error()))
					}
				}))
			default:
				return
			}
		}
	}
}
