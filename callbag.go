// All functions are designed to be async/lazy, which means you need to wait for them,
// or it will be killed once the program ends. One option is to wait for the end of the
// stream using ForEach or ToChannel function. All the examples provided the mechanisum
// to wait for the stream to be ended.
package callbag

import (
	"sync"
	"time"
)

type Type int

const (
	Hello Type = iota
	Data
	Bye
)

func (t Type) String() string {
	switch t {
	case Hello:
		return "Hello"
	case Data:
		return "Data"
	case Bye:
		return "Bye"
	default:
		return "Unknown"
	}
}

type Option[T any] struct {
	Data T
	Err  error
	Func Func[T]
	Type Type
}

type Func[T any] func(opt Option[T])

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

// FromInterval generates a tick and send an incermental integer value every n's `time.Duration`
func FromInterval[T number](period time.Duration) Func[T] {
	var down Func[T]

	return func(opt Option[T]) {
		if opt.Type != Hello {
			panic("Expected Hello")
		}

		down = opt.Func
		reqNextItem := make(chan struct{}, 1)
		go func() {
			var i T
			for {
				_, ok := <-reqNextItem
				if !ok {
					return
				}

				time.Sleep(period)
				i++

				down(Option[T]{
					Type: Data,
					Data: i,
				})
			}
		}()

		down(Option[T]{
			Type: Hello,
			Func: func(opt Option[T]) {
				switch opt.Type {
				case Data:
					reqNextItem <- struct{}{}
				case Bye:
					close(reqNextItem)
				}
			},
		})
	}
}

// FromSlice if the input is a slice, it can be converted into callbag stream by using
// FromSlice function as a source. Because slice is Finite, the stream will autoamatically
// ends once it processes the last item
func FromSlice[T any](values []T) Func[T] {
	var down Func[T]

	return func(opt Option[T]) {
		if opt.Type != Hello {
			panic("Expected Hello")
		}

		down = opt.Func
		reqNextItem := make(chan struct{}, 1)
		go func() {
			i := 0
			for {
				_, ok := <-reqNextItem
				if !ok {
					return
				}

				if i >= len(values) {
					down(Option[T]{
						Type: Bye,
					})
					return
				}

				down(Option[T]{
					Type: Data,
					Data: values[i],
				})
				i++
			}
		}()

		down(Option[T]{
			Type: Hello,
			Func: func(opt Option[T]) {
				switch opt.Type {
				case Data:
					reqNextItem <- struct{}{}
				case Bye:
					close(reqNextItem)
				}
			},
		})
	}
}

// FromRange generates number values from n to m with step value. It will stream will be terminated by the last iteration
func FromRange[T number](a, b, step T) Func[T] {
	var down Func[T]

	return func(opt Option[T]) {
		if opt.Type != Hello {
			panic("Expected Hello")
		}

		down = opt.Func
		reqNextItem := make(chan struct{}, 1)
		go func() {
			for {
				_, ok := <-reqNextItem
				if !ok {
					return
				}

				if !(a < b) {
					down(Option[T]{
						Type: Bye,
					})
					return
				}

				down(Option[T]{
					Type: Data,
					Data: a,
				})
				a += step
			}
		}()

		down(Option[T]{
			Type: Hello,
			Func: func(opt Option[T]) {
				switch opt.Type {
				case Data:
					reqNextItem <- struct{}{}
				case Bye:
					close(reqNextItem)
				}
			},
		})
	}
}

// FromChannel Converts a given channel into callbag stream.
// The stream will be termnated once channel is closed
func FromChannel[T any](ch <-chan T) Func[T] {
	var down Func[T]

	return func(opt Option[T]) {
		if opt.Type != Hello {
			panic("Expected Hello")
		}

		down = opt.Func
		reqNextItem := make(chan struct{}, 1)
		go func() {
			for {
				_, ok := <-reqNextItem
				if !ok {
					return
				}

				value, ok := <-ch
				if !ok {
					down(Option[T]{
						Type: Bye,
					})
					return
				}

				down(Option[T]{
					Type: Data,
					Data: value,
				})
			}
		}()

		down(Option[T]{
			Type: Hello,
			Func: func(opt Option[T]) {
				switch opt.Type {
				case Data:
					reqNextItem <- struct{}{}
				case Bye:
					close(reqNextItem)
				}
			},
		})
	}
}

