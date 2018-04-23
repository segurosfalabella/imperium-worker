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
}

//Execute method
func (job *Job) Execute() error {
	cmd := CreateCommand("docker", "run", "--rm", "redis")
	err := cmd.Run()
	return err
}
