package gdblib

import(
	"strconv"
)

type BreakListResult struct {
	BreakPointTable BreakPointTable
}

type BreakPointTable struct {
	Nr_rows string					`json:"nr_rows"`
	Nr_cols string					`json:"nr_cols"`
	Hdr []BreakPointHeaderElement	`json:"hdr"`
	Body []BreakPoint				`json:"body"`
}

type BreakPointHeaderElement struct {
	Width string					`json:"width"`
	Alignment string				`json:"alignment"`
	Col_name string					`json:"col_name"`
	Colhdr string					`json:"colhdr"`
}

type BreakPoint struct {
	Number string					`json:"number"`
	Type string						`json:"type"`
	FullName string					`json:"fullname"`
	Disp string						`json:"disp"`
	Enabled string					`json:"enabled"`
	Addr string						`json:"addr"`
	Func string						`json:"func"`
	File string						`json:"file"`
	Line string						`json:"line"`
	ThreadGroups []string			`json:"thread-groups"`
	Times string					`json:"times"`
}

func (gdb *GDB) BreakList() (breakList *BreakListResult, _ error){
	descriptor := cmdDescr{}
	
	descriptor.cmd = "-break-list"
	
	descriptor.response = make (chan cmdResultRecord)
	gdb.input <- descriptor
	
	result := <- descriptor.response
	
	resultObj := BreakListResult{}
	err := parseResult(result, &resultObj)
	
	if err != nil {
		return nil, err
	}
	
	return &resultObj, nil
}

type BreakInsertParms struct {
	Temporary bool
	Hardware bool
	Force bool
	Disabled bool
	Tracepoint bool
	Condition string
	IgnoreCount int64
	ThreadId string
	Location string
}

type BreakInsertResult struct {
	BreakPoint BreakPoint	`json:"bkpt"`
}

func (gdb *GDB) BreakInsert(parms BreakInsertParms) (*BreakInsertResult, error) {
	descriptor := cmdDescr{}
	
	descriptor.cmd = "-break-insert"
	if parms.Temporary {
		descriptor.cmd = descriptor.cmd + " -t"
	}
	if parms.Hardware {
		descriptor.cmd = descriptor.cmd + " -h"
	}
	if (parms.Force) {
		descriptor.cmd = descriptor.cmd + " -f"
	}
	if (parms.Disabled) {
		descriptor.cmd = descriptor.cmd + " -d"
	}
	if (parms.Tracepoint) {
		descriptor.cmd = descriptor.cmd + " -a"
	}
	if (parms.Condition != "") {
		descriptor.cmd = descriptor.cmd + " -c " + parms.Condition
	}
	if (parms.IgnoreCount > 0) {
		descriptor.cmd = descriptor.cmd + " -i " + strconv.FormatInt(parms.IgnoreCount, 10)
	}
	if (parms.ThreadId != "") {
		descriptor.cmd = descriptor.cmd + " -p " + parms.ThreadId
	}
	if (parms.Location != "") {
		descriptor.cmd = descriptor.cmd + " " + parms.Location
	}
	
	descriptor.response = make (chan cmdResultRecord)
	gdb.input <- descriptor
	
	result := <- descriptor.response
	resultObj := BreakInsertResult{}
	
	err := parseResult(result, &resultObj)
	
	if (err !=  nil) {
		return nil, err
	}
	
	return &resultObj, nil
}

type BreakEnableParms struct {
	Breakpoints []string
}

func (gdb *GDB) BreakEnable(parms BreakEnableParms) (_ error){
	descriptor := cmdDescr{}
	
	descriptor.cmd = "-break-enable"
	
	for _, id := range(parms.Breakpoints) {
		descriptor.cmd = descriptor.cmd + " " + id
	}
	
	descriptor.response = make (chan cmdResultRecord)
	gdb.input <- descriptor
	
	result := <- descriptor.response
	
	err := parseResult(result, nil)
	
	if err != nil {
		return err
	}
	
	return nil
}

type BreakDisableParms struct {
	Breakpoints []string
}

func (gdb *GDB) BreakDisable(parms BreakDisableParms) (_ error) {
	descriptor := cmdDescr{}
	
	descriptor.cmd = "-break-disable"
	
	for _, id := range(parms.Breakpoints) {
		descriptor.cmd = descriptor.cmd + " " + id
	}
	
	descriptor.response = make (chan cmdResultRecord)
	gdb.input <- descriptor
	
	result := <- descriptor.response
	
	err := parseResult(result, nil)

	if err != nil {
		return err
	}
	
	return nil
}
