# gollection

A generic generic collection library based on go's generic implementation.

> This library relies on the generic features of go 1.18 and is not yet stable, we will make the api stable after 1.18.

## Core Interfaces

```go
type Iterable[T any] interface {
	Iter() Iterator[T]
}

type Iterator[T any] interface {
	Next() Option[T]
}
```

gollection applies the iterator pattern design, the core interface consists of Iterator and Iterable.

Iterator is responsible for providing iterative functionality, Iterator is unidirectional and lazy, each call to next will only return one result.

Iterable is responsible for providing Iterator, the implementation type determines whether the provided Iterator is reusable.

The inert traversal feature allows the combination of higher-order functions without significant overhead and can provide a richer combination of functions.

Here is a simple example of direct traversal:

```go
func printAll[T any](it Iterable[T]) {
	var iter = it.Iter()
	for v, ok := iter.Next().Get(); ok; v, ok = iter.Next().Get() {
		println(v)
	}
}
```

## Streaming operations

gollection provides a rich set of stream manipulation functions that can be used in combination with any Iterable.

Here is an example of a simple combination used:

```go
func foo() {
	show := func(i int) {
		println(i)
	}
	even := func(i int) bool {
		return i%2 == 0
	}
	square := func(i int) int {
		return i * i
	}
	ForEach(show, Mapper(square, Filter(even, ToSliceIter([]int{1, 2, 3, 4, 5, 6, 7}))))
    // Result:
    // 4
    // 16
    // 36
}
```

### Transform Iterable

A series of conversion functions are provided to process one Iterable conversion to another Iterable. these conversions are not executed immediately and only act one at a time when iterating.

```go
func Indexer[T any](it Iterator[T]) Iterator[Pair[int, T]]
func Mapper[T any, R any](transform func(T) R, it Iterator[T]) Iterator[R]
func Filter[T any](predecate func(T) bool, it Iterator[T]) Iterator[T]
func Limit[T any](count int, it Iterator[T]) Iterator[T]
func Skip[T any](count int, it Iterator[T]) Iterator[T]
func Step[T any](count int, it Iterator[T]) Iterator[T]
func Concat[T any](left Iterator[T], right Iterator[T]) Iterator[T]
```

### Terminal Iterable

A set of functions that evaluate the Iterable and are executed immediately.

```go
func Contains[T comparable](target T, it Iterator[T]) bool
func Sum[T Number](it Iterator[T]) T
func Product[T Number](it Iterator[T]) T
func Average[T Number](it Iterator[T]) float64
func Count[T any](it Iterator[T]) int
func Max[T Number](it Iterator[T]) T
func Min[T Number](it Iterator[T]) T
func ForEach[T any](action func(T), it Iterator[T])
func AllMatch[T any](predicate func(T) bool, it Iterator[T]) bool
func NoneMatch[T any](predicate func(T) bool, it Iterator[T]) bool
func AnyMatch[T any](predicate func(T) bool, it Iterator[T]) bool
func First[T any](it Iterator[T]) Option[T]
func Last[T any](it Iterator[T]) Option[T]
func At[T any](index int, it Iterator[T]) Option[T]
func Reduce[T any, R any](initial R, operation func(R, T) R, it Iterator[T]) R
func Fold[T any, R any](initial R, operation func(T, R) R, it Iterator[T]) R
```

## ToString and ToSlice

In order to make go's native string and slice also iterable, we have introduced `ToSlice` and `ToString` to make these two types implement the interface.

```go
var str = "Hello, world!"
var sli = []int{1, 2, 3}
Count(ToString(str).Iter()) // 13
Count(ToSlice(sli).Iter()) // 3
```

We also provide the version that gets the iterator directly.

```go
Count(ToStringIter(str)) // 13
Count(ToSliceIter(sli)) // 3
```

## Collection

We define a unified collection type interface to describe more information than iterators to facilitate performance optimization.

```go
type Collection[T any] interface {
	Iterable[T]

	Size() int
	IsEmpty() bool
	ToSlice() []T
}
```

## List

We provide the List interface to unify the description of ordered sequences, and provide `ArrayList` and `LinkedList` as its implementation types.

```go
type List[T any] interface {
	Collection[T]

	Get(index int) T
	Set(index int, newElement T) T
	GetAndSet(index int, set func(oldElement T) T) Pair[T, T]
	TryGet(index int) Option[T]
	TrySet(index int, newElement T) Option[T]

	Prepend(element T)
	PrependAll(elements Collection[T])
	Append(element T)
	AppendAll(elements Collection[T])
	Insert(index int, element T)
	InsertAll(index int, elements Collection[T])
	Remove(index int) T
	Clear()
}
```

### Map

We provide the Map interface to unify the description of the mapping type, and provide `HashMap` as its implementation type.

```go
type Map[K any, V any] interface {
	Collection[Pair[K, V]]

	Get(key K) V
	Put(key K, value V) Option[V]
	PutAll(elements Collection[Pair[K, V]])
	GetAndPut(key K, set func(oldValue Option[V]) V) Pair[V, Option[V]]
	TryGet(key K) Option[V]

	Remove(key K) Option[V]
	Contains(key K) bool
	Clear()
}
```

### Set

We provide the Set interface to describe the element-unique collection type, and we provide `HashSet` as its implementation type.

```go
type Set[T any] interface {
	Collection[T]

	Put(element T) bool
	PutAll(elements Collection[T])

	Remove(element T) bool
	Contains(element T) bool
	ContainsAll(elements Collection[T]) bool
	Clear()
}
```

### Stack

We provide the Stack interface to describe the stack data structure and provide `ArrayStack` and `LinkedStack` as its implementation types.

```go
type Stack[T any] interface {
	Collection[T]

	Push(element T)
	Pop() T
	Peek() T
	TryPop() Option[T]
	TryPeek() Option[T]
}
```

### Tuple

We have also introduced several convenient tuple types for use, and indeed gollection uses them as well. Including `Void`, `Pair`, `Triple`.

```go
type Void struct{}

func PairOf[T1 any, T2 any](f T1, s T2) Pair[T1, T2]
func (a Pair[T1, T2]) Get() (T1, T2)

func TripleOf[T1 any, T2 any, T3 any](f T1, s T2, t T3) Triple[T1, T2, T3]
func (a Triple[T1, T2, T3]) Get() (T1, T2, T3)
```

### Union

We have also introduced several convenient union types for use, which are actually used by gollection. Including `Option` and `Result`.

```go
func Some[T any](a T) Option[T]
func None[T any]() Option[T]
func (a Option[T]) Get() (value T, ok bool)

func Ok[T any](a T) Result[T]
func Err[T any](a error) Result[T]
func (a Result[T]) Get() (value T, err error)
```