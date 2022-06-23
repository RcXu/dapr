package mytest

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/dapr/components-contrib/mytest"
	"github.com/dapr/dapr/pkg/components"
)

type Mytest struct {
	Name          string
	FactoryMethod func() mytest.Mytest
}

func New(name string, factoryMethod func() mytest.Mytest) Mytest {
	return Mytest{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

type Registry interface {
	Register(components ...Mytest)
	Create(name, version string) (mytest.Mytest, error)
}

type mytestRegistry struct {
	mytests map[string]func() mytest.Mytest
}

func NewRegistry() Registry {
	return &mytestRegistry{
		mytests: map[string]func() mytest.Mytest{},
	}
}

func (t *mytestRegistry) Register(components ...Mytest) {
	for _, component := range components {
		t.mytests[createFullName(component.Name)] = component.FactoryMethod
	}
}

func (t *mytestRegistry) Create(name, version string) (mytest.Mytest, error) {
	if method, ok := t.getMytest(name, version); ok {
		return method(), nil
	}
	return nil, errors.Errorf("couldn't find Mytest %s/%s", name, version)
}

func (t *mytestRegistry) getMytest(name, version string) (func() mytest.Mytest, bool) {
	nameLower := strings.ToLower(name)
	versionLower := strings.ToLower(version)
	mytestFn, ok := t.mytests[nameLower+"/"+versionLower]
	if ok {
		return mytestFn, true
	}
	if components.IsInitialVersion(versionLower) {
		mytestFn, ok = t.mytests[nameLower]
	}
	return mytestFn, ok
}

func createFullName(name string) string {
	return strings.ToLower("mytest." + name)
}
