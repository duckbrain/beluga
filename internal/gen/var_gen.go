//go:generate go run ./var_gen.go

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gobuffalo/flect"
)

var names = []string{
	"Application",
	"Domain",
	"Environment",
	"Image",
	"Registry",
	"RegistryPassword",
	"RegistryUsername",
	"Variant",
	"Version",
}

var t = template.Must(template.New("").Funcs(template.FuncMap{
	"VarName": VarName,
	"EnvName": EnvName,
}).Parse(`
package lib

const (
	{{- range .}}
	{{ . | VarName }} = "{{ . | EnvName }}"
	{{- end }}
)
`))

func main() {
	const filename = "../lib/var.go"
	output, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(output, names)
	if err != nil {
		fmt.Println(err)
	}
	err = output.Close()
	if err != nil {
		fmt.Println(err)
	}
	err = exec.Command("go", "fmt", filename).Run()
	if err != nil {
		fmt.Println(err)
	}
}

func VarName(s string) string {
	return "var" + flect.Pascalize(s)
}

func EnvName(s string) string {
	return "BELUGA_" + strings.ToUpper(flect.Underscore(s))
}
