package key_value

type Value struct {
	Value  any
	Exists bool
}

type Storage interface {
	Get(key string) Value

	Set(key string, value any)

	Delete(key string)
}
