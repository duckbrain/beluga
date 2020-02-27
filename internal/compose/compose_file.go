package compose

import "strings"

type File struct {
	Filename string             `yaml:"-"`
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
	extra
}

func (f *File) MarshalYAML() (interface{}, error) {
	return f.marshalYAML(f)
}
func (f *File) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return f.unmarshalYAML(f, unmarshal)
}

type Service struct {
	Name   string
	Deploy Deploy `yaml:"deploy"`
	Labels Labels `yaml:"labels"`
	extra
}

func (s *Service) MarshalYAML() (interface{}, error) {
	return s.marshalYAML(s)
}
func (s *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return s.unmarshalYAML(s, unmarshal)
}

type Deploy struct {
	Labels Labels `yaml:"labels"`
	extra
}

func (s *Deploy) MarshalYAML() (interface{}, error) {
	return s.marshalYAML(s)
}
func (s *Deploy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return s.unmarshalYAML(s, unmarshal)
}

type Labels map[string]string

func (l Labels) UnmarshalYAML(unmarshal func(interface{}) error) error {
	labels := map[string]string(l)
	if err := unmarshal(labels); err == nil {
		return nil
	}

	var lines []string
	if err := unmarshal(&lines); err != nil {
		return err
	}
	for _, line := range lines {
		i := strings.Index(line, "=")
		if i == -1 {
			labels[line] = ""
			continue
		}
		labels[line[:i]] = line[i+1:]
	}
	return nil
}
