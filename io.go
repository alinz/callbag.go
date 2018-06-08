package callbag

import (
	"io"
)

func ToWriter(writer io.Writer) Sink {
	return func(source Source) {
		buffer := make([]byte, 1)
		source(NewGreets(func(p Payload) {
			var talkback Source

			source(NewGreets(func(p Payload) {
				switch v := p.(type) {
				case Greets:
					talkback = v.Source()
					talkback(NewData(nil))
				case Data:
					buffer[0] = v.Value().(byte)
					writer.Write(buffer)
					talkback(NewData(nil))
				}
			}))
		}))
	}
}

func FromReader(reader io.Reader, size int) Source {
	var sink Source
	var isPumping bool
	var done bool

	buffer := make([]byte, size)

	pump := func() {
		for {
			if n, err := reader.Read(buffer); err == nil {
				sink(NewData(buffer[:n]))
			} else {
				done = true
				sink(NewTerminate(err))
				return
			}
		}
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
		}
	}
}
