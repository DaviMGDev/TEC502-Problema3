package utils 

import (
	// "sync"
)

type List[T any] interface {
	Append(item T)
	Remove(item T) bool
	Contains(item T) bool
	Size() int
	Get(index int) (T, bool)
	Search(item T) int
}
