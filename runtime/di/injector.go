package di

import (
	"fmt"
	"strings"
)

type Injector struct {
	injections      map[string]any
	hydratedModules map[string]bool
	sealed          bool
}

var InjectorToken = NewInjectionToken[Injector]("Injector")

func NewInjector() *Injector {
	return &Injector{map[string]any{}, map[string]bool{}, false}
}

func HydrateMain[V Module](token *InjectionToken[V], injector *Injector) V {
	path := &[]string{}
	result := Hydrate[V](token, injector, path)
	injector.sealed = true
	return result
}

func Hydrate[V any](token *InjectionToken[V], injector *Injector, path *[]string) V {
	for _, p := range *path {
		if p == token.Id() {
			panic(fmt.Sprintf("Circular dependency detected [ %s -> %s ]", strings.Join(*path, " -> "), token.Id()))
		}
	}
	*path = append(*path, token.Id())
	if injector.injections[token.Id()] == nil {
		panic(fmt.Sprintf("Injection not provided [ %s ]", strings.Join(*path, " -> ")))
	}

	value := injector.injections[token.Id()]

	if module, ok := value.(Module); ok {
		if !injector.hydratedModules[token.Id()] {
			module.Hydrate(injector, path)
			injector.hydratedModules[token.Id()] = true
		}
	}

	if castedValue, ok := value.(V); ok {
		if len(*path) > 0 {
			*path = (*path)[:len(*path)-1]
		}
		return castedValue
	} else {
		panic(fmt.Sprintf("Injection type mismatch [ %s -> %s ]", strings.Join(*path, " -> "), token.Id()))
	}
}

func Provide[V any](token *InjectionToken[V], value V, injector *Injector) {
	if injector.sealed {
		panic("Injector sealed")
	}
	injector.injections[token.Id()] = value
}

func Get[V any](token *InjectionToken[V], injector *Injector) V {
	if injector.injections[token.Id()] == nil {
		panic(fmt.Sprintf("Injection not provided %s", token.Id()))
	}

	return injector.injections[token.Id()].(V)
}
