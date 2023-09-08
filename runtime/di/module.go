package di

type Module interface {
	Hydrate(*Injector, *[]string)
	Start() error
	CascadeStart() error
}
