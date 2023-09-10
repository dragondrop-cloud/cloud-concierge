package pyscriptexec

// PyScriptExec is an interface for running python scripts programmatically.
type PyScriptExec interface {

	// ExecutePythonScript is a generic function for executing a python script
	// from within the python_scripts directory.
	ExecutePythonScript(name string, otherArgs []string) error

	// RunStateOfCloudReport is a function that wraps ExecutePythonScript to execute
	// python_scripts/state_of_cloud_report/main.py
	RunStateOfCloudReport(uniqueID string, jobName string) error
}

// pyScriptExec implements the PyScriptExec interface.
type pyScriptExec struct {
}

// NewPyScriptExec returns an instance of the PyScriptExec interface.
func NewPyScriptExec() PyScriptExec {
	return &pyScriptExec{}
}