// Filter this function can be used to filter out any items
// inside stream from going through the stream
func Filter[T any](cond func(value T) bool) func(Func[T]) Func[T] {
	var up Func[T]
	var down Func[T]

	return func(source Func[T]) Func[T] {
		return func(opt Option[T]) { // this func will be called by downstream
			switch opt.Type {
			case Hello:
				if down != nil {
					panic("Too many Hello calls from downstream")
				}

				down = opt.Func
				source(Option[T]{
					Type: Hello,
					Func: func(opt Option[T]) {
						switch opt.Type {
						case Hello:
							up = opt.Func
							down(Option[T]{
								Type: Hello,
								Func: func(opt Option[T]) {
									up(opt)
								},
							})
						case Data:
							if cond(opt.Data) {
								down(opt)
							} else {
								up(Option[T]{
									Type: Data,
								})
							}
						case Bye:
							down(opt)
						}
					},
				})
			case Data:
				up(Option[T]{
					Type: Data,
				})
			case Bye:
				up(Option[T]{
					Type: Bye,
					Err:  opt.Err,
				})
			}
		}
	}
}

// Map transforms one value to another
func Map[T any, E any](fn func(T) E) func(Func[T]) Func[E] {
	var up Func[T]
	var down Func[E]

	return func(source Func[T]) Func[E] {
		return func(opt Option[E]) { // this func will be called by downstream
			switch opt.Type {
			case Hello:
				if down != nil {
					panic("Too many Hello calls from downstream")
				}

				down = opt.Func
				source(Option[T]{
					Type: Hello,
					Func: func(opt Option[T]) {
						switch opt.Type {
						case Hello:
							up = opt.Func
							down(Option[E]{
								Type: Hello,
								Func: func(opt Option[E]) {
									up(Option[T]{
										Type: opt.Type,
										Err:  opt.Err,
									})
								},
							})
						case Data:
							down(Option[E]{
								Type: Data,
								Data: fn(opt.Data),
							})
						case Bye:
							down(Option[E]{
								Type: Bye,
								Err:  opt.Err,
							})
						}
					},
				})
			case Data:
				up(Option[T]{
					Type: Data,
				})
			case Bye:
				up(Option[T]{
					Type: Bye,
					Err:  opt.Err,
				})
			}

		}
	}
}

// ParallelMap is a special case of map that creates N number of goroutines
// to process the data. It put the processed data back to its original index
// location. The number of goroutines is determined by the length of input
func ParallelMap[T any, E any](fn func(T) E) func(Func[[]T]) Func[[]E] {
	process := func(values []T) []E {
		results := make([]E, len(values))

		wg := sync.WaitGroup{}
		for i, value := range values {
			wg.Add(1)
			go func(i int, value T) {
				results[i] = fn(value)
				wg.Done()
			}(i, value)
		}
		wg.Wait()

		return results
	}

	return Map(process)
}

// Group is a function that groups N number of items into a single slice item.
// This function is useful when you want to use ParallelMap to process N number and
// once the proces is done, you can use Flatten to flatten the data back to a stream of items
func Group[T any](n int) func(Func[T]) Func[[]T] {
	var up Func[T]
	var down Func[[]T]
	var i int
	buffer := make([]T, n)

	return func(source Func[T]) Func[[]T] {
		return func(opt Option[[]T]) { // this func will be called by downstream
			switch opt.Type {
			case Hello:
				if down != nil {
					panic("Too many Hello calls from downstream")
				}

				down = opt.Func
				source(Option[T]{
					Type: Hello,
					Func: func(opt Option[T]) {
						switch opt.Type {
						case Hello:
							up = opt.Func
							down(Option[[]T]{
								Type: Hello,
								Func: func(opt Option[[]T]) {
									up(Option[T]{
										Type: opt.Type,
										Err:  opt.Err,
									})
								},
							})
						case Data:
							buffer[i] = opt.Data
							i++

							if i == n {
								i = 0
								down(Option[[]T]{
									Type: Data,
									Data: buffer,
								})
							} else {
								up(Option[T]{
									Type: Data,
								})
							}
						case Bye:
							down(Option[[]T]{
								Type: Bye,
								Err:  opt.Err,
							})
						}
					},
				})
			case Data:
				up(Option[T]{
					Type: Data,
				})
			case Bye:
				up(Option[T]{
					Type: Bye,
					Err:  opt.Err,
				})
			}
		}
	}
}

