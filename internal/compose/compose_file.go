package compose

import (
	"strings"

	"github.com/imdario/mergo"
	"gopkg.in/yaml.v3"
)

func Parse(s string) (*File, error) {
	f := &File{}
	if err := yaml.Unmarshal([]byte(s), f); err != nil {
		return f, err
	}
	for _, s := range f.Services {
		s.file = f
	}
	return f, nil
}

func MustParse(s string) *File {
	f, err := Parse(s)
	if err != nil {
		panic(err)
	}
	return f
}

type File struct {
	Services map[string]*Service `yaml:"services"`
	Fields   Fields              `yaml:"-,inline"`
}

func (f *File) Merge(a *File) error {
	return mergo.MergeWithOverwrite(&f.Fields, a.Fields)
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
	file        *File
	Deploy      Deploy    `yaml:"deploy,omitempty"`
	Labels      StringMap `yaml:"labels,omitempty"`
	Environment StringMap `yaml:"environment,omitempty"`
	Networks    Networks  `yaml:"networks,omitempty"`
	Fields      Fields    `yaml:"-,inline"`
}

func (s *Service) UnmarshalYAML(value *yaml.Node) error {
	type S Service
	v := S(*s)
	err := value.Decode(&v)
	if err != nil {
		return err
	}
	*s = Service(v)
	if len(s.Networks) == 0 {
		s.Networks = Networks{"default": {}}
	}
	return nil
}

func (s *Service) Merge(a *Service) error {
	networks := s.Networks.Clone()
	err := mergo.MergeWithOverwrite(s, *a)
	if err != nil {
		return err
	}

	if networks.IsZero() && len(s.file.Services) == 1 {
		s.Networks = a.Networks.Clone()
	} else {
		for name, net := range a.Networks {
			networks[name] = net
		}
		s.Networks = networks
	}
	return nil
}

type Deploy struct {
	Labels StringMap `yaml:"labels,omitempty"`
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

type Network struct {
	Aliases     []string `yaml:"aliases,omitempty"`
	IPv4Address string   `yaml:"ipv4_address,omitempty"`
	IPv6Address string   `yaml:"ipv6_address,omitempty"`
	Fields      Fields   `yaml:"-,inline"`
}

func (n Network) IsZero() bool {
	return len(n.Aliases) == 0 && n.IPv4Address == "" && n.IPv6Address == "" && len(n.Fields) == 0
}

type Networks map[string]Network

func (n Networks) Clone() Networks {
	c := Networks{}
	for name, net := range n {
		c[name] = net
	}
	return c
}

func (n Networks) Has(s string) bool {
	_, ok := n[s]
	return ok
}

func (n Networks) IsZero() bool {
	d, ok := n["default"]
	return len(n) == 1 && ok && d.IsZero()
}

func (n *Networks) UnmarshalYAML(value *yaml.Node) (err error) {
	values := map[string]Network(*n)
	if values == nil {
		values = map[string]Network{}
	}

	falsyValues := map[string]*Network{}
	err = value.Decode(&falsyValues)
	if err == nil {
		for k, v := range falsyValues {
			if v == nil {
				v = &Network{}
			}
			values[k] = *v
		}
		*n = Networks(values)
		return nil
	}

	var lines []string
	err = value.Decode(&lines)
	if err != nil {
		return err
	}
	for _, name := range lines {
		values[name] = Network{}
	}
	*n = Networks(values)
	return nil
}
