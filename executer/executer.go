package executer

import (
	"encoding/json"
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
	ID        string
	Command   string
	Image     string
	Arguments string
	Envs      map[string]string
	Response  string
	ExitCode  int
}

// FromJSON method
func (job *Job) FromJSON(text string) {
	json.Unmarshal([]byte(text), &job)
}

// ToJSON method
func (job *Job) ToJSON() string {
	binary, _ := json.Marshal(job)
	return string(binary)
}

// Execute method
func (job *Job) Execute() error {
	if job.Command == "health" {
		job.Response = "i am alive"
		return nil
	}
	return executeDocker(job)
}

func executeDocker(job *Job) error {
	cmd := CreateCommand("docker", "run", "--rm", job.Image, job.Arguments)
	err := cmd.Run()
	setExitCode(job, err)
	return err
}

func setExitCode(job *Job, err error) {
	if err == nil {
		job.ExitCode = 0
	}
}
