package di

type InjectionToken[Value any] interface {
	Id() string
}

type injectionToken[Value any] struct {
	name string
}

func (t *injectionToken[Value]) Id() string {
	return t.name
}

func NewInjectionToken[Value any](name string) InjectionToken[Value] {
	return &injectionToken[Value]{name}
}
