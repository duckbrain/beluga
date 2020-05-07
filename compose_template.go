package beluga

import (
	"bytes"
	"text/template"

	"github.com/duckbrain/beluga/internal/compose"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

func composeTemplate(og, tmp string, env Environment) (string, error) {
	if len(tmp) == 0 {
		return og, nil
	}

	t, err := template.New("").Parse(tmp)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}

	file, err := compose.Parse(og)
	if err != nil {
		return "", errors.Wrap(err, "yaml parse source")
	}

	parseTemplate := func(service compose.Service) (compose.File, error) {
		s := new(bytes.Buffer)
		err := t.Execute(s, struct {
			Src     compose.File
			Service compose.Service
			Env     Environment
		}{file, service, env})
		if err != nil {
			return compose.File{}, errors.Wrap(err, "execute template")
		}
		f, err := compose.Parse(s.String())
		return f, errors.Wrap(err, "yaml parse template output")
	}

	tFile, err := parseTemplate(compose.Service{})
	if err != nil {
		return "", err
	}

	err = mergo.MergeWithOverwrite(&file.Fields, tFile.Fields)
	if err != nil {
		return "", errors.Wrap(err, "merging file base")
	}

	for name, service := range file.Services {
		tFile, err := parseTemplate(service)
		if err != nil {
			return "", errors.Wrapf(err, "service \"%v\"", name)
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
