package callbag

// FromValues converts list of items to a callbag pullable source
//
func FromValues(values ...Value) Source {
	var sink Source
	var isPumping bool
	var done bool

	pump := func() {
		for _, value := range values {
			sink(NewData(value))
		}
		done = true
		sink(NewTerminate(nil))
	}

	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if !isPumping && !done {
						isPumping = true
						pump()
					}
				}
			}))

		default:
			return
		}
	}
}

// FromRange generate numbers from a number to a number to a callbag pullable source
//
func FromRange(from, to int) Source {
	var sink Source
	var isPumping bool
	var i int

	pump := func() {
		for i = from; i < to; i++ {
			sink(NewData(i))
		}
		sink(NewTerminate(nil))
	}

	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			sink = v.Source()
			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if !isPumping && i != to {
						isPumping = true
						pump()
					}
				}
			}))
		default:
			return
		}
	}
}

// FromBytes converts list of items to a callbag pullable source
//
func FromBytes(values []byte) Source {
	var sink Source
	var isPumping bool
	var done bool

	pump := func() {
		for _, value := range values {
			sink(NewData(value))
		}
		done = true
		sink(NewTerminate(nil))
	}

	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if !isPumping && !done {
						isPumping = true
						pump()
					}
				}
			}))

		default:
			return
		}
	}
}

// FromStrings converts list of items to a callbag pullable source
//
func FromStrings(values []string) Source {
	var sink Source
	var isPumping bool
	var done bool

	pump := func() {
		for _, value := range values {
			sink(NewData(value))
		}
		done = true
		sink(NewTerminate(nil))
	}

	return func(p Payload) {
		switch v := p.(type) {
		case Greets:
			sink = v.Source()

			sink(NewGreets(func(p Payload) {
				switch p.(type) {
				case Data:
					if !isPumping && !done {
						isPumping = true
						pump()
					}
				}
			}))

		default:
			return
		}
	}
}
