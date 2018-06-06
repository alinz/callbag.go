package callbag

// Payload is a interface which is beening used in
// Greets, Data and Terminate it only carries Type which
// is protocol value of the types.
type Payload interface {
	Type() int
}

// Source is a function type which caries the base communication link
type Source func(Payload) // => (start, sink) => {}
// Transform is a function type which is used for any Callbag functions
// that transforms or sits between source and sink
type Transform func(Source) Source // => source => (start, sink) => {}
// Sink is a function type that identifies the end of the stream
type Sink func(Source) // => source => {}

// Greets is a handshake between Source and Sink
// it carries the Source to next Sink
type Greets interface {
	Payload
	Source() Source
}

// Terminate is an interface which signals the data is not coming or
// I don't need data to be pushed any more.
//
// NOTE: it can also carry error in case of error
//
type Terminate interface {
	Payload
	Error() error
}

// Data is an interface design to send a value to Sink or talkback
//
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

// NewGreets creates a new Greets interface based on given source
func NewGreets(source Source) Payload {
	return &greets{source}
}

// NewData creates a new Data interface based on given value
func NewData(value interface{}) Payload {
	return &data{value}
}

// NewTerminate create a new Terminate interface which wither carries error or not
func NewTerminate(err error) Payload {
	return &terminate{err}
}
