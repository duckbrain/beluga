package compose

import "gopkg.in/yaml.v2"

// Extra is exported because of a dependency bug. https://github.com/go-yaml/yaml/issues/463 Do not depend on it.
type Extra struct {
	// additionally stores the fields in the yaml so when marshalling them back,
	// they will be preserved
	fields map[string]interface{} `yaml:"-"`
}

func (e *Extra) marshalYAML(s interface{}) (interface{}, error) {
	b, err := yaml.Marshal(s)
	if err != nil {
		return nil, err
	}
	x := map[string]interface{}{}
	err = yaml.Unmarshal(b, &x)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	for key, value := range e.fields {
		m[key] = value
	}
	for key, value := range x {
		m[key] = value
	}
	return m, nil
}

func (e *Extra) unmarshalYAML(s interface{}, unmarshal func(interface{}) error) error {
	err := unmarshal(s)
	if err != nil {
		return err
	}
	return unmarshal(&e.fields)
}
