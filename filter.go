package callbag

// Filter is a Callbag operator that conditionally lets data pass through
//
// As an example, the fillowing code print out the event number from 1 to 5
//
// 		callbag.Pipe(
// 			callbag.FromIter(1, 2, 3, 4, 5),
// 			callbag.Filter(func(val interface{}) bool {
// 				n := val.(int)
// 				return n%2 == 0
// 			}),
// 			callbag.ForEach(func(val interface{}) {
// 				fmt.Println(val)
// 			}),
// 		)
//
func Filter(cond func(val Value) bool) Transform {
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
