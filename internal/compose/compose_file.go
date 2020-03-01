package compose

import (
	"strings"

	"gopkg.in/yaml.v3"
)

type File struct {
	Filename string                 `yaml:"-"`
	Version  string                 `yaml:"version"`
	Services map[string]Service     `yaml:"services"`
	Fields   map[string]interface{} `yaml:"-,inline"`
}

type Service struct {
	Deploy      Deploy                 `yaml:"deploy"`
	Labels      StringMap              `yaml:"labels"`
	Environment StringMap              `yaml:"environment"`
	Fields      map[string]interface{} `yaml:"-,inline"`
}
type Deploy struct {
	Labels StringMap              `yaml:"labels"`
	Fields map[string]interface{} `yaml:"-,inline"`
}

// StringMap represents a map of string to string or array of key/values
// seperated by "=". Values can be left clear represented by nil.
//
// label and environment fields in the service definitions use this type.
type StringMap map[string]*string

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
