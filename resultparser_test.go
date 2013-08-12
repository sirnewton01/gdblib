package gdblib

import (
	"encoding/json"
	"testing"

//	"fmt"
)

func TestBreakListOutput1(t *testing.T) {
	input := `BreakpointTable={nr_rows="0",nr_cols="6",hdr=[{width="7",alignment="-1",col_name="number",colhdr="Num"},{width="14",alignment="-1",col_name="type",colhdr="Type"},{width="4",alignment="-1",col_name="disp",colhdr="Disp"},{width="3",alignment="-1",col_name="enabled",colhdr="Enb"},{width="10",alignment="-1",col_name="addr",colhdr="Address"},{width="40",alignment="2",col_name="what",colhdr="What"}],body=[]}`
	result, _ := createKeyValueNode(input)

	jsonObj1 := make(map[string]interface{})
	err := json.Unmarshal([]byte("{"+result.toJSON()+"}"), &jsonObj1)
	if err != nil {
		t.Error(err)
	}
	jsonOutput1, err := json.Marshal(jsonObj1)
	if err != nil {
		t.Error(err)
	}

	jsonObj2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(`{"BreakpointTable":{"nr_rows":"0","nr_cols":"6","hdr":[{"width":"7","alignment":"-1","col_name":"number","colhdr":"Num"},{"width":"14","alignment":"-1","col_name":"type","colhdr":"Type"},{"width":"4","alignment":"-1","col_name":"disp","colhdr":"Disp"},{"width":"3","alignment":"-1","col_name":"enabled","colhdr":"Enb"},{"width":"10","alignment":"-1","col_name":"addr","colhdr":"Address"},{"width":"40","alignment":"2","col_name":"what","colhdr":"What"}],"body":[]}}`), &jsonObj2)
	if err != nil {
		t.Error(err)
	}
	jsonOutput2, err := json.Marshal(jsonObj2)
	if err != nil {
		t.Error(err)
	}

	if string(jsonOutput1) != string(jsonOutput2) {
		t.Error("Breakpoint list is not output does not match exemplar")
	}
}

func TestBreakpointHit1(t *testing.T) {
	input := `reason="breakpoint-hit",disp="keep",bkptno="1",frame={addr="0x0000000000400c00",func="main.printHello",args=[],file="/home/cmcgee/godev/src/hello/hello.go",fullname="/home/cmcgee/godev/src/hello/hello.go",line="8"},thread-id="2",stopped-threads=["2"],core="3"`
	result, _ := createObjectNode("{" + input + "}")

	jsonObj1 := make(map[string]interface{})
	err := json.Unmarshal([]byte(result.toJSON()), &jsonObj1)
	if err != nil {
		t.Error(err)
	}
	_, err = json.Marshal(jsonObj1)
	if err != nil {
		t.Error(err)
	}

	//	fmt.Printf("JSON:%v\n", string(jsonOutput1))

	//	jsonObj2 := make(map[string]interface{})
	//	err = json.Unmarshal([]byte(`{"BreakpointTable":{"nr_rows":"0","nr_cols":"6","hdr":[{"width":"7","alignment":"-1","col_name":"number","colhdr":"Num"},{"width":"14","alignment":"-1","col_name":"type","colhdr":"Type"},{"width":"4","alignment":"-1","col_name":"disp","colhdr":"Disp"},{"width":"3","alignment":"-1","col_name":"enabled","colhdr":"Enb"},{"width":"10","alignment":"-1","col_name":"addr","colhdr":"Address"},{"width":"40","alignment":"2","col_name":"what","colhdr":"What"}],"body":[]}}`), &jsonObj2)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//	jsonOutput2,err := json.Marshal(jsonObj2)
	//	if err != nil {
	//		t.Error(err)
	//	}
	//
	//	if string(jsonOutput1) != string(jsonOutput2) {
	//		t.Error("Breakpoint list is not output does not match exemplar")
	//	}
}

