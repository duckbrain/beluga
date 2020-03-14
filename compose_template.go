package beluga

import (
	"text/template"
	"strconv"
	"bytes"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"github.com/duckbrain/beluga/internal/compose"
	 "github.com/imdario/mergo"
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

	err = mergo.MergeWithOverwrite(&file.Fields, templateBase.Fields)
	if err != nil {
		return "", errors.Wrap(err, "merging file base")
	}

	for name, service := range file.Services {
		info := composeTemplateData{}
		port, _ := strconv.ParseUint(service.Labels.Get("us.duckfam.beluga.port"), 10, 16)
		info.Service.Port = uint16(port)

		if info.Service.Port == 0 {
			continue
		}

		s := new(bytes.Buffer)
		err := t.Execute(s, info)
		if err != nil {
			return "", errors.Wrap(err, "execute template")
		}
		tFile := compose.File{}
		err = yaml.Unmarshal(s.Bytes(), &tFile)
		if err != nil {
			return "", errors.Wrap(err, "parse templated yaml")
		}

		if tService, ok := tFile.Services["BELUGA"]; ok {
			err = mergo.MergeWithOverwrite(&service, tService)
			if err != nil {
				return "", errors.Wrapf(err, "merge service %v", name)
			}
		}

		file.Services[name] = service
	}

	for name, tService := range templateBase.Services {
		if name == "BELUGA" {
			continue
		}
		service, ok := file.Services[name]
		if ok {
			err = mergo.MergeWithOverwrite(&service, tService)
			if err != nil {
				return "", errors.Wrapf(err, "merge template service %v", name)
			}
			file.Services[name] = service
		} else {
			file.Services[name] = tService
		}
	}

	data, err := yaml.Marshal(file)
	return string(data), errors.Wrap(err, "marshal yaml")
}