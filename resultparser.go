// Copyright 2013 Chris McGee <sirnewton_01@yahoo.ca>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdblib

import (
//"fmt"
)

type gdbResultNode struct {
	children []gdbResultNode

	key      string
	value    interface{}
	nodeType string
}

func createObjectNode(input string) (gdbResultNode, int) {
	//	fmt.Printf("OBJECT\n")

	node := gdbResultNode{nodeType: "object"}

	i := 1

	for ; i < len(input); i++ {
		c := input[i]

		if c == '}' {
			break
		} else if c == ',' {
			// skip and continue to the next key value node
		} else if c == ' ' || c == '\t' || c == '\r' || c == '\n' {
			// Whitespace, skip and proceed to the next character
		} else {
			value, size := createKeyValueNode(input[i:])
			i = i + size - 1
			node.children = append(node.children, value)
		}
	}

	//	fmt.Printf("OBJECT : %v %v\n", input[:i+1], i+1)
	return node, i + 1
}

func createStringNode(input string) (string, int) {
	//	fmt.Printf("STRING\n")

	i := 1
	buffer := `"`

	for ; i < len(input); i++ {
		c := input[i]

		buffer = buffer + string(c)

		if c == '"' && (i <= 1 || input[i-1] != '\\') {
			break
		}
	}

	//	fmt.Printf("STRING : %v %v\n", input[:i+1], i+1)
	return buffer, i + 1
}

func createArrayNode(input string) (gdbResultNode, int) {
	//	fmt.Printf("ARRAY\n")

	node := gdbResultNode{nodeType: "array"}

	i := 1

	for ; i < len(input); i++ {
		c := input[i]

		if c == '"' {
			str, size := createStringNode(input[i:])
			i = i + size - 1
			childNode := gdbResultNode{nodeType: "keyvalue"}
			childNode.value = str
			node.children = append(node.children, childNode)
		} else if c == '{' {
			objectNode, size := createObjectNode(input[i:])
			i = i + size - 1
			childNode := gdbResultNode{nodeType: "keyvalue"}
			childNode.value = objectNode
			node.children = append(node.children, childNode)
		} else if c == '\t' || c == '\n' || c == ' ' || c == '\r' {
			// Whitespace, ignore it
		} else if c == ',' {
			// A new array element is beginning
		} else if c == ']' {
			break
		} else {
			valueNode, size := createKeyValueNode(input[i:])
			valueNode.key = ""
			i = i + size - 1
			childNode := gdbResultNode{nodeType: "keyvalue"}
			childNode.value = valueNode
			node.children = append(node.children, childNode)
		}
	}

	//	fmt.Printf("ARRAY : %v %v\n", input[:i+1], i+1)
	return node, i + 1
}

func createKeyValueNode(input string) (gdbResultNode, int) {
	//	fmt.Printf("KEYVALUE\n")

	node := gdbResultNode{nodeType: "keyvalue"}

	buffer := ""

	i := 0

	for ; i < len(input); i++ {
		c := input[i]

		if c == '=' {
			node.key = buffer
			buffer = ""
		} else if c == '{' {
			// Beginning of an object
			objectNode, size := createObjectNode(input[i:])
			i = i + size - 1
			node.value = objectNode
			break
		} else if c == '[' {
			// Beginning of an array
			arrayNode, size := createArrayNode(input[i:])
			i = i + size - 1
			node.value = arrayNode
			break
		} else if c == '"' {
			// Beginning of a string
			str, size := createStringNode(input[i:])
			i = i + size - 1
			node.value = str
			break
		} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			// Whitespace, ignore it
		} else {
			buffer = buffer + string(c)
		}
	}

	//	if buffer != "" {
	//		node.value = buffer
	//	}

	//	fmt.Printf("KEYVALUE : %v %v\n", input[:i+1], i+1)
	return node, i + 1
}

func (node *gdbResultNode) toJSON() string {
	buffer := ""

	if node.nodeType == "array" {
		buffer = buffer + "["

		for idx, child := range node.children {
			if idx > 0 {
				buffer = buffer + ","
			}

			buffer = buffer + child.toJSON()
		}

		buffer = buffer + "]"
	} else if node.nodeType == "object" {
		buffer = buffer + "{"

		for idx, child := range node.children {
			if idx > 0 {
				buffer = buffer + ","
			}

			buffer = buffer + child.toJSON()
		}

		buffer = buffer + "}"
	} else if node.nodeType == "keyvalue" {
		key := node.key
		if key != "" {
			buffer = buffer + "\"" + key + "\":"
		}

		stringValue, ok := node.value.(string)
		if ok {
			buffer = buffer + stringValue
		} else {
			valueNode, ok := node.value.(gdbResultNode)
			if ok {
				buffer = buffer + valueNode.toJSON()
			}
		}
	}

	return buffer
}
