package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const Display_bytecode_mode bool = true
const Display_stack bool = false

func main() {
	dir, _ := os.Getwd()

	options := []string{}
	fOptions := []string{}

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".py" {
				options = append(options, path)
			}
			return nil
		})

	if err != nil {
		log.Println(err)
	}

	for i, option := range options {
		indexStr := strconv.Itoa(i + 1)
		fOptions = append(fOptions, "File "+indexStr+": "+option)
	}

	fmt.Println("Select a Python file to interpret:")
	for _, option := range fOptions {
		fmt.Println(option)
	}

	fmt.Print("Enter your choice (1, 2, 3): ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(options) {
		fmt.Println("Invalid choice. Please select a valid option.")
		return
	}

	da_option := options[choice-1]

	cmd := exec.Command("python3", "toBytecode.py", da_option)

	output, _ := cmd.CombinedOutput()
	if Display_bytecode_mode {
		fmt.Println(string(output))
	}

	instructions, _ := ParseBytecode(string(output))

	da_vm := makeDaVM(instructions)

	da_vm.run()
}