// Flatten converts a stream of slices into a stream of single item
func Flatten[T any]() func(Func[[]T]) Func[T] {
	var up Func[[]T]
	var down Func[T]

	var values []T
	var i int = -1

	return func(source Func[[]T]) Func[T] {
		return func(opt Option[T]) { // this func will be called by downstream
			switch opt.Type {
			case Hello:
				if down != nil {
					panic("Too many Hello calls from downstream")
				}

				down = opt.Func
				source(Option[[]T]{
					Type: Hello,
					Func: func(opt Option[[]T]) {
						switch opt.Type {
						case Hello:
							up = opt.Func
							down(Option[T]{
								Type: Hello,
								Func: func(opt Option[T]) {
									switch opt.Type {
									case Data:
										if i == -1 {
											//request data from upstream
											up(Option[[]T]{
												Type: Data,
											})
											return
										}

										i++

										if i == len(values) {
											i = 0
											up(Option[[]T]{
												Type: Data,
											})
										} else {
											down(Option[T]{
												Type: Data,
												Data: values[i],
											})
										}
									}
								},
							})
						case Data:
							i = 0
							values = opt.Data
							down(Option[T]{
								Type: Data,
								Data: values[i],
							})
						case Bye:
							down(Option[T]{
								Type: Bye,
								Err:  opt.Err,
							})
						}
					},
				})
			case Bye:
				up(Option[[]T]{
					Type: Bye,
					Err:  opt.Err,
				})
			}
		}
	}
}

// Take gets the N number of items and stop the stream.
// This function is usful if you are dealing with large or infinite streams
func Take[T any](n int) func(Func[T]) Func[T] {
	var up Func[T]
	var down Func[T]
	var i int

	return func(source Func[T]) Func[T] {
		return func(opt Option[T]) { // this func will be called by downstream
			switch opt.Type {
			case Hello:
				if down != nil {
					panic("Too many Hello calls from downstream")
				}

				down = opt.Func
				source(Option[T]{
					Type: Hello,
					Func: func(opt Option[T]) {
						switch opt.Type {
						case Hello:
							up = opt.Func
							down(Option[T]{
								Type: Hello,
								Func: func(opt Option[T]) {
									up(Option[T]{
										Type: opt.Type,
										Err:  opt.Err,
									})
								},
							})
						case Data:
							if !(i < n) {
								up(Option[T]{
									Type: Bye,
								})
								down(Option[T]{
									Type: Bye,
								})
								return
							}
							i++
							down(Option[T]{
								Type: Data,
								Data: opt.Data,
							})
						case Bye:
							down(Option[T]{
								Type: Bye,
								Err:  opt.Err,
							})
						}
					},
				})
			case Data:
				up(Option[T]{
					Type: Data,
				})
			case Bye:
				up(Option[T]{
					Type: Bye,
					Err:  opt.Err,
				})
			}

		}
	}
}

// ForEach is a sink function and should be used as the last function in a stream
// once the stream is terminated, the done argument will be set to true
func ForEach[T any](fn func(item T, done bool)) func(Func[T]) {
	return func(source Func[T]) {

		var next Func[T]

		source(Option[T]{
			Type: Hello,
			Func: func(opt Option[T]) {
				switch opt.Type {
				case Hello:
					next = opt.Func
					next(Option[T]{
						Type: Data,
					})
				case Data:
					fn(opt.Data, false)
					next(Option[T]{
						Type: Data,
					})
				case Bye:
					fn(opt.Data, true)
				}
			},
		})
	}
}

// ToChannel is a special ForEach function that push all the items to a given channel.
// Once the stream is terminated, This fucntion will automatically close the channel
func ToChannel[T any](ch chan<- T) func(Func[T]) {
	return ForEach(func(value T, done bool) {
		if done {
			close(ch)
			return
		}

		ch <- value
	})
}

// Pipe is a helper function that let's combined all types of callbag functions.
// There are number of varity of this function, `Pipe2`, `Pipe3`, `Pipe4`, ..., `Pipe10`.
// The N represents the number of callbag function you can combined. The first argument
// must be a source, and the last argument must be a sink function. The rest of the arguemnts
// can be set with all modifier callbag functions.

func Pipe2[A any](
	src Func[A],
	sink func(Func[A])) {
	sink(src)
}

