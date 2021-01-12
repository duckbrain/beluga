package beluga

import (
	"bytes"
	"text/template"

	"github.com/duckbrain/beluga/internal/compose"
	"github.com/pkg/errors"
)

type composeTemplateData struct {
	Src         *compose.File
	Service     *compose.Service
	ServiceName string
	Volume      *compose.VolumeDefinition
	VolumeName  string
	Env         Environment
}

const tmplKey = "BELUGA"

func composeTemplate(src, tmpl string, env Environment) (string, error) {
	if len(tmpl) == 0 {
		return src, nil
	}

	t, err := template.New("").Parse(tmpl)
	if err != nil {
		return "", errors.Wrap(err, "parse template")
	}

	file, err := compose.Parse(src)
	if err != nil {
		return "", errors.Wrap(err, "yaml parse source")
	}

	parseTemplate := func(data composeTemplateData) (*compose.File, error) {
		s := new(bytes.Buffer)
		data.Env = env
		data.Src = file
		if data.Service == nil {
			data.Service = &compose.Service{}
		}
		if data.Volume == nil {
			data.Volume = &compose.VolumeDefinition{}
		}
		err := t.Execute(s, data)
		if err != nil {
			return nil, errors.Wrap(err, "execute template")
		}
		f, err := compose.Parse(s.String())
		return f, errors.Wrap(err, "yaml parse template output")
	}

	tFile, err := parseTemplate(composeTemplateData{})
	if err != nil {
		return "", err
	}

	err = file.Merge(tFile)
	if err != nil {
		return "", errors.Wrap(err, "merging file base")
	}

	for name, service := range file.Services {
		tFile, err := parseTemplate(composeTemplateData{Service: service, ServiceName: name})
		if err != nil {
			return "", errors.Wrapf(err, "service \"%v\"", name)
		}

		if tService, ok := tFile.Services[tmplKey]; ok {
			err = service.Merge(tService)
			if err != nil {
				return "", errors.Wrapf(err, "merge service %v", name)
			}
		}

		file.Services[name] = service
	}

	for name, tService := range tFile.Services {
		if name == tmplKey {
			continue
		}
		service, ok := file.Services[name]
		if ok {
			err = service.Merge(tService)
			if err != nil {
				return "", errors.Wrapf(err, "merge template service %v", name)
			}
			file.Services[name] = service
		} else {
			file.Services[name] = tService
		}
	}

	for name, volume := range file.Volumes {
		tFile, err := parseTemplate(composeTemplateData{Volume: &volume, VolumeName: name})
		if err != nil {
			return "", errors.Wrapf(err, "volume \"%v\"", name)
		}

		if tVolume, ok := tFile.Volumes[tmplKey]; ok {
			err = volume.Merge(&tVolume)
			if err != nil {
				return "", errors.Wrapf(err, "merge volume %v", name)
			}
		}

		file.Volumes[name] = volume
	}

	for name, tVolume := range tFile.Volumes {
		if name == tmplKey {
			continue
		}
		volume, ok := file.Volumes[name]
		if ok {
			err = volume.Merge(&tVolume)
			if err != nil {
				return "", errors.Wrapf(err, "merge template volume %v", name)
			}
			file.Volumes[name] = volume
		} else {
			file.Volumes[name] = tVolume
		}
	}

	return file.TryString()
}
