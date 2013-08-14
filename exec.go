// Copyright 2013 Chris McGee <sirnewton_01@yahoo.ca>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdblib

import ()

type ExecRunParms struct {
	ThreadGroup  string
	AllInferiors bool
}

func (gdb *GDB) ExecRun(parms ExecRunParms) error {
	descriptor := cmdDescr{}

	descriptor.cmd = "-exec-run"
	if parms.AllInferiors {
		descriptor.cmd = descriptor.cmd + " --all"
	} else if parms.ThreadGroup != "" {
		descriptor.cmd = descriptor.cmd + " --thread-group " + parms.ThreadGroup
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor

	result := <-descriptor.response
	err := parseResult(result, nil)

	return err
}

type ExecInterruptParms struct {
	ThreadGroup  string
	AllInferiors bool
}

func (gdb *GDB) ExecInterrupt(parms ExecInterruptParms) error {
	descriptor := cmdDescr{}

	descriptor.cmd = "-exec-interrupt"
	if parms.AllInferiors {
		descriptor.cmd = descriptor.cmd + " --all"
	} else if parms.ThreadGroup != "" {
		descriptor.cmd = descriptor.cmd + " --thread-group " + parms.ThreadGroup
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response
	err := parseResult(result, nil)

	return err
}

type ExecNextParms struct {
	Reverse bool
}

func (gdb *GDB) ExecNext(parms ExecNextParms) error {
	descriptor := cmdDescr{}

	descriptor.cmd = "-exec-next"
	if parms.Reverse {
		descriptor.cmd = descriptor.cmd + " --reverse"
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response
	err := parseResult(result, nil)

	return err
}

type ExecStepParms struct {
	Reverse bool
}

func (gdb *GDB) ExecStep(parms ExecStepParms) error {
	descriptor := cmdDescr{}

	descriptor.cmd = "-exec-step"
	if parms.Reverse {
		descriptor.cmd = descriptor.cmd + " --reverse"
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response
	err := parseResult(result, nil)

	return err
}

type ExecContinueParms struct {
	Reverse      bool
	ThreadGroup  string
	AllInferiors bool
}

func (gdb *GDB) ExecContinue(parms ExecContinueParms) error {
	descriptor := cmdDescr{}

	descriptor.cmd = "-exec-continue"

	if parms.Reverse {
		descriptor.cmd = descriptor.cmd + " --reverse"
	}
	if parms.AllInferiors {
		descriptor.cmd = descriptor.cmd + " --all"
	} else if parms.ThreadGroup != "" {
		descriptor.cmd = descriptor.cmd + " --thread-group " + parms.ThreadGroup
	}
	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response
	err := parseResult(result, nil)

	return err
}
