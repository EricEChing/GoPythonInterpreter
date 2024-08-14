package main

import (
	"fmt"
	"strconv"
	"strings"
)

func GetArg(s string) string {
	start := 0

	startIndex := strings.Index(s, "(")
	if startIndex == -1 {
		return ""
	}

	startIndex += start
	endIndex := strings.Index(s, ")")
	if endIndex == -1 {
		return ""
	}
	return s[startIndex+1 : endIndex]
}

func IsInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// BytecodeInstruction represents a single bytecode instruction
type BytecodeInstruction struct {
	LineNumber string
	Opcode     string
	ArgIndex   string
	Argument   string
}

// ParseBytecode parses a bytecode string into a list of BytecodeInstructions
func ParseBytecode(bytecode string) ([]BytecodeInstruction, error) {
	// Assumes first line with 2 ints has () arg
	var instructions []BytecodeInstruction
	lines := strings.Split(bytecode, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var instruction BytecodeInstruction
		elements := strings.Fields(line)

		switch len(elements) {
		case 2:
			instruction.LineNumber, instruction.Opcode = elements[0], elements[1]
		case 3:
			instruction.LineNumber, instruction.Opcode, instruction.ArgIndex = elements[0], elements[1], elements[2]
		default:
			arg := GetArg(line)

			if string(arg[0]) == "'" && string(arg[len(arg)-1]) == "'" {
				arg = arg[1 : len(arg)-1]
			}

			if arg != "" {
				instruction.Argument = arg
			}

			if (IsInt(elements[0]) && IsInt(elements[1])) || elements[0] == ">>" {
				instruction.LineNumber, instruction.Opcode, instruction.ArgIndex = elements[1], elements[2], elements[3]
			} else {
				instruction.LineNumber, instruction.Opcode, instruction.ArgIndex = elements[0], elements[1], elements[2]
			}

		}

		instructions = append(instructions, instruction)
	}
	return instructions, nil
}

// PrintInstructions prints the decomposed bytecode instructions
func PrintInstructions(instructions []BytecodeInstruction) {
	for _, inst := range instructions {
		fmt.Printf("Line %s: Opcode %s, Argument %s\n", inst.LineNumber, inst.Opcode, inst.Argument)
	}
}
