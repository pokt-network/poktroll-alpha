package di

type InjectionToken[Value any] struct {
	name string
}

func (t *InjectionToken[Value]) Id() string {
	return t.name
}

func NewInjectionToken[Value any](name string) *InjectionToken[Value] {
	return &InjectionToken[Value]{name}
}
