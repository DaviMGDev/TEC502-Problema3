package utils

import ()

type Mux[fn any] struct {
	Map[string, fn]
	defaultFn	fn
}

func (m *Mux[fn]) Get(key string) (fn, bool) {
	if !m.Has(key) {
		return m.defaultFn, false
	}
	return m.Map.Get(key)
}
