package executer

type Job struct {
	Name        string
	Description string
	Command     string
}

func (job Job) Execute() (bool, error) {
	return true, nil
}
