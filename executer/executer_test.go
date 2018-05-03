package executer_test

import (
	"errors"
	"testing"

	"github.com/segurosfalabella/imperium-worker/executer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCmd struct {
	mock.Mock
}

func (mockCmd *MockCmd) Run() error {
	args := mockCmd.Called()
	return args.Error(0)
}

func TestExecuteShouldFailWhenRunCommandFail(t *testing.T) {
	job := new(executer.Job)
	oldCreateCommand := executer.CreateCommand
	defer func() { executer.CreateCommand = oldCreateCommand }()
	mock := new(MockCmd)

	mock.On("Run").Return(errors.New("527c090d-4102-4671-9033-b3363f78b343"))
	executer.CreateCommand = func(name string, arg ...string) executer.Commander {
		return mock
	}

	err := job.Execute()

	assert.NotNil(t, err)
	assert.Equal(t, "527c090d-4102-4671-9033-b3363f78b343", err.Error())
}

func TestExecuteShouldSuccessWhenRunCommandSuccess(t *testing.T) {
	job := new(executer.Job)
	oldCreateCommand := executer.CreateCommand
	defer func() { executer.CreateCommand = oldCreateCommand }()
	mock := new(MockCmd)
	mock.On("Run").Return(nil)
	executer.CreateCommand = func(name string, arg ...string) executer.Commander {
		return mock
	}

	err := job.Execute()

	assert.Nil(t, err)
}

func TestShouldHandleHealthCommand(t *testing.T) {
	job := &executer.Job{
		Command: "health",
	}

	err := job.Execute()

	assert.Nil(t, err)
	assert.Equal(t, "i am alive", job.Response)
}

func TestShouldReturnExitCodeZeroWhenExecuteJobSuccess(t *testing.T) {
	job := &executer.Job{
		Image:     "hogwarts/mirror",
		Arguments: "0",
		ExitCode:  -1,
	}
	oldCreateCommand := executer.CreateCommand
	defer func() { executer.CreateCommand = oldCreateCommand }()
	mock := new(MockCmd)
	mock.On("Run").Return(nil)
	executer.CreateCommand = func(name string, arg ...string) executer.Commander {
		return mock
	}

	job.Execute()

	exitCode := 0
	assert.Equal(t, exitCode, job.ExitCode)
}
