package compose

import (
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type File struct {
	Services map[string]Service `yaml:"services"`
	Fields   Fields             `yaml:"-,inline"`
}

type Service struct {
	Deploy      Deploy    `yaml:"deploy"`
	Labels      StringMap `yaml:"labels"`
	Environment StringMap `yaml:"environment"`
	Fields      Fields    `yaml:"-,inline"`
}

func (s *Service) Merge(b Service) {
	s.Deploy.Labels.Merge(b.Deploy.Labels)
	s.Environment.Merge(b.Environment)
	s.Fields.Merge(b.Fields)
	s.Labels.Merge(b.Labels)
}

type Deploy struct {
	Labels StringMap `yaml:"labels"`
	Fields Fields    `yaml:"-,inline"`
}

type Fields map[string]interface{}

func (l Fields) Merge(b Fields) {
	for k, v := range b {
		merge(&l[k], v)
	}
}

// StringMap represents a map of string to string or array of key/values
// seperated by "=". Values can be left clear represented by nil.
//
// label and environment fields in the service definitions use this type.
type StringMap map[string]*string

func (l StringMap) Merge(b StringMap) {
	for k, v := range b {
		l[k] = v
	}
}

func merge(a *interface{}, b interface{}) error {
	switch x := b.(type) {
	case string, int64, float64:
		*a = x
	case map[interface{}]interface{}:
		y, ok := a.(map[interface{}]interface{})
		if !ok {
			return errors.New("incompatible types")
		}
		for k, v := range x {
			y[k] = v
		}
	default:
		return errors.Errorf("unknown type %T", x)
	}
	return nil
}

func (l StringMap) Get(key string) string {
	s := l[key]
	if s == nil {
		return ""
	}
	return *s
}
func (l StringMap) Unset(key string) {
	delete(l, key)
}
func (l StringMap) Clear(key string) {
	l[key] = nil
}
func (l StringMap) Set(key, value string) {
	l[key] = &value
}

func (l *StringMap) UnmarshalYAML(value *yaml.Node) error {
	values := map[string]*string(*l)
	if values == nil {
		values = make(map[string]*string)
	}
	if err := value.Decode(values); err == nil {
		*l = StringMap(values)
		return nil
	}

	var lines []string
	if err := value.Decode(&lines); err != nil {
		return err
	}
	for _, line := range lines {
		i := strings.Index(line, "=")
		if i == -1 {
			values[line] = nil
			continue
		}
		s := line[i+1:]
		values[line[:i]] = &s
	}
	*l = StringMap(values)
	return nil
}
