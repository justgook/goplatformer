package util

type Stack[T any] []T

func (q *Stack[T]) Push(n T) {
	*q = append(*q, n)
}

func (q *Stack[T]) Pop() T {
	x := q.Len() - 1
	n := (*q)[x]
	*q = (*q)[:x]

	return n
}
func (q *Stack[T]) Len() int {
	return len(*q)
}
