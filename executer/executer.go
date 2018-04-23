package executer

import "os/exec"

// Commander interface
type Commander interface {
	Run() error
}

// CreateCommand function
var CreateCommand = func(name string, arg ...string) Commander {
	return exec.Command(name, arg...)
}

//Job struct
type Job struct {
	Name        string
	Description string
	Command     string
	Image       string
	Arguments   string
}

//Execute method
func (job *Job) Execute() error {
	cmd := CreateCommand("docker", "run", "--rm", job.Image, job.Arguments)
	err := cmd.Run()
	return err
}
