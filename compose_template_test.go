package beluga

import (
	"strings"
	"testing"

	"github.com/duckbrain/beluga/internal/compose"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func t2s(s string) string {
	return strings.ReplaceAll(strings.TrimSpace(s), "\t", "    ")
}

func TestComposeTemplate(t *testing.T) {
	testCases := []struct {
		desc     string
		original string
		template string
		env      Environment
		output   string
	}{
		{
			desc: "existing networks and adding ports",
			original: `
services:
	hello:
		environment:
			VERSION: "2"
		image: hello-world
		labels:
		- us.duckfam.beluga.port=8080
networks:
	default:
version: "3.0"
`,
			template: `
services:
	BELUGA:
		deploy:
			labels:
				"traefik.enable": "true"
				"traefik.http.services.{{ .Env.StackName }}.loadbalancer.server.port": "{{ .Service.Port }}"
	backup:
		image: example/backups
networks:
	traefik:
		external: true
`,
			env: Environment{varStackName: "foobar"},
			output: `
services:
	backup:
		image: example/backups
	hello:
		deploy:
			labels:
				traefik.enable: "true"
				traefik.http.services.foobar.loadbalancer.server.port: "8080"
		labels:
			us.duckfam.beluga.port: "8080"
		environment:
			VERSION: "2"
		image: hello-world
networks:
	default: null
	traefik:
		external: true
version: "3.0"
`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			res, err := composeTemplate(t2s(tC.original), t2s(tC.template), tC.env)
			if err != nil {
				t.Fatal(err)
			}

			expect := compose.MustParse(t2s(tC.output)).String()

			if res == expect {
				return
			}

			dmp := diffmatchpatch.New()
			diff := dmp.DiffMain(expect, res, false)
			if len(diff) > 0 {
				t.Log(diff)
				t.Error("yaml doesn't match\n", dmp.DiffPrettyText(diff))
			}
		})
	}
}
