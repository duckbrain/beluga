package compose

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
	Name string
	extra
}

func (s *Service) MarshalYAML() (interface{}, error) {
	return s.marshalYAML(s)
}
func (s *Service) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return s.unmarshalYAML(s, unmarshal)
}

type Deploy struct {
	Name string
	extra
}

func (s *Deploy) MarshalYAML() (interface{}, error) {
	return s.marshalYAML(s)
}
func (s *Deploy) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return s.unmarshalYAML(s, unmarshal)
}

type Labels map[string]string

//TODO marshal and unmarshal for labels
