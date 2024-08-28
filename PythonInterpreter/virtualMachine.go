package main

import (
	"fmt"
	"strconv"
	"strings"
)

type VirtualMachine struct {
	stack       []string
	bytecodes   []BytecodeInstruction
	namespace   map[string]any
	currentLine int
}

func (vm *VirtualMachine) add2Stack(s string) {
	vm.stack = append(vm.stack, s)
}

func (vm *VirtualMachine) popStack() string {
	topOfDaStack := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	return topOfDaStack
}

func (vm *VirtualMachine) handleCOMPARE_OP(b BytecodeInstruction) {
	arg1 := vm.popStack()
	arg2 := vm.popStack()

	if b.Argument == "==" {
		if arg1 == arg2 {
			vm.add2Stack("true")
		} else {
			vm.add2Stack("false")
		}
	} else if b.Argument == "!=" {
		if arg1 != arg2 {
			vm.add2Stack("true")
		} else {
			vm.add2Stack("false")
		}
	}

}

func (vm *VirtualMachine) handleBUILD_LIST() {
	// assumes BUILD_LIST has no args
	vm.add2Stack("[]")
}

func (vm *VirtualMachine) handleLIST_EXTEND(b BytecodeInstruction) {
	num_args, _ := strconv.Atoi(b.ArgIndex)
	args := make([]string, num_args)
	for i := 0; i < num_args; i++ {
		arg := vm.popStack()
		arg = stripPara(arg)
		args[i] = arg
	}
	list := vm.popStack()

	protoList := ""

	if list == "[]" {
		protoList += "["

		for _, v := range args {
			protoList += v
			protoList += ", "
		}
		protoList = protoList[:len(protoList)-2]
		protoList += "]"
		vm.add2Stack(protoList)
	}
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

func (vm *VirtualMachine) handlePOP_JUMP_IF_FALSE(b BytecodeInstruction) {
	boo_lean := vm.popStack()
	if boo_lean == "false" {
		vm.currentLine, _ = strconv.Atoi(b.ArgIndex)
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

func (vm *VirtualMachine) handleGET_ITER() {
	iterStr := vm.popStack()
	
}

func (vm *VirtualMachine) handleJUMP_FORWARD(b BytecodeInstruction) {
	elems := strings.Split(b.Argument, " ")
	line := elems[1]
	lineInt, _ := strconv.Atoi(line)
	vm.currentLine = lineInt
}

func (vm *VirtualMachine) run() {

	for i, instruction := range vm.bytecodes {
		if vm.currentLine != i*2 {
			continue
		} else {
			vm.currentLine += 2
		}
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
			vm.handleCOMPARE_OP(instruction)
		case "POP_JUMP_IF_FALSE":
			vm.handlePOP_JUMP_IF_FALSE(instruction)
		case "JUMP_FORWARD":
			vm.handleJUMP_FORWARD(instruction)
		case "BUILD_LIST":
			vm.handleBUILD_LIST()
		case "LIST_EXTEND":
			vm.handleLIST_EXTEND(instruction)
		default:
			fmt.Println("WHOOPS, " + opcode + " NOT RECOGNIZED")
		}
		if Display_stack {
			fmt.Println(vm.stack)
		}

	}
}

func makeDaVM(b []BytecodeInstruction) *VirtualMachine {
	var namespace = map[string]any{
		"print": "print",
		"range": "range",
	}
	return &VirtualMachine{
		namespace: namespace,
		bytecodes: b,
	}
}