func TestBreakListOutput2(t *testing.T) {
	input := `BreakpointTable={nr_rows="2",nr_cols="6",hdr=[{width="7",alignment="-1",col_name="number",colhdr="Num"},{width="14",alignment="-1",col_name="type",colhdr="Type"},{width="4",alignment="-1",col_name="disp",colhdr="Disp"},{width="3",alignment="-1",col_name="enabled",colhdr="Enb"},{width="18",alignment="-1",col_name="addr",colhdr="Address"},{width="40",alignment="2",col_name="what",colhdr="What"}],body=[{number="1",type="breakpoint",disp="keep",enabled="y",addr="0x0000000000400c3d",func="main.main",file="/home/cmcgee/godev/src/hello/hello.go",fullname="/home/cmcgee/godev/src/hello/hello.go",line="12",times="0",original-location="main.main"},{number="2",type="breakpoint",disp="keep",enabled="y",addr="0x0000000000400c00",func="main.printHello",file="/home/cmcgee/godev/src/hello/hello.go",fullname="/home/cmcgee/godev/src/hello/hello.go",line="8",times="0",original-location="main.printHello"}]}`
	result, _ := createKeyValueNode(input)

	jsonObj1 := make(map[string]interface{})
	err := json.Unmarshal([]byte("{"+result.toJSON()+"}"), &jsonObj1)
	if err != nil {
		t.Error(err)
	}
	jsonOutput1, err := json.Marshal(jsonObj1)
	if err != nil {
		t.Error(err)
	}

	jsonObj2 := make(map[string]interface{})
	err = json.Unmarshal([]byte(`{"BreakpointTable":{"nr_rows":"2","nr_cols":"6","hdr":[{"width":"7","alignment":"-1","col_name":"number","colhdr":"Num"},{"width":"14","alignment":"-1","col_name":"type","colhdr":"Type"},{"width":"4","alignment":"-1","col_name":"disp","colhdr":"Disp"},{"width":"3","alignment":"-1","col_name":"enabled","colhdr":"Enb"},{"width":"18","alignment":"-1","col_name":"addr","colhdr":"Address"},{"width":"40","alignment":"2","col_name":"what","colhdr":"What"}],"body":[{"number":"1","type":"breakpoint","disp":"keep","enabled":"y","addr":"0x0000000000400c3d","func":"main.main","file":"/home/cmcgee/godev/src/hello/hello.go","fullname":"/home/cmcgee/godev/src/hello/hello.go","line":"12","times":"0","original-location":"main.main"},{"number":"2","type":"breakpoint","disp":"keep","enabled":"y","addr":"0x0000000000400c00","func":"main.printHello","file":"/home/cmcgee/godev/src/hello/hello.go","fullname":"/home/cmcgee/godev/src/hello/hello.go","line":"8","times":"0","original-location":"main.printHello"}]}}`), &jsonObj2)
	if err != nil {
		t.Error(err)
	}
	jsonOutput2, err := json.Marshal(jsonObj2)
	if err != nil {
		t.Error(err)
	}

	if string(jsonOutput1) != string(jsonOutput2) {
		t.Error("Breakpoint list is not output does not match exemplar")
	}
}

func TestString(t *testing.T) {
	// TRIVIAL
	input := `""`
	result, size := createStringNode(input)

	if result != `""` {
		t.Errorf("String equal to '%v' instead of '\"\"'", result)
	}

	if size != 2 {
		t.Errorf("String size equal to '%v' instead of '2'", size)
	}

	// REASONABLE
	input = `"value"`
	result, size = createStringNode(input)

	if result != `"value"` {
		t.Errorf("String equal to '%v' instead of '\"value\"'", result)
	}

	if size != 7 {
		t.Errorf("String size equal to '%v' instead of '7'", size)
	}

	// EXTRA STUFF AT THE END
	input = `"value",[]{}`
	result, size = createStringNode(input)

	if result != `"value"` {
		t.Errorf("String equal to '%v' instead of '\"value\"'", result)
	}

	if size != 7 {
		t.Errorf("String size equal to '%v' instead of '7'", size)
	}

	// ESCAPED QUOTE
	input = `"val\"ue"`
	result, size = createStringNode(input)
	if result != `"val\"ue"` {
		t.Errorf("String equal to '%v' instead of '\"val\\\"ue\"'", result)
	}

	if size != 9 {
		t.Errorf("String size equal to '%v' instead of '9'", size)
	}
}

func TestObject(t *testing.T) {
	// TRIVIAL
	input := `{}`
	result, size := createObjectNode(input)

	if len(result.children) != 0 {
		t.Errorf("Number of children equal to '%v' instead of '0'", len(result.children))
	}

	if size != 2 {
		t.Errorf("Size equal to '%v' instead of '2'", size)
	}

	// SINGLE VALUE
	input = `{key1="value1"}`
	result, size = createObjectNode(input)

	if len(result.children) != 1 {
		t.Errorf("Number of children equal to '%v' instead of '1'", len(result.children))
	}

	if size != 15 {
		t.Errorf("Size equal to '%v' instead of '29'", size)
	}

	// MULTIPLE VALUES
	input = `{key1="value1",key2="value2"}`
	result, size = createObjectNode(input)

	if len(result.children) != 2 {
		t.Errorf("Number of children equal to '%v' instead of '2'", len(result.children))
	}

	if size != 29 {
		t.Errorf("Size equal to '%v' instead of '29'", size)
	}

	// EXTRA STUFF
	input = `{key1="value1",key2="value2"},[]{}`
	result, size = createObjectNode(input)

	if len(result.children) != 2 {
		t.Errorf("Number of children equal to '%v' instead of '2'", len(result.children))
	}

	if size != 29 {
		t.Errorf("Size equal to '%v' instead of '29'", size)
	}
}

