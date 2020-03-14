package compose

import (
	"strings"

	"gopkg.in/yaml.v3"
)

func Parse(s string) (File, error) {
	f := File{}
	return f, yaml.Unmarshal([]byte(s), &f)
}

func MustParse(s string) File {
	f, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return f
}

type File struct {
	Services map[string]Service `yaml:"services"`
	Fields   Fields             `yaml:"-,inline"`
}

func (f File) TryString() (string, error) {
	b, err := yaml.Marshal(f)
	return string(b), err
}

func (f File) String() string {
	b, _ := yaml.Marshal(f)
	return string(b)
}

type Service struct {
	Deploy      Deploy    `yaml:"deploy,omitempty"`
	Labels      StringMap `yaml:"labels,omitempty"`
	Environment StringMap `yaml:"environment,omitempty"`
	Fields      Fields    `yaml:"-,inline"`
}

type Deploy struct {
	Labels StringMap `yaml:"labels"`
	Fields Fields    `yaml:"-,inline"`
}

type Fields map[string]interface{}

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
