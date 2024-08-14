package main

import (
	"fmt"
	"strconv"
)

type VirtualMachine struct {
	stack     []string
	bytecodes []BytecodeInstruction
	namespace map[string]any
}

func (vm *VirtualMachine) add2Stack(s string) {
	vm.stack = append(vm.stack, s)
}

func (vm *VirtualMachine) popStack() string {
	topOfDaStack := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return topOfDaStack
}

func (vm *VirtualMachine) handleCOMPARE_OP() {
	arg1 := vm.popStack()
	arg2 := vm.popStack()

	if arg1 == arg2 {
		vm.add2Stack("true")
	}
	vm.add2Stack("false")

}

func (vm *VirtualMachine) handlePOP_TOP() {
	vm.stack = vm.stack[:len(vm.stack)-1]
}

func (vm *VirtualMachine) handleSTORE_NAME(b BytecodeInstruction) {
	thingy := vm.popStack()
	vm.namespace[b.Argument] = thingy
}

func (vm *VirtualMachine) handleLOAD_NAME(b BytecodeInstruction) {
	vm.add2Stack(vm.namespace[b.Argument].(string))
}

func (vm *VirtualMachine) handleLOAD_CONST(b BytecodeInstruction) {
	vm.add2Stack(b.Argument)
}

func (vm *VirtualMachine) handleBINARY_ADD() {
	arg1 := vm.popStack()
	arg2 := vm.popStack()

	intArg1, _ := strconv.Atoi(arg1)
	intArg2, _ := strconv.Atoi(arg2)

	sum := intArg1 + intArg2

	strSum := strconv.Itoa(sum)

	vm.add2Stack(strSum)
}

func (vm *VirtualMachine) handleBINARY_SUBTRACT() {
	arg1 := vm.popStack()
	arg2 := vm.popStack()

	intArg1, _ := strconv.Atoi(arg1)
	intArg2, _ := strconv.Atoi(arg2)

	difference := intArg2 - intArg1

	strDifference := strconv.Itoa(difference)

	vm.add2Stack(strDifference)
}

func (vm *VirtualMachine) handleBINARY_MULTIPLY() {
	arg1 := vm.popStack()
	arg2 := vm.popStack()

	intArg1, _ := strconv.Atoi(arg1)
	intArg2, _ := strconv.Atoi(arg2)

	product := intArg2 * intArg1

	strProduct := strconv.Itoa(product)

	vm.add2Stack(strProduct)
}

func (vm *VirtualMachine) handleBINARY_TRUE_DIVIDE() {
	arg1 := vm.popStack()
	arg2 := vm.popStack()

	intArg1, _ := strconv.Atoi(arg1)
	intArg2, _ := strconv.Atoi(arg2)

	dividend := float64(intArg2) / float64(intArg1)

	strDividend := strconv.FormatFloat(dividend, 'f', -1, 64)

	vm.add2Stack(strDividend)
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
		vm.add2Stack("None")
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
		case "BINARY_ADD":
			vm.handleBINARY_ADD()
		case "BINARY_SUBTRACT":
			vm.handleBINARY_SUBTRACT()
		case "BINARY_MULTIPLY":
			vm.handleBINARY_MULTIPLY()
		case "BINARY_TRUE_DIVIDE":
			vm.handleBINARY_TRUE_DIVIDE()
		case "STORE_NAME":
			vm.handleSTORE_NAME(instruction)
		case "COMPARE_OP":
			vm.handleCOMPARE_OP()
		default:
			fmt.Println("WHOOPS, " + opcode + " NOT RECOGNIZED")
		}
	}
}

func makeDaVM(b []BytecodeInstruction) *VirtualMachine {
	var namespace = map[string]any{
		"print": "print",
	}
	return &VirtualMachine{
		namespace: namespace,
		bytecodes: b,
	}
}
