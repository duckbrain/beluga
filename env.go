package beluga

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/gobuffalo/envy"
)

func init() {
	// Load all environment variables
	_ = envy.Load()
}

const (
	envProduction = "production"
	envStaging    = "staging"
	envReview     = "review"
)

var envs = []func(Environment){}

func Env(logger logrus.StdLogger) Environment {
	// Start with the environment variables
	e := Environment(envy.Map())

	// CI environments
	for _, read := range envs {
		read(e)
	}

	envReadEnvOverrides(e)

	// Parse the image tempate
	err := envParseImageTemplate(e)
	if err != nil {
		logger.Println(err)
	}

	// Try to fill in blanks with git
	gitEnvRead(e)

	return e
}

type Environment map[string]string

// Target branch for PRs/MRs. Defaults to master.
func (e Environment) DefaultBranch() string {
	return e[varDefaultBranch]
}

func (e Environment) Validate() error {
	errs := Errors{}

	// Must be set
	for _, key := range []string{
		varStackName,
	} {
		if v := e[key]; v == "" {
			errs.Append(errors.Errorf("%v must be non-empty", key))
		}
	}

	return errs.Err()
}

func (e Environment) Get(key, fallback string) string {
	v := e[key]
	if v == "" {
		return fallback
	}
	return v
}

func (e Environment) Keys() []string {
	keys := []string{}
	for key := range e {
		keys = append(keys, key)
	}
	return []string(keys)
}

func (e Environment) CommonKeys() []string {
	keys := []string{}
	for _, key := range knownVarNames {
		if v, ok := e[key]; ok && v != "" {
			keys = append(keys, key)
		}
	}
	return []string(keys)
}

func (e Environment) Merge(src Environment) {
	for key, value := range src {
		if value != "" {
			e[key] = value
		}
	}
}

func (e Environment) MergeMissing(src Environment) {
	for key, value := range src {
		if e[key] == "" && value != "" {
			e[key] = value
		}
	}
}

func (e Environment) clone() Environment {
	n := make(Environment)
	for key, value := range e {
		if value != "" {
			n[key] = value
		}
	}
	return n
}

// envReadEnvOverrides overrides values for a specific environment by reading
// the value of the BELUGA_* variable as an env file.
func envReadEnvOverrides(e Environment) {
	v := e["BELUGA_"+strings.ToUpper(e[varEnvironment])]
	s, err := godotenv.Unmarshal(v)
	if err != nil {
		log.Println("beluga: env parse: ", err)
	}

	e.Merge(Environment(s))
}

func envParseImageTemplate(e Environment) error {
	t := e.ImagesTemplate()
	if t == "" {
		return nil
	}
	imgTmpl, err := template.New("").Parse(t)
	if err != nil {
		return err
	}
	data := struct{ Env Environment }{e}
	buf := new(bytes.Buffer)
	err = imgTmpl.Execute(buf, data)
	if err != nil {
		return err
	}
	images := buf.String()
	e.MergeMissing(Environment{varImages: images})

	imgs := strings.Fields(e.Images())
	if len(imgs) > 0 {
		e.MergeMissing(Environment{varImage: imgs[0]})
		e[varImage] = imgs[0]
	}
	return nil
}

type EnvFormat func(key, val string) (string, error)

var BashFormat EnvFormat = func(key, val string) (string, error) {
	// TODO escape special characters like newline and single quote
	return fmt.Sprintf("export %v='%v'", key, val), nil
}
var EnvFileFormat EnvFormat = func(key, val string) (string, error) {
	// ENV file format doesn't support newlines.
	if strings.Contains(val, "\n") {
		return "", errors.New("cannot represent newlines in env file format")
	}
	return fmt.Sprintf("%v=%v", key, val), nil
}
var GoEnvFormat = func(key, val string) (string, error) {
	// Newlines are okay because go format keeps strings in slice
	return fmt.Sprintf("%v=%v", key, val), nil
}

// Format returns the env variables in the specified format. If allKeys is true,
// it will return all values found in the environment, otherwise it only outputs
// keys known to Beluga. Errors may occur for individual key outputs, in these
// cases the function will return a valid slice of strings, and the error.
func (env Environment) Format(format EnvFormat, keys []string) ([]string, error) {
	var err error
	var values = make([]string, len(keys))
	for i, key := range keys {
		values[i], err = format(key, env[key])
	}
	return values, err
}
