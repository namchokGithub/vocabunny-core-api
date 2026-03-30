package domain

type EntityField[T any] struct {
	Value T
	Set   bool
}

func NewEntityField[T any](value T) EntityField[T] {
	return EntityField[T]{
		Value: value,
		Set:   true,
	}
}
