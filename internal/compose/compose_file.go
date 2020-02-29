package compose

import "strings"

type FileFields struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

type File struct {
	Filename string
	FileFields
	extra
}

func (f *File) MarshalYAML() (interface{}, error) {
	panic("I'm hit!!!")
	return f.marshalYAML(&f.FileFields)
}
func (f *File) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return f.unmarshalYAML(&f.FileFields, unmarshal)
}

type serviceFields struct {
	Deploy Deploy `yaml:"deploy"`
	Labels Labels `yaml:"labels"`
}
type Service struct {
	Name string
	serviceFields
	extra
}

func (s *Service) MarshalYAML() (interface{}, error) {
	return s.marshalYAML(&s.serviceFields)
}
func (s *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := s.unmarshalYAML(&s.serviceFields, unmarshal)
	if err != nil {
		return err
	}
	if s.Labels == nil {
		s.Labels = Labels{}
	}
	return nil
}

type deployFields struct {
	Labels Labels `yaml:"labels"`
}
type Deploy struct {
	deployFields
	extra
}

func (d *Deploy) MarshalYAML() (interface{}, error) {
	return d.marshalYAML(&d.deployFields)
}
func (d *Deploy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	err := d.unmarshalYAML(&d.deployFields, unmarshal)
	if err != nil {
		return err
	}
	if d.Labels == nil {
		d.Labels = Labels{}
	}
	return nil
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
