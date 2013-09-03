// Copyright 2013 Chris McGee <sirnewton_01@yahoo.ca>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package gdblib

import (
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	sendSignalPath string
)

func init() {
	gopath := build.Default.GOPATH
	gopaths := strings.Split(gopath, filepath.ListSeparator)
	for _,path := range(gopaths) {
		p := path + "\\src\\github.com\\sirnewton01\\gdblib\\SendSignal.exe"
		_,err := os.Stat(p)
		if err == nil {
			sendSignalPath = p
			break
		}
	}
}

func fixCmd(cmd *exec.Cmd) {
	// No process group separation is required on Windows.
	// Processes do not share signals like they can on Unix.
}

func interruptInferior(process *os.Process, pid string) {
	// Invoke the included "sendsignal" program to send the
	// Ctrl-break to the inferior process to interrupt it

	initCommand := exec.Command("cmd", "/c", "start", sendSignalPath, pid)
	initCommand.Run()
}
