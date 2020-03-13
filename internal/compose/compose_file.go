package compose

import (
	"strings"

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

func (s Service) Merge(b Service) Service {
	r := Service{}
	r.Deploy.Fields = s.Deploy.Fields.Merge(b.Deploy.Fields)
	r.Deploy.Labels = s.Deploy.Labels.Merge(b.Deploy.Labels)
	r.Fields = s.Fields.Merge(b.Fields)
	r.Environment = s.Environment.Merge(b.Environment)
	r.Labels = s.Labels.Merge(b.Labels)
	return r
}

type Deploy struct {
	Labels StringMap `yaml:"labels"`
	Fields Fields    `yaml:"-,inline"`
}

type Fields map[string]interface{}

func (f Fields) Merge(b Fields) Fields {
	if f == nil && b == nil {
		return nil
	}
	r := make(map[string]interface{})
	for key, value := range f {
		r[key] = value
	}
	for key, value := range b {
		r[key] = value
	}
	return r
}

// StringMap represents a map of string to string or array of key/values
// seperated by "=". Values can be left clear represented by nil.
//
// label and environment fields in the service definitions use this type.
type StringMap map[string]*string

func (m StringMap) Merge(b StringMap) StringMap {
	if m == nil && b == nil {
		return nil
	}
	r := make(map[string]*string)
	for key, value := range m {
		r[key] = value
	}
	for key, value := range b {
		r[key] = value
	}
	return r
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
