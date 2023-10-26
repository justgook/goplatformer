package util

type Queue[T any] []T

func (q *Queue[T]) Push(n T) {
	*q = append(*q, n)
}

func (q *Queue[T]) Pop() T {
	var n T
	if len(*q) < 1 {
		return n
	}
	n = (*q)[0]
	*q = (*q)[1:]

	return n
}

func (q *Queue[T]) Len() int {
	return len(*q)
}
