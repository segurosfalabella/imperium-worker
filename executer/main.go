package executer

//Job struct
type Job struct {
	Name        string
	Description string
	Command     string
}

//Execute method
func (job Job) Execute() error {
	return nil
}
