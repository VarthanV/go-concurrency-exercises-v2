package generators

// Generate is a generic function that can be used to generate stream of any type T.
func Generate[T any](done chan any, values ...T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		for _, value := range values {
			select {
			case <-done:
				return
			case stream <- value:
				// Value sent to the stream
			}
		}
	}()

	return stream
}

func GenerateFromFunc[T any](done chan any, fn func() []T) <-chan T {
	stream := make(chan T)
	go func() {
		defer close(stream)
		values := fn()
		for _, v := range values {
			select {
			case <-done:
				return
			case stream <- v:
			}
		}
	}()
	return stream
}
