package utils

import ()

type Map[K comparable, V any] interface {
	Get(key K) (V, bool)
	Set(key K, value V)
	Delete(key K)
	Has(key K) bool
}
