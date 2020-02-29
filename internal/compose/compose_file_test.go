package compose

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v2"
)

func assertEqual(t *testing.T, expected, received interface{}) {
	if diff := deep.Equal(expected, received); diff != nil {
		t.Errorf(`Results do not match: %v`, diff)
	}
}
func assertYamlEqual(t *testing.T, expected, received interface{}) {
	expectedYaml, err := yaml.Marshal(expected)
	if err != nil {
		t.Error("marshal expected failed", err)
	}
	receivedYaml, err := yaml.Marshal(received)
	if err != nil {
		t.Error("marshal expected failed", err)
	}

	if string(expectedYaml) == string(receivedYaml) {
		return
	}

	dmp := diffmatchpatch.New()
	diff := dmp.DiffMain(string(expectedYaml), string(receivedYaml), false)
	if len(diff) > 0 {
		t.Error("marshalled yaml doesn't match", dmp.DiffPrettyText(diff))
	}
}

func TestFilesUnmarshal(t *testing.T) {
	for _, set := range []struct {
		Name  string
		Input string
		File  File
	}{
		{
			Name: "Simple file with no services",
			Input: `
version: '3.0'
x-extra: 12345`,
			File: File{
				FileFields: FileFields{
					Version: "3.0",
				},
				Extra: Extra{fields: map[string]interface{}{
					"x-extra": 12345,
				}},
			},
		},
		{
			"Full example file",
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
				FileFields: FileFields{
					Version: "3.0",
					Services: map[string]Service{
						"foo": Service{
							ServiceFields: ServiceFields{
								Labels: Labels{
									"beluga-foo": "hello-world",
								},
							},
							Extra: Extra{
								fields: map[string]interface{}{
									"image": "hello-world",
								},
							},
						},
					},
				},
			},
		},
	} {
		t.Run(set.Name, func(t *testing.T) {
			file := File{}
			y := strings.ReplaceAll(set.Input, "\t", "  ")
			if err := yaml.Unmarshal([]byte(y), &file); err != nil {
				t.Error("parse failed", err)
				return
			}
			assertEqual(t, set.File, file)
			assertYamlEqual(t, set.File, file)
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
			var labels = Labels{}
			if err := yaml.Unmarshal([]byte(set.Input), &labels); err != nil {
				t.Error("parse failed", err)
				return
			}
			assertEqual(t, set.Labels, labels)
			assertYamlEqual(t, set.Labels, labels)
		})
	}
}
