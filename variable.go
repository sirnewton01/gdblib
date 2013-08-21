// Copyright 2013 Chris McGee <sirnewton_01@yahoo.ca>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdblib

import (
)

type VarCreateParms struct {
	// Name for the new variable or empty for gdb to assign a new unique variable name.
	Name string
	
	// Frame address can be a valid frame id, "*" (current frame), or "@" (floating variable)
	// An empty frame address defaults to the current frame.
	FrameAddr string
	
	// Expression to assign this variable
	Expression string
}

type VarCreateResult struct {
	Name string `json:"name"`
	NumChild string `json:"numchild"`
	Value string `json:"value"`
	Type string `json:"type"`
	ThreadId string `json:"thread-id"`
	HasMore string `json:"has_more"`
}

func (gdb *GDB) VarCreate(parms VarCreateParms) (*VarCreateResult, error) {
	descriptor := cmdDescr{}
	
	descriptor.cmd = "-var-create"
	if parms.Name != "" {
		descriptor.cmd = descriptor.cmd + " " + parms.Name
	} else {
		descriptor.cmd = descriptor.cmd + " -"
	}
	
	if parms.FrameAddr == "" {
		descriptor.cmd = descriptor.cmd + " *"
	} else {
		descriptor.cmd = descriptor.cmd + " " + parms.FrameAddr
	}

	descriptor.cmd = descriptor.cmd + " " + parms.Expression
	
	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response
	
	resultObj := VarCreateResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}
	
	return &resultObj, nil
}

type VarDeleteParms struct {
	Name string
	ChildrenOnly bool
}

func (gdb *GDB) VarDelete(parms VarDeleteParms) (error) {
	descriptor := cmdDescr{}
	
	descriptor.cmd = "-var-delete"
	
	if parms.ChildrenOnly {
		descriptor.cmd = descriptor.cmd + " -c"
	}
	
	descriptor.cmd = descriptor.cmd + " " + parms.Name
	
	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response
	
	err := parseResult(result, nil)
	
	return err
}

type VarListChildrenParms struct {
	Name string
	AllValues bool
	From string
	To string
}

type VarListChildrenResult struct {
	NumChild string `json:"num-child"`
	Children []ChildVar `json:"children"`
}

type ChildVar struct {
	Name string `json:"name"`
	Exp string `json:"exp"`
	NumChild string `json:"numchild"`
	Type string `json:"type"`
	Value string `json:"value"`
	ThreadId string `json:"thread-id"`
	Frozen string `json:"frozen"`
}

func (gdb *GDB) VarListChildren(parms VarListChildrenParms) (*VarListChildrenResult, error) {
	descriptor := cmdDescr{}
	
	descriptor.cmd = "-var-list-children"
	if parms.AllValues {
		descriptor.cmd = descriptor.cmd + " --all-values"
	}
	descriptor.cmd = descriptor.cmd + " " + parms.Name
	if parms.From != "" && parms.To != "" {
		descriptor.cmd = descriptor.cmd + " " + parms.From + " " + parms.To
	}
	
	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response
	
	resultObj := VarListChildrenResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}
	
	return &resultObj, nil
}
