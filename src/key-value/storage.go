package key_value

type Value struct {
	Value  any
	Exists bool
}

type Storage interface {
	Get(key string) (Value, error)

	Set(key string, value any) error

	Delete(key string) error
}
