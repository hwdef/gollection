package stack

import . "github.com/kulics/gollection"

type Stack[T any] interface {
	Collection[T]

	Push(element T)
	Pop() T
	Peek() T
	TryPop() Option[T]
	TryPeek() Option[T]
}
