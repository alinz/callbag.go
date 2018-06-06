package callbag

import "reflect"

func Flatten() Transform {
	return func(source Source) Source {
		var sink Source
		var outerEnded bool
		var outerTalkback Source
		var innerTalkback Source

		isNil := func(val interface{}) bool {
			return val == nil || reflect.ValueOf(val).IsNil()
		}

		talkback := func(p Payload) {
			switch v := p.(type) {
			case Data:
				if !isNil(innerTalkback) {
					innerTalkback(v)
				} else if !isNil(outerTalkback) {
					outerTalkback(v)
				}
			case Terminate:
				if !isNil(innerTalkback) {
					innerTalkback(v)
				}
				if !isNil(outerTalkback) {
					outerTalkback(v)
				}
			}
		}

		return func(p Payload) {

			switch v := p.(type) {
			case Greets:
				sink = v.Source()

				source(NewGreets(func(P Payload) {
					switch V := P.(type) {
					case Greets:
						outerTalkback = V.Source()
						sink(NewGreets(talkback))
					case Data:
						innerSource := V.Value().(Source)
						if !isNil(innerTalkback) {
							innerTalkback(NewTerminate(nil))
						}

						innerSource(NewGreets(func(p Payload) {
							switch v := p.(type) {
							case Greets:
								innerTalkback = v.Source()
								innerTalkback(NewData(nil))
							case Data:
								sink(v)
							case Terminate:
								if isNil(v.Error()) {
									if outerEnded {
										sink(NewTerminate(nil))
									} else {
										innerTalkback = nil
										outerTalkback(NewData(nil))
									}
								} else {
									sink(v)
								}
							}
						}))
					case Terminate:
						if isNil(V.Error()) {
							if isNil(innerTalkback) {
								sink(NewTerminate(nil))
							} else {
								outerEnded = true
							}
						} else {
							sink(V)
						}
					}
				}))

			default:
				return
			}
		}
	}
}
