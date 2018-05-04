package executer

import (
	"encoding/json"
	"os/exec"
	"syscall"
)

// Commander interface
type Commander interface {
	Run() ExitErrorInterface
}

// ExitErrorInterface interface
type ExitErrorInterface interface {
	Error() string
	Sys() interface{}
}

// CmdShim struct
type CmdShim struct {
	*exec.Cmd
}

// Run function
func (c *CmdShim) Run() ExitErrorInterface {
	err := c.Cmd.Run()
	if exitError, ok := err.(*exec.ExitError); ok {
		return &ExitErrorShim{exitError}
	}

	return nil
}

// ExitErrorShim struct
type ExitErrorShim struct {
	*exec.ExitError
}

func (e *ExitErrorShim) Error() string {
	return e.ProcessState.String()
}

// Sys function
func (e *ExitErrorShim) Sys() interface{} {
	return e.ExitError.Sys()
}

// CreateCommand function
var CreateCommand = func(name string, arg ...string) Commander {
	cmd := exec.Command(name, arg...)
	return &CmdShim{cmd}
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
	if exitError, ok := err.(ExitErrorInterface); ok {
		waitStatus := exitError.Sys().(syscall.WaitStatus)
		job.ExitCode = waitStatus.ExitStatus()
	}
}
