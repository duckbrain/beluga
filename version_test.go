package beluga

import "testing"

func TestParseVersion(t *testing.T) {
	type test struct {
		Input       string
		Output      string
		Application string
		Environment string
		IsInvalid   bool
	}
	testCases := []test{
		{
			Input:  "v1.2.3",
			Output: "1.2.3",
		},
		{
			Input:       "app/1.2.3",
			Output:      "1.2.3",
			Application: "app",
		},
		{
			Input:       "app/v1.2.3-beta.0",
			Output:      "1.2.3-beta.0",
			Application: "app",
			Environment: "beta",
		},
		{
			Input:       "1.2.3-beta",
			Output:      "1.2.3-beta",
			Environment: "beta",
		},
	}
	for _, c := range testCases {
		t.Run(c.Input, func(t *testing.T) {
			v := parseVersion(c.Input)
			if c.IsInvalid != (v == nil) {
				if c.IsInvalid {
					t.Errorf("got version \"%v\" expected not to parse", v.Version())
				} else {
					t.Errorf("did not parse expected to parse")
				}
				return
			}
			if v.Version() != c.Output {
				t.Errorf("got version \"%v\" expected \"%v\"", v.Version(), c.Output)
			}
			if v.Application() != c.Application {
				t.Errorf("got application \"%v\" expected \"%v\"", v.Application(), c.Application)
			}
			if v.Environment() != c.Environment {
				t.Errorf("got environment \"%v\" expected \"%v\"", v.Environment(), c.Environment)
			}
		})
	}
}
