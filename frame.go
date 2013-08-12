package gdblib

import ()

type StackInfoFrameResult struct {
	Frame Frame `json:"frame"`
}

type Frame struct {
	Level    string `json:"level"`
	Addr     string `json:"addr"`
	Func     string `json:"func"`
	File     string `json:"file"`
	Fullname string `json:"fullname"`
	Line     string `json:"line"`
	From     string `json:"from"`
}

func (gdb *GDB) StackInfoFrame() (*StackInfoFrameResult, error) {
	descriptor := cmdDescr{}

	descriptor.cmd = "-stack-info-frame"

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response

	resultObj := StackInfoFrameResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}

	return &resultObj, nil
}

type StackListFramesParms struct {
	NoFrameFilters bool
	LowFrame       string
	HighFrame      string
}

type StackListFramesResult struct {
	Stack []Frame `json:"stack"`
}

func (gdb *GDB) StackListFrames(parms StackListFramesParms) (*StackListFramesResult, error) {
	descriptor := cmdDescr{}

	descriptor.cmd = "-stack-list-frames"
	if parms.NoFrameFilters {
		descriptor.cmd = descriptor.cmd + " --no-frame-filters"
	}
	if parms.LowFrame != "" && parms.HighFrame != "" {
		descriptor.cmd = descriptor.cmd + " " + parms.LowFrame + " " + parms.HighFrame
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response

	resultObj := StackListFramesResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}

	return &resultObj, nil
}

type StackListVariablesParms struct {
	AllValues bool
	Thread    string
	Frame     string
}

type StackListVariablesResult struct {
	Variables []Variable `json:"variables"`
}

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (gdb *GDB) StackListVariables(parms StackListVariablesParms) (*StackListVariablesResult, error) {
	descriptor := cmdDescr{}

	descriptor.cmd = "-stack-list-variables"
	descriptor.cmd = descriptor.cmd + " --thread " + parms.Thread
	descriptor.cmd = descriptor.cmd + " --frame " + parms.Frame
	if parms.AllValues {
		descriptor.cmd = descriptor.cmd + " --all-values"
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response

	resultObj := StackListVariablesResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}

	return &resultObj, nil
}
