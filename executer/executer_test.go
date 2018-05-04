package executer_test

import (
	"syscall"
	"testing"

	"github.com/segurosfalabella/imperium-worker/executer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCmd struct {
	mock.Mock
}

type ExitErrorDouble struct {
	mock.Mock
}

func (e *ExitErrorDouble) Error() string {
	return ""
}

// Sys function
func (e *ExitErrorDouble) Sys() interface{} {
	args := e.Called()
	return args.Get(0).(syscall.WaitStatus)
	// return uint32(args.Int(0))
}

func (mockCmd *MockCmd) Run() executer.ExitErrorInterface {
	args := mockCmd.Called()
	return args.Get(0).(executer.ExitErrorInterface)
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
	exitErrorDouble := new(ExitErrorDouble)
	exitErrorDouble.On("Sys").Return(syscall.WaitStatus(0))
	mock := new(MockCmd)
	mock.On("Run").Return(exitErrorDouble)
	executer.CreateCommand = func(name string, arg ...string) executer.Commander {
		return mock
	}

	job.Execute()

	exitCode := 0
	assert.Equal(t, exitCode, job.ExitCode)
}

func TestShouldReturnExitCodeNotZeroWhenExecuteJobFail(t *testing.T) {
	job := &executer.Job{
		Image:     "hogwarts/mirror",
		Arguments: "0",
		ExitCode:  -1,
	}
	oldCreateCommand := executer.CreateCommand
	defer func() { executer.CreateCommand = oldCreateCommand }()
	exitErrorDouble := new(ExitErrorDouble)
	exitErrorDouble.On("Sys").Return(syscall.WaitStatus(00003000))
	mock := new(MockCmd)
	mock.On("Run").Return(exitErrorDouble)
	executer.CreateCommand = func(name string, arg ...string) executer.Commander {
		return mock
	}

	job.Execute()

	exitCode := 6
	assert.Equal(t, exitCode, job.ExitCode)
}
