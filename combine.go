package callbag

// Combine is a Callbag factory that combines the latest data points from multiple
// (2 or more) callbag sources. It delivers those latest values as an array. Works
// with both pullable and listenable sources.
//
func Combine(sources ...Source) Source {
	EMPTY := &struct{}{}

	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			n := len(sources)
			sink := v.Source()

			if n == 0 {
				sink(NewGreets(func(p Payload) {}))
				sink(NewData(make([]interface{}, 0)))
				sink(NewTerminate(nil))
				return
			}

			Ns := n
			Nd := n
			Ne := n

			vals := make([]interface{}, n)
			sourceTalkbacks := make([]Source, n)

			var talkback Source = func(p Payload) {
				switch p.(type) {
				case Greets:
					return
				default:
					for _, sourceTalkback := range sourceTalkbacks {
						sourceTalkback(p)
					}
				}
			}

			for i, source := range sources {
				func(source Source, i int) {
					vals[i] = EMPTY
					source(NewGreets(func(p Payload) {
						switch v := p.(type) {
						case Greets:
							sourceTalkbacks[i] = v.Source()
							Ns--
							if Ns == 0 {

								sink(NewGreets(talkback))
							}
						case Data:
							_Nd := 0
							if Nd != 0 {
								if vals[i] == EMPTY {
									Nd--
								}
								_Nd = Nd
							}
							vals[i] = v.Value()
							if _Nd == 0 {
								arr := make([]interface{}, n)
								for j := 0; j < n; j++ {
									arr[j] = vals[j]
								}
								sink(NewData(arr))
							}
						case Terminate:
							Ne--
							if Ne == 0 {
								sink(NewTerminate(nil))
							}
						default:
							sink(p)
						}
					}))
				}(source, i)
			}

		default:
			return
		}
	}
}
