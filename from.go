package callbag

// FromIter converts list of items to a callbag pullable source
//
func FromIter(arr ...Value) Source {
	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			i := 0
			sink := v.Source()

			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if i < len(arr) {
						val := arr[i]
						i++
						sink(NewData(val))
					} else {
						sink(NewTerminate(nil))
					}
				}
			}))

		default:
			return
		}
	}
}

func FromRange(from, to int) Source {
	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			sink := v.Source()
			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if from < to {
						from++
						sink(NewData(from))
					} else {
						sink(NewTerminate(nil))
					}
				}
			}))
		default:
			return
		}
	}
}

func FromRangeChan(from, to int) Source {

}
