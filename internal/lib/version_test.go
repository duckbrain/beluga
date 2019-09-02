package lib

import "testing"

func TestParseVersion(t *testing.T) {
	tests := struct{
		Input: "v1.2.3",
		Output: "1.2.3",
		Flavor: "",
		Prerelease: "" 
	}
}