func TestArray(t *testing.T) {
	// TRIVIAL
	input := `[]`
	result, size := createArrayNode(input)

	if len(result.children) != 0 {
		t.Errorf("Number of children equal to '%v' instead of '0'", len(result.children))
	}

	if size != 2 {
		t.Errorf("Size equal to '%v' instead of '2'", size)
	}

	// SINGLE KEYED VALUE
	input = `[key1="value1"]`
	result, size = createArrayNode(input)

	if len(result.children) != 1 {
		t.Errorf("Number of children equal to '%v' instead of '1'", len(result.children))
	}

	if size != 15 {
		t.Errorf("Size equal to '%v' instead of '29'", size)
	}

	// MULTIPLE KEYED VALUES
	input = `[key1="value1",key2="value2"]`
	result, size = createArrayNode(input)

	if len(result.children) != 2 {
		t.Errorf("Number of children equal to '%v' instead of '2'", len(result.children))
	}

	if size != 29 {
		t.Errorf("Size equal to '%v' instead of '29'", size)
	}

	// EXTRA STUFF
	input = `[key1="value1",key2="value2"],[]{}`
	result, size = createArrayNode(input)

	if len(result.children) != 2 {
		t.Errorf("Number of children equal to '%v' instead of '2'", len(result.children))
	}

	if size != 29 {
		t.Errorf("Size equal to '%v' instead of '29'", size)
	}
}

func TestKeyValue1(t *testing.T) {
	// TRIVIAL 1
	input := `""`
	result, size := createKeyValueNode(input)

	if result.key != "" {
		t.Errorf("Key equal to '%v' instead of ''", result.key)
	}

	if result.value != `""` {
		t.Errorf("Value equal to '%v' instead of ''", result.value)
	}

	if size != 2 {
		t.Errorf("Size equal to '%v' instead of '2'", size)
	}

	// TRIVIAL 2
	input = `key={}`
	result, size = createKeyValueNode(input)

	if result.key != "key" {
		t.Errorf("Key equal to '%v' instead of 'key'", result.key)
	}

	//	if result.value != `""` {
	//		t.Errorf("Value equal to '%v' instead of ''", result.value)
	//	}

	if size != 6 {
		t.Errorf("Size equal to '%v' instead of '6'", size)
	}

	// TRIVIAL 3
	input = `key=[]`
	result, size = createKeyValueNode(input)

	if result.key != "key" {
		t.Errorf("Key equal to '%v' instead of 'key'", result.key)
	}

	//	if result.value != `""` {
	//		t.Errorf("Value equal to '%v' instead of ''", result.value)
	//	}

	if size != 6 {
		t.Errorf("Size equal to '%v' instead of '6'", size)
	}

	// REASONABLE
	input = `key="value"`
	result, size = createKeyValueNode(input)

	if result.key != "key" {
		t.Errorf("Key equal to '%v' instead of 'key'", result.key)
	}

	if result.value != `"value"` {
		t.Errorf("Value equal to '%v' instead of 'value'", result.value)
	}

	if size != 11 {
		t.Errorf("Size equal to '%v' instead of '11'", size)
	}

	// REASONABLE
	input = `key="value",{}[]`
	result, size = createKeyValueNode(input)

	if result.key != "key" {
		t.Errorf("Key equal to '%v' instead of 'key'", result.key)
	}

	if result.value != `"value"` {
		t.Errorf("Value equal to '%v' instead of 'value'", result.value)
	}

	if size != 11 {
		t.Errorf("Size equal to '%v' instead of '11'", size)
	}
}

func TestArrayChildKeys(t *testing.T) {
	input := `[foo="bar",foo="baz"]`
	result, _ := createArrayNode(input)

	for _, child := range result.children {
		if child.value.(gdbResultNode).key != "" {
			t.Errorf("Key value node underneath an array node has a key")
		}
	}
}
