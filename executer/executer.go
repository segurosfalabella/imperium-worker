package executer

import (
	"os/exec"

	"github.com/sirupsen/logrus"
)

// Commander interface
type Commander interface {
	Run() error
}

var log = logrus.New()

// CreateCommand function
var CreateCommand = func(name string, arg ...string) Commander {
	return exec.Command(name, arg...)
}

// Job struct
type Job struct {
	Name        string
	Description string
	Command     string
	Image       string
	Arguments   string
	Response    string
}

// Execute method
func (job *Job) Execute() error {
	if job.Command == "health" {
		job.Response = "i am alive"
		return nil
	}
	return executeDocker(job)
}

// GetResponse method
func (job *Job) GetResponse() string {
	return job.Response
}

func executeDocker(job *Job) error {
	cmd := CreateCommand("docker", "run", "--rm", job.Image, job.Arguments)
	err := cmd.Run()
	return err
}
