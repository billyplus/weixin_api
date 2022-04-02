package weixin_api

// A simple stack
type stack[T any] struct {
	data []T
}

func (s *stack[T]) empty() bool {
	return len(s.data) == 0
}

func (s *stack[T]) push(value T) {
	s.data = append(s.data, value)
}

func (s *stack[T]) pop() T {
	n := len(s.data) - 1
	value := s.data[n]
	s.data = s.data[:n]
	return value
}

func (s *stack[T]) peek() T {
	return s.data[len(s.data)-1]
}
