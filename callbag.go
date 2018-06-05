package callbag

type Payload interface {
	Type() int
}

type Source func(Payload)          // => (start, sink) => {}
type Transform func(Source) Source // => source => (start, sink) => {}
type Sink func(Source)             // => source => {}

type Greets interface {
	Payload
	Source() Source
}

type Terminate interface {
	Payload
	Error() error
}

type Data interface {
	Payload
	Value() interface{}
}

type greets struct {
	source Source
}

func (g *greets) Type() int {
	return 0
}

func (g *greets) Source() Source {
	return g.source
}

type data struct {
	value interface{}
}

func (d *data) Type() int {
	return 1
}

func (d *data) Value() interface{} {
	return d.value
}

type terminate struct {
	err error
}

func (t *terminate) Type() int {
	return 2
}

func (t *terminate) Error() error {
	return t.err
}

func NewGreets(source Source) Payload {
	return &greets{source}
}

func NewData(value interface{}) Payload {
	return &data{value}
}

func NewTerminate(err error) Payload {
	return &terminate{err}
}
