package beluga

import (
	"text/template"
	"strconv"
	"bytes"
	"github.com/pkg/errors"
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
	tFile, err := compose.Parse(tmp)
	if err != nil {
		return "", errors.Wrap(err, "parse template yaml")
	}

	t, err := template.New("").Parse(tmp)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}

	file, err := compose.Parse(og)
	if err != nil {
		return "", errors.Wrap(err, "parse source")
	}

	err = mergo.MergeWithOverwrite(&file.Fields, tFile.Fields)
	if err != nil {
		return "", errors.Wrap(err, "merging file base")
	}

	for name, service := range file.Services {
		info := composeTemplateData{Env: env}
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
		tFile, err := compose.Parse(s.String())
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

	for name, tService := range tFile.Services {
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

	return file.TryString()
}