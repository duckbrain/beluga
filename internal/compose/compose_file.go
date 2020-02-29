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
	Deploy Deploy                 `yaml:"deploy"`
	Labels Labels                 `yaml:"labels"`
	Fields map[string]interface{} `yaml:"-,inline"`
}
type Deploy struct {
	Labels Labels                 `yaml:"labels"`
	Fields map[string]interface{} `yaml:"-,inline"`
}

type Labels map[string]*string

func (l *Labels) UnmarshalYAML(value *yaml.Node) error {
	labels := map[string]*string(*l)
	if labels == nil {
		labels = make(map[string]*string)
	}
	if err := value.Decode(labels); err == nil {
		*l = Labels(labels)
		return nil
	}

	var lines []string
	if err := value.Decode(&lines); err != nil {
		return err
	}
	for _, line := range lines {
		i := strings.Index(line, "=")
		if i == -1 {
			labels[line] = nil
			continue
		}
		s := line[i+1:]
		labels[line[:i]] = &s
	}
	*l = Labels(labels)
	return nil
}
