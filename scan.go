package callbag

func Scan(reducer func(prev interface{}, curr interface{}) interface{}, seed interface{}) Transform {
	return func(source Source) Source {
		return func(p Payload) {
			switch v := p.(type) {
			case Greets:
				acc := seed
				sink := v.Source()

				source(NewGreets(func(p Payload) {
					switch v := p.(type) {
					case Data:
						acc = reducer(acc, v.Value())
						sink(NewData(acc))
					default:
						sink(p)
					}
				}))
			}
		}
	}
}
