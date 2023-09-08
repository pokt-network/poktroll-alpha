package di_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"poktroll/runtime/di"
)

var moduleDepInjectionToken = di.NewInjectionToken[DependencyModule]("depModule")
var mainModuleInjectionToken = di.NewInjectionToken[MainModule]("module")
var configInjectionToken = di.NewInjectionToken[int]("config")

type DependencyModule interface {
	di.Module
	DoThis(string) string
}

type depModuleImpl struct {
	prefix string
}

func (m *depModuleImpl) Module() DependencyModule                      { return m }
func (m *depModuleImpl) Hydrate(injector *di.Injector, path *[]string) {}
func (m *depModuleImpl) Start() error                                  { return nil }
func (m *depModuleImpl) CascadeStart() error                           { return nil }
func (m *depModuleImpl) DoThis(s string) string {
	return fmt.Sprintf("%s%s", m.prefix, s)
}

type MainModule interface {
	di.Module
	DoThat(int) int
}

type mainModuleImpl struct {
	timeout   int
	moduleDep DependencyModule
}

func (m *mainModuleImpl) Hydrate(injector *di.Injector, path *[]string) {
	m.timeout = di.Hydrate(configInjectionToken, injector, path)
	m.moduleDep = di.Hydrate(moduleInjectionToken, injector, path)
}

func (m *mainModuleImpl) Module() MainModule { return m }
func (m *mainModuleImpl) Start() error       { return nil }
func (m *mainModuleImpl) CascadeStart() error {
	if err := m.moduleDep.CascadeStart(); err != nil {
		return err
	}
	return m.Start()
}
func (m *mainModuleImpl) DoThat(n int) int { return n }

func Test_DI_Works(t *testing.T) {
	injector := di.NewInjector()
	di.Provide(mainModuleInjectionToken, (&mainModuleImpl{}).Module(), injector)
	di.Provide(moduleDepInjectionToken, (&depModuleImpl{}).Module(), injector)
	di.Provide(configInjectionToken, 10, injector)

	mainMod := di.HydrateMain(mainModuleInjectionToken, injector)
	cfg := di.Get(configInjectionToken, injector)
	mainMod.DoThat(cfg)
	dep := di.Get(moduleDepInjectionToken, injector)

	assert.Equal(t, 10, cfg)
	assert.Nil(t, mainMod.Start())
	assert.Equal(t, 10, mainMod.DoThat(cfg))
	assert.Equal(t, "hello", dep.DoThis("hello"))
}

func Test_DI_MissingDependency(t *testing.T) {
	injector := di.NewInjector()
	di.Provide(mainModuleInjectionToken, (&mainModuleImpl{}).Module(), injector)
	di.Provide(configInjectionToken, 10, injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()

	di.HydrateMain(mainModuleInjectionToken, injector)
}

type circDepModuleImpl struct {
	moduleDeps MainModule
}

func (m *circDepModuleImpl) Module() DependencyModule { return m }
func (m *circDepModuleImpl) Start() error             { return nil }
func (m *circDepModuleImpl) CascadeStart() error {
	if err := m.moduleDeps.CascadeStart(); err != nil {
		return err
	}
	return m.Start()
}
func (m *circDepModuleImpl) DoThis(s string) string { return s }
func (m *circDepModuleImpl) Hydrate(injector *di.Injector, path *[]string) {
	m.moduleDeps = di.Hydrate(mainModuleInjectionToken, injector, path)
}

func Test_DI_CircularDependencies(t *testing.T) {

	injector := di.NewInjector()
	di.Provide(mainModuleInjectionToken, (&mainModuleImpl{}).Module(), injector)
	di.Provide(moduleDepInjectionToken, (&circDepModuleImpl{}).Module(), injector)
	di.Provide(configInjectionToken, 10, injector)

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic")
		}
	}()

	di.HydrateMain(mainModuleInjectionToken, injector)
}
