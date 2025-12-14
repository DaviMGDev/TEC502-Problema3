package utils 

type Comparable interface {
	Equals(other Comparable) bool
}

type String string 

func (s String) Equals(other Comparable) bool {
	otherString, ok := other.(String)
	if !ok {
		return false
	}
	return s == otherString
}
