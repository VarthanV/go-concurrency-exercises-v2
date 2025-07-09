package main

type Slice[T any] struct {
	s []T
}

func New[T any](length int) Slice[T] {
	defaultLen := 1
	if length > defaultLen {
		defaultLen = length
	}
	return Slice[T]{
		s: make([]T, defaultLen),
	}
}

func (s *Slice[T]) Append(value T) {
	s.s = append(s.s, value)
}
func (s *Slice[T]) Get(index int) T {
	if index < 0 || index >= len(s.s) {
		panic("index out of bounds")
	}
	return s.s[index]
}

func (s *Slice[T]) InsertInIndex(index int, value T) {
	if index < 0 || index >= len(s.s) {
		panic("index out of bounds")
	}
	s.s = append(append(s.s[:index], value), s.s[index+1:]...)
}

func (s *Slice[T]) Pop() T {
	if len(s.s) == 0 {
		panic("cannot pop from a empty slice")
	}
	length := len(s.s) - 1
	val := s.s[length]
	s.s = s.s[:length]
	return val
}

func (s *Slice[T]) RemoveFromIndex(index int) {
	if index < 0 || index >= len(s.s) {
		panic("index out of bounds")
	}
	s.s = append(s.s[:index], s.s[index+1:]...)
}

func (s *Slice[T]) SwapIndex(from, to int) {
	if (from < 0 || from >= len(s.s)) || (to < 0 || from >= len(s.s)) {
		panic("index out of bounds")
	}
	s.s[from], s.s[to] = s.s[to], s.s[from]
}

func (s *Slice[T]) Slice(from, to int) []T {
	if (from < 0 || from >= len(s.s)) || (to < 0 || from >= len(s.s)) {
		panic("index out of bounds")
	}
	return s.s[from:to]
}

func (s *Slice[T]) All() []T {
	return s.s
}
