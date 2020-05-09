package beluga

import (
	"testing"
)

func TestEnvironmentExpand(t *testing.T) {
	e := Environment{"HELLO": "world", "FOO": "bar"}

	testCases := []struct {
		input string
		output string
	}{
		{
			input: "http://$HELLO/${FOO}",
			output: "http://world/bar",
		},
		{
			input: "http://$HELLO/A$BAZ",
			output: "http://world/A",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			received := e.Expand(GitLabVarMatcher, tC.input)
			if received != tC.output {
				t.Errorf("expected %v, received %v", tC.output, received)
			}
		})
	}
}

func TestEnvironmentExand(t *testing.T) {


}