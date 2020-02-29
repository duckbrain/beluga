package compose

import "strings"

// FileFields is exported because of a dependency bug. https://github.com/go-yaml/yaml/issues/463 Do not depend on it.
type FileFields struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

type File struct {
	Filename string
	FileFields
	Extra
}

func (f *File) MarshalYAML() (interface{}, error) {
	panic("I'm hit!!!")
	return f.marshalYAML(&f.FileFields)
}
func (f *File) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return f.unmarshalYAML(&f.FileFields, unmarshal)
}

// ServiceFields is exported because of a dependency bug. https://github.com/go-yaml/yaml/issues/463 Do not depend on it.
type ServiceFields struct {
	Deploy Deploy `yaml:"deploy"`
	Labels Labels `yaml:"labels"`
}
type Service struct {
	Name string
	ServiceFields
	Extra
}

func (s *Service) MarshalYAML() (interface{}, error) {
	return s.marshalYAML(&s.ServiceFields)
}
func (s *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return s.unmarshalYAML(&s.ServiceFields, unmarshal)
}

// DeployFields is exported because of a dependency bug. https://github.com/go-yaml/yaml/issues/463 Do not depend on it.
type DeployFields struct {
	Labels Labels `yaml:"labels"`
}
type Deploy struct {
	DeployFields
	Extra
}

func (d *Deploy) MarshalYAML() (interface{}, error) {
	return d.marshalYAML(&d.DeployFields)
}
func (d *Deploy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return d.unmarshalYAML(&d.DeployFields, unmarshal)
}

type Labels map[string]string

func (l *Labels) UnmarshalYAML(unmarshal func(interface{}) error) error {
	labels := map[string]string(*l)
	if labels == nil {
		labels = make(map[string]string)
	}
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
	*l = Labels(labels)
	return nil
}
