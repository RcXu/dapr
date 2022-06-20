package logstorage

import (
	"strings"

	"github.com/dapr/components-contrib/logstorage"
	"github.com/dapr/dapr/pkg/components"
	"github.com/pkg/errors"
)

type Logstorage struct {
	Name          string
	FactoryMethod func() logstorage.Logstorage
}

func New(name string, factoryMethod func() logstorage.Logstorage) Logstorage {
	return Logstorage{
		Name:          name,
		FactoryMethod: factoryMethod,
	}
}

type Registry interface {
	Register(components ...Logstorage)
	Create(name, version string) (logstorage.Logstorage, error)
}

type logstorageRegistry struct {
	logstorages map[string]func() logstorage.Logstorage
}

func NewRegistry() Registry {
	return &logstorageRegistry{
		logstorages: map[string]func() logstorage.Logstorage{},
	}
}

func (t *logstorageRegistry) Register(components ...Logstorage) {
	for _, component := range components {
		t.logstorages[createFullName(component.Name)] = component.FactoryMethod
	}
}

func (t *logstorageRegistry) Create(name, version string) (logstorage.Logstorage, error) {
	if method, ok := t.getLogstorage(name, version); ok {
		return method(), nil
	}
	return nil, errors.Errorf("couldn't find logstorage %s/%s", name, version)
}

func (t *logstorageRegistry) getLogstorage(name, version string) (func() logstorage.Logstorage, bool) {
	nameLower := strings.ToLower(name)
	versionLower := strings.ToLower(version)
	logstorageFn, ok := t.logstorages[nameLower+"/"+versionLower]
	if ok {
		return logstorageFn, true
	}
	if components.IsInitialVersion(versionLower) {
		logstorageFn, ok = t.logstorages[nameLower]
	}
	return logstorageFn, ok
}

func createFullName(name string) string {
	return strings.ToLower("logstorage." + name)
}
