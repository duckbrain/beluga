package compose

import "gopkg.in/yaml.v2"

type extra struct {
	// additionally stores the fields in the yaml so when marshalling them back,
	// they will be preserved
	fields map[string]interface{} `yaml:"-"`
}

func (e *extra) marshalYAML(s interface{}) (interface{}, error) {
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

func (e *extra) unmarshalYAML(s interface{}, unmarshal func(interface{}) error) error {
	err := unmarshal(s)
	if err != nil {
		return err
	}
	return unmarshal(&e.fields)
}
