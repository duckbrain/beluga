package compose

import (
	"testing"

	"gopkg.in/yaml.v2"
)

func compareLabels(t *testing.T, expected, received Labels) {
	match := len(received) == len(expected)
	for key, value := range received {
		if expected[key] != value {
			match = false
			break
		}
	}
	if !match {
		t.Errorf(`Results do not match; expected %v; received %v`, expected, received)
	}
}

func TestFilesUnmarshal(t *testing.T) {
	for _, set := range []struct {
		Input string
		File  File
	}{
		{
			`
version: '3.0'
x-extra: 12345
services:
	foo:
		image: duckbrain/foo
		labels:
			- beluga-foo=hello-world
		deploy:
			labels:
				beluga-bar: world-hello`,
			File{
				Version: "3.0",
				Services: map[string]Service{
					"foo": Service{
						Labels: Labels{
							"beluga-foo": "hello-world",
						},
						extra: extra{
							fields: map[string]interface{}{
								"image": "hello-world",
							},
						},
					},
				},
			},
		},
	} {
		t.Run(set.Input, func(t *testing.T) {
			file := File{}
			if err := yaml.Unmarshal([]byte(set.Input), &file); err != nil {
				t.Error("parse failed", err)
				return
			}

		})
	}
}

func TestLabels(t *testing.T) {
	for _, set := range []struct {
		Input  string
		Labels Labels
	}{
		{"a:\nb: 3\nc:", Labels{"a": "", "b": "3", "c": ""}},
		{"- a\n- b=3\n- c", Labels{"a": "", "b": "3", "c": ""}},
	} {
		t.Run(set.Input, func(t *testing.T) {
			labels := Labels{}
			if err := yaml.Unmarshal([]byte(set.Input), &labels); err != nil {
				t.Error("parse failed", err)
				return
			}

			compareLabels(t, set.Labels, labels)
		})
	}
}
