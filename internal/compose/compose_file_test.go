package compose

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/sergi/go-diff/diffmatchpatch"
	"gopkg.in/yaml.v3"
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

var helloWorld = "hello-world"
var emptyString = ""

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
				Fields: map[string]interface{}{
					"version": "3.0",
					"x-extra": 12345,
				},
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
				beluga-bar: hello-world
`,
			File{
				Services: map[string]Service{
					"foo": {
						Labels: StringMap{
							"beluga-foo": &helloWorld,
						},
						Deploy: Deploy{
							Labels: StringMap{
								"beluga-bar": &helloWorld,
							},
						},
						Networks: Networks{
							"default": {},
						},
						Fields: map[string]interface{}{
							"image": "duckbrain/foo",
						},
					},
				},
				Fields: map[string]interface{}{
					"version": "3.0",
					"x-extra": 12345,
				},
			},
		},
	} {
		t.Run(set.Name, func(t *testing.T) {
			file := File{}
			// yaml doesn't allow tabs. We'll replace them with spaces to avoid
			// mixed indentation that may confuse editors.
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
		Name   string
		Input  string
		Labels StringMap
	}{
		{"dictionary format", "a:\nb: hello-world\nc: ''", StringMap{"a": nil, "b": &helloWorld, "c": &emptyString}},
		{"array format", "- a\n- b=hello-world\n- c=", StringMap{"a": nil, "b": &helloWorld, "c": &emptyString}},
	} {
		t.Run(set.Name, func(t *testing.T) {
			var labels = StringMap{}
			if err := yaml.Unmarshal([]byte(set.Input), &labels); err != nil {
				t.Error("parse failed", err)
				return
			}
			assertEqual(t, set.Labels, labels)
			assertYamlEqual(t, set.Labels, labels)
		})
	}
}
