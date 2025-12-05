package utils 

type Mux[fn any] struct {
	table map[string]fn
}

func NewMux[fn any]() *Mux[fn] {
	return &Mux[fn]{
		table: make(map[string]fn),
	}
}

func (mux *Mux[fn]) Register(command string, handler fn) {
	mux.table[command] = handler
}

func (mux *Mux[fn]) Handle(command string) (fn, bool) {
	handler, exists := mux.table[command]
	return handler, exists
}

