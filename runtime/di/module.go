package di

type Module interface {
	Resolve(*Injector, *[]string)
	Start() error
	CascadeStart() error
}

type ModuleInternals[Deps any] struct {
	deps *Deps
}

func (m *ModuleInternals[Deps]) Deps() *Deps {
	return m.deps
}

func (m *ModuleInternals[Deps]) ResolveDeps(deps *Deps) {
	m.deps = deps
}

type Uninjectable interface {
	Uninjectable()
}