func Pipe3[A any, B any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	sink func(Func[B])) {
	sink(m1(src))
}

func Pipe4[A any, B any, C any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	sink func(Func[C])) {
	sink(m2(m1(src)))
}

func Pipe5[A any, B any, C any, D any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	sink func(Func[D])) {
	sink(m3(m2(m1(src))))
}

func Pipe6[A any, B any, C any, D any, E any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	sink func(Func[E])) {
	sink(m4(m3(m2(m1(src)))))
}

func Pipe7[A any, B any, C any, D any, E any, F any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	sink func(Func[F])) {
	sink(m5(m4(m3(m2(m1(src))))))
}

func Pipe8[A any, B any, C any, D any, E any, F any, G any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
	sink func(Func[G])) {
	sink(m6(m5(m4(m3(m2(m1(src)))))))
}

func Pipe9[A any, B any, C any, D any, E any, F any, G any, H any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
	m7 func(Func[G]) Func[H],
	sink func(Func[H])) {
	sink(m7(m6(m5(m4(m3(m2(m1(src))))))))
}

func Pipe10[A any, B any, C any, D any, E any, F any, G any, H any, I any](
	src Func[A],
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
	m7 func(Func[G]) Func[H],
	m8 func(Func[H]) Func[I],
	sink func(Func[I])) {
	sink(m8(m7(m6(m5(m4(m3(m2(m1(src)))))))))
}

// Similar to Pipe, Compose is a helper function to compose multiple modifier callbag functions.
// The compose functions can not accept either of soucre or sink functions

func Compose2[A any, B any, C any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
) func(Func[A]) Func[C] {
	return func(in Func[A]) Func[C] {
		return m2(m1(in))
	}
}

func Compose3[A any, B any, C any, D any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
) func(Func[A]) Func[D] {
	return func(in Func[A]) Func[D] {
		return m3(m2(m1(in)))
	}
}

func Compose4[A any, B any, C any, D any, E any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
) func(Func[A]) Func[E] {
	return func(in Func[A]) Func[E] {
		return m4(m3(m2(m1(in))))
	}
}

func Compose5[A any, B any, C any, D any, E any, F any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
) func(Func[A]) Func[F] {
	return func(in Func[A]) Func[F] {
		return m5(m4(m3(m2(m1(in)))))
	}
}

func Compose6[A any, B any, C any, D any, E any, F any, G any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
) func(Func[A]) Func[G] {
	return func(in Func[A]) Func[G] {
		return m6(m5(m4(m3(m2(m1(in))))))
	}
}

func Compose7[A any, B any, C any, D any, E any, F any, G any, H any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
	m7 func(Func[G]) Func[H],
) func(Func[A]) Func[H] {
	return func(in Func[A]) Func[H] {
		return m7(m6(m5(m4(m3(m2(m1(in)))))))
	}
}

func Compose8[A any, B any, C any, D any, E any, F any, G any, H any, I any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
	m7 func(Func[G]) Func[H],
	m8 func(Func[H]) Func[I],
) func(Func[A]) Func[I] {
	return func(in Func[A]) Func[I] {
		return m8(m7(m6(m5(m4(m3(m2(m1(in))))))))
	}
}

func Compose9[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
	m7 func(Func[G]) Func[H],
	m8 func(Func[H]) Func[I],
	m9 func(Func[I]) Func[J],
) func(Func[A]) Func[J] {
	return func(in Func[A]) Func[J] {
		return m9(m8(m7(m6(m5(m4(m3(m2(m1(in)))))))))
	}
}

func Compose10[A any, B any, C any, D any, E any, F any, G any, H any, I any, J any, K any](
	m1 func(Func[A]) Func[B],
	m2 func(Func[B]) Func[C],
	m3 func(Func[C]) Func[D],
	m4 func(Func[D]) Func[E],
	m5 func(Func[E]) Func[F],
	m6 func(Func[F]) Func[G],
	m7 func(Func[G]) Func[H],
	m8 func(Func[H]) Func[I],
	m9 func(Func[I]) Func[J],
	m10 func(Func[J]) Func[K],
) func(Func[A]) Func[K] {
	return func(in Func[A]) Func[K] {
		return m10(m9(m8(m7(m6(m5(m4(m3(m2(m1(in))))))))))
	}
}
