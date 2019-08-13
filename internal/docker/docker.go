package docker

import "os/exec"

func Build(context, dockerfile, tag string) error {
	return exec.Command(
		"docker", "build",
		context,
		"-f", dockerfile,
		"-t", tag,
	).Run()
}

func Tag(src, dst string) error {
	return exec.Command("docker", "tag", src, dst).Run()
}

func Push(tag string) error {
	return exec.Command("docker", "push", tag).Run()
}

func Login(hostname, username, password) error {
	return exec.Command(
		"docker", "login",
		hostname,
		"-u", username,
		"-p", password,
	).Run()
}
