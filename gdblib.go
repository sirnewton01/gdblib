// Copyright 2013 Chris McGee <sirnewton_01@yahoo.ca>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdblib

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type cmdDescr struct {
	cmd      string
	response chan cmdResultRecord
}

type cmdResultRecord struct {
	id         int64
	indication string
	result     string
}

type AsyncResultRecord struct {
	Indication string
	Result     map[string]interface{}
}

type GDB struct {
	// Channel of gdb console lines
	Console chan string
	// Channel of target process lines
	Target chan string
	// Channel of internal GDB log lines
	InternalLog chan string
	// Channel of async result records
	AsyncResults chan AsyncResultRecord

	gdbCmd *exec.Cmd

	// Internal channel to send a command to the gdb interpreter
	input chan cmdDescr
	// Internal channel to send result records to callers waiting for a response
	result chan cmdResultRecord

	// Registry of command descriptors for synchronous commands
	cmdRegistry map[int64]cmdDescr
	nextId      int64
}

func convertCString(cstr string) string {
	str := cstr

	if str[0] == '"' && str[len(str)-1] == '"' {
		str = str[1 : len(str)-1]
	}
	str = strings.Replace(str, `\"`, `"`, -1)
	str = strings.Replace(str, `\n`, "\n", -1)

	return str
}

func NewGDB(program string, workingDir string) (*GDB, error) {
	gdb := &GDB{}

	gdb.gdbCmd = exec.Command("gdb", program, "-i=mi")
	gdb.gdbCmd.Dir = workingDir

	gdb.Console = make(chan string)
	gdb.Target = make(chan string)
	gdb.InternalLog = make(chan string)
	gdb.AsyncResults = make(chan AsyncResultRecord)

	gdb.input = make(chan cmdDescr)
	gdb.result = make(chan cmdResultRecord)
	gdb.cmdRegistry = make(map[int64]cmdDescr)
	gdb.nextId = 0

	wg := sync.WaitGroup{}
	wg.Add(2)
	
	wg2 := sync.WaitGroup{}
	wg2.Add(1)

	writer := func() {
		inPipe, err := gdb.gdbCmd.StdinPipe()
		
		wg2.Done()
		
		if err != nil {
			return
		}
		
		wg.Done()
		
		// Force GDB into asynchronous non-stop mode so that it can accept
		//  commands in the middle of execution.
		inPipe.Write([]byte("-gdb-set target-async 1\n"))
		inPipe.Write([]byte("-gdb-set non-stop on\n"))
		
		for {
			select {
			case newInput := <-gdb.input:
				if newInput.response != nil {
					gdb.nextId++
					id := gdb.nextId
					gdb.cmdRegistry[id] = newInput

					inPipe.Write([]byte(strconv.FormatInt(id, 10) + newInput.cmd + "\n"))
				} else {
					inPipe.Write([]byte(newInput.cmd + "\n"))
				}
			case resultRecord := <-gdb.result:
				descriptor := gdb.cmdRegistry[resultRecord.id]

				if descriptor.cmd != "" {
					descriptor.response <- resultRecord
				}
			}
		}
	}

	reader := func() {
		wg2.Wait()
		outPipe, err := gdb.gdbCmd.StdoutPipe()
		gdb.gdbCmd.StderrPipe()
		
		if err != nil {
			return
		}
		
		wg.Done()
		
		reader := bufio.NewReader(outPipe)
		resultRecordRegex := regexp.MustCompile(`^(\d*)\^(\S+?)(,(.*))?$`)
		asyncRecordRegex := regexp.MustCompile(`^([*=])(\S+?),(.*)$`)

		for {
			// TODO what about truncated lines, we should check isPrefix and manage the line better
			lineBytes, _, err := reader.ReadLine()
			if err != nil {
				break
			}

			line := string(lineBytes)

			// TODO unescape the quotes, newlines, etc.
			// stream outputs
			if line[0] == '~' {
				line = convertCString(line[1:])
				gdb.Console <- line
			} else if line[0] == '@' {
				line = convertCString(line[1:])
				gdb.Target <- line
			} else if line[0] == '&' {
				line = convertCString(line[1:])
				gdb.InternalLog <- line + "\n"
				// result record
			} else if matches := resultRecordRegex.FindStringSubmatch(line); matches != nil {
				commandId := matches[1]
				resultIndication := matches[2]
				result := ""
				if len(matches) > 4 {
					result = matches[4]
				}

				if commandId != "" {
					id, err := strconv.ParseInt(commandId, 10, 64)

					if err == nil {
						resultRecord := cmdResultRecord{id: id, indication: resultIndication, result: result}
						gdb.result <- resultRecord
					}

					// TODO handle the parse error case
				}
				//				else {
				//					fmt.Printf("[RESULT RECORD] ID:%v %v %v\n", commandId, resultIndication, result)
				//				}
				// async record
			} else if matches := asyncRecordRegex.FindStringSubmatch(line); matches != nil {
				// recordType := matches[1]
				resultIndication := matches[2]
				result := matches[3]

				resultNode, _ := createObjectNode("{" + result + "}")
				resultObj := make(map[string]interface{})
				jsonStr := resultNode.toJSON()
				err := json.Unmarshal([]byte(jsonStr), &resultObj)

				if err == nil {
					resultRecord := AsyncResultRecord{Indication: resultIndication, Result: resultObj}
					gdb.AsyncResults <- resultRecord
				} else {
					fmt.Printf("[ORIGINAL] %v\n", result)
					fmt.Printf("[JSON] %v\n", jsonStr)
					fmt.Printf("Error unmarshalling JSON for async result record: %v %v\n", err.Error(), resultNode.toJSON())
				}
				// TODO handle the parse error case
				//				fmt.Printf("[ASYNC RESULT RECORD] %v %v\n", resultIndication, result)
			} else if line == "(gdb) " {
				// This is the gdb prompt. We can just throw it out
			} else {
				//fmt.Printf("%v\n", line)
				gdb.Target <- line + "\n"
			}
		}
	}

	go reader()
	go writer()
	
	wg.Wait()

	err := gdb.gdbCmd.Start()
	if err != nil {
		return nil, err
	}

	return gdb, nil
}

func (gdb *GDB) Wait() error {
	return gdb.gdbCmd.Wait()
}

func parseResult(result cmdResultRecord, resultObj interface{}) error {
	if result.indication == "error" {
		msg := strings.Replace(result.result, `msg="`, "", 1)
		msg = msg[:len(msg)-1]

		return errors.New(msg)
	}

	if resultObj != nil {
		//		fmt.Printf("[ORIGINAL] %v\n", result.result)

		gdbNode, _ := createObjectNode("{" + result.result + "}")
		jsonStr := gdbNode.toJSON()

		//		fmt.Printf("[JSON DUMP] %v\n", jsonStr)

		err := json.Unmarshal([]byte(jsonStr), &resultObj)
		if err != nil {
			fmt.Printf("[ORIGINAL] %v\n", result.result)
			fmt.Printf("[JSON DUMP] %v\n", jsonStr)
			return err
		}
	}

	return nil
}

func (gdb *GDB) GdbExit() {
	descriptor := cmdDescr{}
	descriptor.cmd = "-gdb-exit"
	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	<-descriptor.response
}
