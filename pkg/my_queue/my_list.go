package my_list

type MyList[T any] struct {
	list     []T
	len      int
	capacity int
}

func NewMyList[T any](capacity int) MyList[T] {
	return MyList[T]{capacity: capacity, len: 0, list: make([]T, 0)}
}

func (q *MyList[T]) GetIndex(index int) T {
	return q.list[index]
}

func (q *MyList[T]) SetIndex(index int, value T) {
	newSlice := make([]T, len(q.list)+1)

	copy(newSlice[:index], q.list[:index])
	newSlice[index] = value
	copy(newSlice[index+1:], q.list[index:])

	if len(newSlice) > q.capacity {
		newSlice = newSlice[:q.capacity]
	}

	q.list = newSlice
}

func (q *MyList[T]) RemoveIndex(index int) {
	newSlice := make([]T, len(q.list)-1)

	copy(newSlice[:index], q.list[:index])
	copy(newSlice[index:], q.list[index+1:])

	q.list = newSlice

}
