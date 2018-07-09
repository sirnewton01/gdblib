// Copyright 2013 Chris McGee <sirnewton_01@yahoo.ca>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdblib

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type cmdDescr struct {
	cmd            string
	response       chan cmdResultRecord
	forceInterrupt bool
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

	// Inferior process (if running)
	inferiorLock    sync.Mutex
	inferiorProcess *os.Process
	inferiorPid     string
	inferiorRunning bool

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

// NewGDBWithPID creates a new gdb debugging session.
//  Provide the process ID of the program to debug.
//  The source root directory is optional in order to resolve
//  the source file references.
func NewGDBWithPID(pid int, srcRoot string) (*GDB, error) {
	args := []string{
		"-p", fmt.Sprintf("%d", pid), "--interpreter", "mi2",
	}
	return newGDB(args, srcRoot)
}

// NewGDB creates a new gdb debugging session.
//  Provide the full OS path
//  to the program to debug. The source root directory is optional in
//  order to resolve the source file references.
func NewGDB(program string, srcRoot string) (*GDB, error) {
	args := []string{program, "--interpreter", "mi2"}
	return newGDB(args, srcRoot)
}

// newGDB creates a new gdb debugging session.
//  Provide the arguments to the gdb process incantation.
//  The source root directory is optional in order to resolve
//  the source file references.
func newGDB(cmd []string, srcRoot string) (*GDB, error) {
	gdb := &GDB{}

	gdb.gdbCmd = exec.Command("gdb", cmd...)
	if srcRoot != "" {
		gdb.gdbCmd.Dir = srcRoot
	}

	// Perform any os-specific customizations on the command before launching it
	fixCmd(gdb.gdbCmd)

	gdb.Console = make(chan string)
	gdb.Target = make(chan string)
	gdb.InternalLog = make(chan string)
	gdb.AsyncResults = make(chan AsyncResultRecord)

	gdb.input = make(chan cmdDescr)
	gdb.result = make(chan cmdResultRecord)
	gdb.cmdRegistry = make(map[int64]cmdDescr)
	gdb.nextId = 0

	wg := sync.WaitGroup{}
	wg.Add(3)

	wg2 := sync.WaitGroup{}
	wg2.Add(1)

	writer := func() {
		inPipe, err := gdb.gdbCmd.StdinPipe()

		wg2.Done()

		if err != nil {
			return
		}

		wg.Done()

		// Add a default "main" breakpoint (works in C and Go) to force execution to pause
		//  waiting for user to add breakpoints, etc.
		inPipe.Write([]byte("-break-insert main\n"))

		for {
			select {
			case newInput := <-gdb.input:
				// Interrupt the process so that we can send the command
				gdb.inferiorLock.Lock()
				interrupted := false
				if newInput.forceInterrupt && gdb.inferiorProcess != nil && gdb.inferiorRunning {
					interrupted = true
					interruptInferior(gdb.inferiorProcess, gdb.inferiorPid)
				}
				gdb.inferiorLock.Unlock()

				if newInput.response != nil {
					gdb.nextId++
					id := gdb.nextId
					gdb.cmdRegistry[id] = newInput

					inPipe.Write([]byte(strconv.FormatInt(id, 10) + newInput.cmd + "\n"))
				} else {
					inPipe.Write([]byte(newInput.cmd + "\n"))
				}

				// If it is an empty command then it is because the client is requesting
				//  plain interrupt without continuing.
				if interrupted && newInput.cmd != "" {
					inPipe.Write([]byte("-exec-continue\n"))
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

		if err != nil {
			return
		}

		wg.Done()

		reader := bufio.NewReader(outPipe)
		resultRecordRegex := regexp.MustCompile(`^(\d*)\^(\S+?)(,(.*))?$`)
		asyncRecordRegex := regexp.MustCompile(`^([*=])(\S+?),(.*)$`)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("ERROR: %v\n", err.Error())
				break
			}
			line = strings.Replace(line, "\r", "", -1)
			line = strings.Replace(line, "\n", "", -1)

			if len(line) == 0 {
				continue
			}

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
				//				fmt.Printf("[ASYNC RESULT RECORD] %v %v\n", resultIndication, result)
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

					gdb.inferiorLock.Lock()
					if resultIndication == "thread-group-started" {
						pidStr, ok := resultObj["pid"].(string)

						if ok {
							pid, err := strconv.ParseInt(pidStr, 10, 32)
							if err == nil {
								gdb.inferiorProcess, err = os.FindProcess(int(pid))
								gdb.inferiorPid = pidStr
							}
						}
					} else if resultIndication == "thread-group-exited" {
						gdb.inferiorProcess = nil
					} else if resultIndication == "running" {
						gdb.inferiorRunning = true
					} else if resultIndication == "stopped" {
						gdb.inferiorRunning = false
					}
					gdb.inferiorLock.Unlock()

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

	// Handle standard error as if it comes from the target
	errReader := func() {
		wg2.Wait()
		errPipe, err := gdb.gdbCmd.StderrPipe()

		if err != nil {
			return
		}

		wg.Done()

		reader := bufio.NewReader(errPipe)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("ERROR: %v\n", err.Error())
				break
			}

			gdb.Target <- "[stderr] " + line
		}
	}

	go reader()
	go errReader()
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
	descriptor := cmdDescr{forceInterrupt: true}
	descriptor.cmd = "-gdb-exit"
	descriptor.response = make(chan cmdResultRecord)
	gdb.input <- descriptor
	<-descriptor.response
}

func (gdb *GDB) GdbSet(name, value string) error {
	descriptor := cmdDescr{}
	descriptor.cmd = fmt.Sprintf("-gdb-set %s %s", name, value)
	descriptor.response = make(chan cmdResultRecord)

	gdb.input <- descriptor
	rsp := <-descriptor.response

	return parseResult(rsp, nil)
}

func (gdb *GDB) GdbShow(name string) (string, error) {
	descriptor := cmdDescr{}
	descriptor.cmd = fmt.Sprintf("-gdb-show %s", name)
	descriptor.response = make(chan cmdResultRecord)

	gdb.input <- descriptor
	result := <-descriptor.response

	resultMap := make(map[string]string)
	err := parseResult(result, &resultMap)

	return resultMap["value"], err
}
