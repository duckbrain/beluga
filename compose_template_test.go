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
		# implied networks: [ default ]
	foo:
		image: hello-world
		labels:
		- us.duckfam.beluga.port=7080
		networks:
			foo:
	db:
		image: postgres
networks:
	default:
version: "3.0"
`,
			template: `
services:
	{{ if .Service.Labels.Get "us.duckfam.beluga.port" }}
	BELUGA:
		deploy:
			labels:
				"traefik.enable": "true"
				"traefik.http.services.{{ .Env.StackName }}.loadbalancer.server.port": "{{ .Service.Labels.Get "us.duckfam.beluga.port" }}"
		networks:
			- traefik
	{{ end }}
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
		networks:
			traefik: {}
			default: {}
		environment:
			VERSION: "2"
		image: hello-world
	foo:
		deploy:
			labels:
				traefik.enable: "true"
				traefik.http.services.foobar.loadbalancer.server.port: "7080"
		labels:
			us.duckfam.beluga.port: "7080"
		networks:
			foo: {}
			traefik: {}
		image: hello-world
	db:
		image: postgres
networks:
	default: null
	traefik:
		external: true
version: "3.0"
`,
		},
		{
			desc: "existing networks and adding ports",
			original: `
version: '3.0'

services: 
	hello:
		image: hello-world
		environment:
			VERSION: '2'
		labels:
			com.example.deploy.port: 8080
`,
			template: `
services:
	{{ $port := .Service.Labels.Get "com.example.deploy.port" }}
	{{ if $port }}
	BELUGA:
		deploy:
			labels:
				- traefik.enable=true
				- traefik.http.routers.{{ .Env.StackName }}.rule=Host(` + "{{ .Env.Domain }}" + `)
				- traefik.http.services.{{ .Env.StackName }}.loadbalancer.server.port={{ $port }}
				- traefik.http.routers.{{ .Env.StackName }}.entrypoints=web
		networks:
			- traefik
	{{ end }}
networks:
	traefik:
		external: true
`,
			env: Environment{varStackName: "foobar"},
			output: `
version: '3.0'

services: 
	hello:
		deploy:
			labels:
				traefik.enable: "true"
				traefik.http.routers.foobar.entrypoints: web
				traefik.http.routers.foobar.rule: Host()
				traefik.http.services.foobar.loadbalancer.server.port: "8080"
		image: hello-world
		environment:
			VERSION: '2'
		labels:
			com.example.deploy.port: 8080
		networks:
		- traefik
networks:
	traefik:
		external: true
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
