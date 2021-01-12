package beluga

import (
	"fmt"
	"testing"
)

func TestEnvironmentExpand(t *testing.T) {
	e := Environment{"HELLO": "world", "FOO": "bar"}

	testCases := []struct {
		input  string
		output string
	}{
		{
			input:  "http://$HELLO/${FOO}",
			output: "http://world/bar",
		},
		{
			input:  "http://$HELLO/A$BAZ",
			output: "http://world/A",
		},
		{
			input:  "http$$://$HELLO}/${FOO",
			output: "http$://world}/${FOO",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			received := e.expand(gitLabVarMatcher, tC.input)
			if received != tC.output {
				t.Errorf("expected %v, received %v", tC.output, received)
			}
		})
	}
}

func TestEnvironmentOverrides(t *testing.T) {
	e := Environment{
		varOverrides: `
production:
  SECRET: 'Hahaha'
staging:
  SECRET: 'hello ${PLACE}'
`,
		"SECRET": "bar",
		"PLACE":  "world!",
	}

	testCases := []struct {
		environment string
		secret      string
	}{
		{
			secret: "bar",
		},
		{
			environment: "production",
			secret:      "Hahaha",
		},
		{
			environment: "staging",
			secret:      "hello world!",
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("environment: \"%v\"", tC.environment), func(t *testing.T) {
			e[varEnvironment] = tC.environment
			err := envReadEnvOverrides(e)
			if err != nil {
				t.Fatal(err)
			}
			secret := e["SECRET"]
			if secret != tC.secret {
				t.Errorf("expected \"%v\", received \"%v\"", tC.secret, secret)
			}
		})
	}
}
