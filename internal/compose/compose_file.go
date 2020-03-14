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
	Deploy      Deploy    `yaml:"deploy,omitempty"`
	Labels      StringMap `yaml:"labels,omitempty"`
	Environment StringMap `yaml:"environment,omitempty"`
	Fields      Fields    `yaml:"-,inline"`
}

func (s *Service) Merge(b Service) {
	s.Deploy.Fields.Merge(b.Deploy.Fields)
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

func (l *Fields) Merge(b Fields) error {
	if l == nil {
		*l = make(Fields)
	}
	for k, v := range b {
		x, err := merge((*l)[k], v)
		if err != nil {
			return err
		}
		(*l)[k] = x
	}
	return nil
}

// StringMap represents a map of string to string or array of key/values
// seperated by "=". Values can be left clear represented by nil.
//
// label and environment fields in the service definitions use this type.
type StringMap map[string]*string

func (l *StringMap) Merge(b StringMap) {
	if *l == nil {
		*l = make(StringMap)
	}
	for k, v := range b {
		(*l)[k] = v
	}
}

func merge(a interface{}, b interface{}) (interface{}, error) {
	switch x := b.(type) {
	case string, int, int64:
		return x, nil
	case map[interface{}]interface{}:
		y, ok := a.(map[interface{}]interface{})
		if !ok {
			return nil, errors.Errorf("incompatible types %T and %T", a, b)
		}
		for k, v := range x {
			z, err := merge(y[k], v)
			if err != nil {
				return nil, err
			}
			y[k] = z
		}
		return y, nil
	case map[string]interface{}:
		y, ok := a.(map[string]interface{})
		if !ok {
			return nil, errors.Errorf("incompatible types %T and %T", a, b)
		}
		for k, v := range x {
			z, err := merge(y[k], v)
			if err != nil {
				return nil, err
			}
			y[k] = z
		}
		return y, nil
	default:
		return nil, errors.Errorf("unknown type %T", x)
	}
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
