package gdblib

import (
	"strings"
)

type ThreadListIdsResult struct {
	ThreadIds       []string `json:"thread-ids"`
	CurrentThreadId string   `json:"current-thread-id"`
	NumThreads      string   `json:"number-of-threads"`
}

func (gdb *GDB) ThreadListIds() (*ThreadListIdsResult, error) {
	descriptor := cmdDescr{}

	descriptor.cmd = "-thread-list-ids"

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response

	// Swap out the thread-ids because they don't work with the normal JSON
	//  mapping
	resultStr := result.result

	beginThreadIds := strings.Index(resultStr, "thread-ids={")
	endThreadIds := strings.Index(resultStr[beginThreadIds:], "}") + beginThreadIds
	threadIds := resultStr[beginThreadIds : endThreadIds+1]

	// Change the object block into an array block and remote the "thread-id" fields to use
	//  those as the string literal
	newThreadIds := strings.Replace(threadIds, "thread-id=", "", -1)
	newThreadIds = strings.Replace(newThreadIds, "{", "[", -1)
	newThreadIds = strings.Replace(newThreadIds, "}", "]", -1)

	resultStr = strings.Replace(resultStr, threadIds, newThreadIds, 1)
	result.result = resultStr

	resultObj := ThreadListIdsResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}

	return &resultObj, nil
}

type ThreadInfoParms struct {
	ThreadId string
}

type ThreadInfoResult struct {
	Threads         []ThreadInfo `json:"threads"`
	CurrentThreadId string       `json:"current-thread-id"`
}

type ThreadInfo struct {
	Id       string    `json:"id"`
	TargetId string    `json:"target-id"`
	Frame    FrameInfo `json:"frame"`
	State    string    `json:"state"`
}

type FrameInfo struct {
	Level    string    `json:"level"`
	Addr     string    `json:"addr"`
	Func     string    `json:"func"`
	Args     []ArgInfo `json:"args"`
	File     string    `json:"file"`
	Fullname string    `json:"fullname"`
	Line     string    `json:"line"`
}

type ArgInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (gdb *GDB) ThreadInfo(parms ThreadInfoParms) (*ThreadInfoResult, error) {
	descriptor := cmdDescr{}

	descriptor.cmd = "-thread-info"
	if parms.ThreadId != "" {
		descriptor.cmd = descriptor.cmd + " " + parms.ThreadId
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response

	resultObj := ThreadInfoResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}

	return &resultObj, nil
}

type ThreadSelectParms struct {
	ThreadId string
}

type ThreadSelectResult struct {
	NewThreadId string `json:"new-thread-id"`
}

func (gdb *GDB) ThreadSelect(parms ThreadSelectParms) (*ThreadSelectResult, error) {
	descriptor := cmdDescr{}

	descriptor.cmd = "-thread-select"
	if parms.ThreadId != "" {
		descriptor.cmd = descriptor.cmd + " " + parms.ThreadId
	}

	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	result := <-descriptor.response

	resultObj := ThreadSelectResult{}
	err := parseResult(result, &resultObj)
	if err != nil {
		return nil, err
	}

	return &resultObj, nil
}
