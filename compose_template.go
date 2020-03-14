package beluga

import (
	"text/template"
	"strconv"
	"bytes"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"github.com/duckbrain/beluga/internal/compose"
)

type composeTemplateData struct {
	Service struct{
		Port uint16
	}
	Env Environment
}

func composeTemplate(og, tmp string, env Environment) (string, error) {
	if len(tmp) == 0 {
		return og, nil
	}
	templateBase := compose.File{}
	err := yaml.Unmarshal([]byte(tmp), &templateBase)
	if err != nil {
		return "", errors.Wrap(err, "parse template yaml")
	}

	t, err := template.New("").Parse(tmp)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}

	file := compose.File{}
	err = yaml.Unmarshal([]byte(og), &file)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal original")
	}

	err = file.Fields.Merge(templateBase.Fields)
	if err != nil {
		return "", errors.Wrap(err, "merging file base")
	}

	for name, service := range file.Services {
		info := composeTemplateData{}
		port, _ := strconv.ParseUint(service.Labels.Get("us.duckfam.beluga.port"), 10, 16)
		info.Service.Port = uint16(port)

		s := new(bytes.Buffer)
		err := t.Execute(s, info)
		if err != nil {
			return "", errors.Wrap(err, "execute template")
		}

		file.Services[name] = service
	}

	data, err := yaml.Marshal(file)
	return string(data), errors.Wrap(err, "marshal yaml")
}