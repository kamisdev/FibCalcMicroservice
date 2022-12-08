package stack

type Stack []float32

func (s *Stack) IsEmpty() bool {
	return len(*s) == 0
}

func (s *Stack) Push(v float32) {
	*s = append(*s, v)
}

func (s *Stack) Pop() (float32, bool) {
	if s.IsEmpty() {
		return -1.0, false
	}

	item := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return item, true
}
