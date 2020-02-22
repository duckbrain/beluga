package beluga

import (
	"errors"

	"gopkg.in/yaml.v2"
)

type ComposeFileProcessor struct {
}

func yamlField(s yaml.MapSlice, name string) (res interface{}, ok bool) {
	for _, pair := range s {
		if pair.Key == name {
			return pair.Value, true
		}
	}
	return nil, false
}
func yamlSlice(s yaml.MapSlice, name string) yaml.MapSlice {
	res, ok := yamlField(s, name)
	if !ok {
		panic("field not found")
	}
	return res.(yaml.MapSlice)
}

func (e *ComposeFileProcessor) Process(input []byte) (err error) {
	data := yaml.MapSlice{}
	err := yaml.Unmarshal(input, &data)
	if err != nil {
		return err
	}

	defer func() {
		r := recover()
		if r != nil {
			err = errors.New("failed")
			return
		}
	}()

	services := yamlSlice(data, "services")

	for _, servicePair := range services {
		service := servicePair.Value.(yaml.MapSlice)
		if labels, ok := yamlField(service, "labels"); ok {
			belugaHost := yamlField(labels.(yaml.MapSlice), "us.duckfam.beluga.host"
		}
	}

	return nil
}
