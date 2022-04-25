package callbag

import (
	"time"
)

type Op int

const (
	Greets Op = iota
	Data
	Bye
)

type Func[T any] func(op Op, start Func[T], data T, end func(error))
type Modifier[T any, E any] func(Func[T]) Func[E]

func Interval(period time.Duration) Func[int] {
	return func(op Op, start Func[int], data int, end func(error)) {
		if start == nil {
			return
		}

		i := 0
		ticker := time.NewTicker(period)
		done := make(chan struct{})

		go func() {
			for {
				select {
				case <-ticker.C:
					start(Data, nil, i, nil)
					i++
				case <-done:
					ticker.Stop()
					return
				}
			}
		}()

		start(Greets, func(op Op, start Func[int], data int, end func(error)) {
			if end != nil {
				close(done)
			}
		}, 0, nil)
	}
}

func FromSlice[T any](s []T) Func[T] {
	i := -1
	n := len(s)

	var emptyT T

	return func(op Op, down Func[T], data T, end func(error)) {
		if op != Greets {
			return
		}

		down(Greets, func(op Op, _ Func[T], _ T, _ func(error)) {
			if op == Data {
				for i < n {
					down(Data, nil, s[i], nil)
					i++
				}
				down(Bye, nil, emptyT, nil)
			}
		}, data, end)
	}
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

func FromRange[T number](a, b, step T) Func[T] {
	return func(op Op, down Func[T], data T, end func(error)) {
		if op != Greets {
			return
		}

		down(Greets, func(op Op, _ Func[T], _ T, _ func(error)) {
			if op == Data {
				for a < b {
					down(Data, nil, a, nil)
					a += step
				}
				down(Bye, nil, 0, nil)
			}
		}, data, end)
	}
}

func FromChannel[T any](ch <-chan T) Func[T] {
	var emptyT T

	return func(op Op, down Func[T], data T, end func(error)) {
		if op != Greets {
			return
		}

		start := make(chan struct{})
		done := make(chan struct{})

		go func() {
			<-start

			for {
				select {
				case <-done:
					down(Bye, nil, data, nil)
					return
				case data, ok := <-ch:
					if !ok {
						down(Bye, nil, data, nil)
						return
					}
					down(Data, nil, data, nil)
				}
			}
		}()

		down(Greets, func(op Op, _ Func[T], _ T, _ func(error)) {
			if op == Bye {
				close(done)
				return
			}

			if op == Data {
				close(start)
				return
			}
		}, emptyT, nil)
	}
}

func Filter[T any](cond func(value T) bool) Modifier[T, T] {
	var emptyT T

	return func(source Func[T]) Func[T] {
		return func(op Op, down Func[T], _ T, _ func(error)) {
			if op != Greets {
				return
			}

			source(Greets, func(op Op, up Func[T], data T, end func(error)) {
				if op == Data {
					if cond(data) {
						down(Data, nil, data, nil)
					}
					return
				}

				down(op, up, data, end)
			}, emptyT, nil)
		}
	}
}

func Map[T any, E any](fn func(T) E) Modifier[T, E] {
	var emptyT T
	var emptyE E

	return func(source Func[T]) Func[E] {
		return func(op Op, down Func[E], _ E, _ func(error)) {
			if op != Greets {
				return
			}

			source(Greets, func(op Op, up Func[T], data T, end func(error)) {
				switch op {
				case Greets:
					up(Data, nil, emptyT, nil)
				case Data:
					down(Data, nil, fn(data), nil)
				case Bye:
					down(Bye, nil, emptyE, nil)
				}
			}, emptyT, nil)
		}
	}
}

func Reduce[T any, E any](fn func(base E, value T) E, base E) Modifier[T, E] {
	var emptyT T
	var emptyE E

	value := base

	return func(source Func[T]) Func[E] {
		return func(op Op, down Func[E], _ E, _ func(error)) {
			if op != Greets {
				return
			}

			source(Greets, func(op Op, up Func[T], data T, end func(error)) {
				switch op {
				case Greets:
					up(Data, nil, emptyT, nil)
				case Data:
					value = fn(value, data)
				case Bye:
					down(Data, nil, value, nil)
					down(Bye, nil, emptyE, nil)
				}
			}, emptyT, nil)
		}
	}
}

func Take[T any](n int) Modifier[T, T] {
	var emptyT T
	var talkup Func[T]

	return func(source Func[T]) Func[T] {
		return func(op Op, down Func[T], _ T, _ func(error)) {
			if op != Greets {
				return
			}

			source(Greets, func(op Op, up Func[T], data T, end func(error)) {
				switch op {
				case Greets:
					talkup = up
					up(Data, nil, emptyT, nil)
				case Data:
					n--
					if n >= 0 {
						down(Data, nil, data, nil)
					}

					if n == 0 {
						talkup(Bye, nil, emptyT, nil)
					}
				case Bye:
					down(Bye, nil, emptyT, nil)
				}
			}, emptyT, nil)
		}
	}
}

func ForEach[T any](fn func(T, bool)) func(Func[T]) {
	var emptyT T

	return func(source Func[T]) {
		source(Greets, func(op Op, up Func[T], data T, end func(error)) {
			switch op {
			case Greets:
				up(Data, nil, emptyT, nil) // ready, sent us everything
			case Data:
				fn(data, true)
			case Bye:
				fn(data, false)
			}
		}, emptyT, nil)
	}
}

func Pipe2[A any](
	src Func[A],
	sink func(Func[A])) {
	sink(src)
}

func Pipe3[A any, B any](
	src Func[A],
	m1 Modifier[A, B],
	sink func(Func[B])) {
	sink(m1(src))
}

func Pipe4[A any, B any, C any](
	src Func[A],
	m1 Modifier[A, B],
	m2 Modifier[B, C],
	sink func(Func[C])) {
	sink(m2(m1(src)))
}

func Pipe5[A any, B any, C any, D any](
	src Func[A],
	m1 Modifier[A, B],
	m2 Modifier[B, C],
	m3 Modifier[C, D],
	sink func(Func[D])) {
	sink(m3(m2(m1(src))))
}

func Pipe6[A any, B any, C any, D any, E any](
	src Func[A],
	m1 Modifier[A, B],
	m2 Modifier[B, C],
	m3 Modifier[C, D],
	m4 Modifier[D, E],
	sink func(Func[E])) {
	sink(m4(m3(m2(m1(src)))))
}

func Pipe7[A any, B any, C any, D any, E any, F any](
	src Func[A],
	m1 Modifier[A, B],
	m2 Modifier[B, C],
	m3 Modifier[C, D],
	m4 Modifier[D, E],
	m5 Modifier[E, F],
	sink func(Func[F])) {
	sink(m5(m4(m3(m2(m1(src))))))
}

func Pipe8[A any, B any, C any, D any, E any, F any, G any](
	src Func[A],
	m1 Modifier[A, B],
	m2 Modifier[B, C],
	m3 Modifier[C, D],
	m4 Modifier[D, E],
	m5 Modifier[E, F],
	m6 Modifier[F, G],
	sink func(Func[G])) {
	sink(m6(m5(m4(m3(m2(m1(src)))))))
}

func Pipe9[A any, B any, C any, D any, E any, F any, G any, H any](
	src Func[A],
	m1 Modifier[A, B],
	m2 Modifier[B, C],
	m3 Modifier[C, D],
	m4 Modifier[D, E],
	m5 Modifier[E, F],
	m6 Modifier[F, G],
	m7 Modifier[G, H],
	sink func(Func[H])) {
	sink(m7(m6(m5(m4(m3(m2(m1(src))))))))
}

func Pipe10[A any, B any, C any, D any, E any, F any, G any, H any, I any](
	src Func[A],
	m1 Modifier[A, B],
	m2 Modifier[B, C],
	m3 Modifier[C, D],
	m4 Modifier[D, E],
	m5 Modifier[E, F],
	m6 Modifier[F, G],
	m7 Modifier[G, H],
	m8 Modifier[H, I],
	sink func(Func[I])) {
	sink(m8(m7(m6(m5(m4(m3(m2(m1(src)))))))))
}
