package key_value

type Value struct {
	Value  any
	Exists bool
}

type Storage interface {
	Get(key string) Value

	Set(key string, value any)

	CompareAndSet(key string, oldValue any, newValue any) (bool, error)

	Delete(key string)

	AddElement(key string, subKey string, value any)
}
