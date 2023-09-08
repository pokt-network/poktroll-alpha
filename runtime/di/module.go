package di

type Module interface {
	Resolve(*Injector, *[]string)
	Start() error
	CascadeStart() error
}
