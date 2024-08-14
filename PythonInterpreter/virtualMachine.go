package main

import (
	"fmt"
	"strconv"
)

func (vm *VirtualMachine) popStack() string {
	topOfDaStack := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return topOfDaStack
}

type VirtualMachine struct {
	stack     []string
	bytecodes []BytecodeInstruction
	namespace map[string]any
}

func (vm *VirtualMachine) handlePOP_TOP() {
	vm.stack = vm.stack[:len(vm.stack)-1]
}

// Figure out name vs const

func (vm *VirtualMachine) handleLOAD_NAME(b BytecodeInstruction) {
	vm.stack = append(vm.stack, b.Argument)
}

func (vm *VirtualMachine) handleLOAD_CONST(b BytecodeInstruction) {
	vm.stack = append(vm.stack, b.Argument)
}

func (vm *VirtualMachine) handleBINARY_ADD() {
	arg1 := vm.popStack()
	arg2 := vm.popStack()

	intArg1, _ := strconv.Atoi(arg1)
	intArg2, _ := strconv.Atoi(arg2)

	sum := intArg1 + intArg2

	strSum := strconv.Itoa(sum)

	vm.stack = append(vm.stack, strSum)
}

func (vm *VirtualMachine) handleRETURN_VALUE() {

	if vm.stack[len(vm.stack)-1] == "None" {
		vm.stack = vm.stack[:len(vm.stack)-1]
		return
	}
}

func (vm *VirtualMachine) handleCALL_FUNCTION(b BytecodeInstruction) {
	num_args, _ := strconv.Atoi(b.ArgIndex)
	args := make([]interface{}, num_args)

	for i := range num_args {
		args[i] = vm.popStack()
	}
	function := vm.popStack()

	if function == "print" && len(args) == 1 {
		fmt.Println(args[0])
		vm.stack = append(vm.stack, "None")
	}
}

func (vm *VirtualMachine) run() {
	for _, instruction := range vm.bytecodes {
		switch opcode := instruction.Opcode; opcode {
		case "LOAD_NAME":
			vm.handleLOAD_NAME(instruction)
		case "LOAD_CONST":
			vm.handleLOAD_CONST(instruction)
		case "POP_TOP":
			vm.handlePOP_TOP()
		case "CALL_FUNCTION":
			vm.handleCALL_FUNCTION(instruction)
		case "RETURN_VALUE":
			vm.handleRETURN_VALUE()
		default:
			fmt.Println("WHOOPS")
		}
	}
}

func makeDaVM(b []BytecodeInstruction) *VirtualMachine {
	var namespace = map[string]any{
		"print": fmt.Print,
	}
	return &VirtualMachine{
		namespace: namespace,
		bytecodes: b,
	}
}
