package portainer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalNilEnv(t *testing.T) {
	var env Env
	const input = `[{"name":"MY_KEY","value":"my_val"}]`
	err := json.Unmarshal([]byte(input), &env)
	require.NoError(t, err)
	assert.Equal(t, env["MY_KEY"], "my_val")
}
