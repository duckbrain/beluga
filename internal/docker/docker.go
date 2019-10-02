package docker

import (
	"os"
	"os/exec"
)

func run(s string, args ...string) error {
	c := exec.Command(s, args...)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func Build(context, dockerfile, tag string) error {
	return run(
		"docker", "build",
		context,
		"-f", dockerfile,
		"-t", tag,
	)
}

func Tag(src, dst string) error {
	return run("docker", "tag", src, dst)
}

func Push(tag string) error {
	return run("docker", "push", tag)
}

func Login(hostname, username, password string) error {
	return run(
		"docker", "login",
		hostname,
		"-u", username,
		"-p", password,
	)
}
