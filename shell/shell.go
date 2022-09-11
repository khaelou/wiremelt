package shell

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"wiremelt/worker"
)

func InitShell(spec worker.MacroSpec) {
	reader := bufio.NewReader(os.Stdin)

	for {
		hostName, err := os.Hostname() // Get Hostname
		if err != nil {
			log.Println(err)
		}

		wireMeltShell := fmt.Sprintf(">_ ] wiremelt@%s\n", hostName)
		fmt.Println(wireMeltShell)

		// Read the keyboad input.
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		if strings.Contains(input, "*force") { // *force used to execute command on underlying system
			// Prevent self-execution of software within shell
			if !strings.Contains(input, "go run") && !strings.Contains(input, "wiremelt") {
				cmdOpr := strings.Replace(input, "*force", "", -1) // Remove operator for cmd execution
				input = strings.TrimSpace(cmdOpr)

				// Handle the execution of the input
				if err = execSysInput(input); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		} else {
			// Handle the execution of the input
			input, err := execInput(input)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			inputHandler(input, spec)
		}
	}
}

func execInput(input string) (string, error) {
	if input == "" {
		return "", errors.New("no input")
	}

	input = strings.TrimSuffix(input, "\n") // Remove the newline character.

	return input, nil
}

func execSysInput(input string) error {
	if input == "" {
		return errors.New("no input")
	}

	// Remove the newline character.
	input = strings.TrimSuffix(input, "\n")

	// Split the input to separate the command and the arguments.
	args := strings.Split(input, " ")

	// Check for default commands.
	switch args[0] {
	case "cd":
		// 'cd' to home dir with empty path not yet supported.
		if len(args) < 2 {
			return errors.New("path required")
		}
		// Change the directory and return the error.
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	case "quit":
		os.Exit(0)
	}

	// Pass the program and the arguments separately.
	cmd := exec.Command(args[0], args[1:]...)

	// Set the correct output device.
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// Execute the command and return the error.
	return cmd.Run()
}

func inputHandler(input string, macroLibrary worker.MacroSpec) {
	fmt.Println(">", input)
	lowerInput := strings.ToLower(input)

	switch lowerInput {
	case "macros":
		fmt.Println()
		fmt.Println(macroLibrary)
	case "exit":
		os.Exit(0)
	case "quit":
		os.Exit(0)
	}
}
