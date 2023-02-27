package program

import (
	"fmt"
	"github.com/onuryurdupak/gomod/v2/core"
	"github.com/onuryurdupak/gomod/v2/slice"

	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func Main(args []string) {
	if len(args) == 1 {
		switch {
		case args[0] == "version" || args[0] == "--version" || args[0] == "-v":
			fmt.Println(versionInfo())
			os.Exit(ErrSuccess)
			return
		case args[0] == "help" || args[0] == "--help" || args[0] == "-h":
			fmt.Println(helpMessageStyled())
			os.Exit(ErrSuccess)
		case args[0] == "-hr":
			fmt.Println(helpMessageUnstyled())
			os.Exit(ErrSuccess)
		default:
			fmt.Println(helpPrompt)
			os.Exit(ErrInput)
		}
	}

	recursiveMode := slice.RemoveString(&args, "-r")
	replaceCount := core.ConditionalInt(recursiveMode, -1, 1)

	if len(args) < 3 {
		fmt.Println("Define a file path, a separator and at least one key=value pair.")
		fmt.Println(helpPrompt)
		os.Exit(ErrInput)
	}

	_, err := os.Stat(args[0])
	if err != nil {
		fmt.Printf("Error checking file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(ErrUnknown)
	}

	if os.IsNotExist(err) {
		fmt.Printf("File does not exist: '%s'.\n", args[0])
		os.Exit(ErrInput)
	}

	fileContent, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Printf("Error reading file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(ErrUnknown)
	}
	fileContentStr := string(fileContent)
	separator := args[1]

	replaces := make(map[string]string, len(args)-3)

	for i := 2; i < len(args); i++ {
		split := strings.Split(args[i], separator)

		if len(split) != 2 {
			fmt.Printf("Invalid key=value assignment: '%s'.\n", args[i])
			os.Exit(ErrInput)
		}

		_, ok := replaces[split[0]]
		if ok {
			fmt.Printf("Key: '%s' is defined multiple times in arguments.\n", split[0])
			os.Exit(ErrInput)
		}

		if split[0] == "" {
			fmt.Printf("A key can not be an empty string.\n")
			os.Exit(ErrInput)
		}

		regKeyMatch, err := regexp.Compile(split[0])

		if err != nil {
			fmt.Printf("Key '%s' is not regex compliant. Error: '%s'.\n", split[0], err.Error())
			os.Exit(ErrInput)
		}

		matches := regKeyMatch.FindAllString(fileContentStr, -1)
		if len(matches) == 0 {
			fmt.Printf("Key: '%s' not found in file.\n", split[0])
			os.Exit(ErrInput)
		} else if len(matches) > 1 && !recursiveMode {
			fmt.Printf("Key: '%s' is defined %v times in the file: '%s'.\nRun with -r for recursive mode.\n", split[0], len(matches), args[0])
			os.Exit(ErrInput)
		}

		replaces[matches[0]] = split[1]
	}

	for k, v := range replaces {
		fileContentStr = strings.Replace(fileContentStr, k, v, replaceCount)
	}

	const maxTries = 10
	currentTries := maxTries
	tempFilename := ""

	for {
		if currentTries == 0 {
			fmt.Printf("Could not create a unique temp file name after %d tries.\n", maxTries)
			os.Exit(ErrInternal)
		}

		guid := uuid.New().String()
		tempFilename = fmt.Sprintf("%s_%s", args[0], guid)

		_, err := os.Stat(tempFilename)
		if err != nil {
			if os.IsNotExist(err) {
				break
			}
		}

		currentTries--
	}

	err = os.WriteFile(tempFilename, []byte(fileContentStr), 0644)
	if err != nil {
		fmt.Printf("Error writing file at: '%s' error: '%s.\n", args[0], err.Error())
		os.Exit(ErrInternal)
	}

	err = os.Remove(args[0])
	if err != nil {
		fmt.Printf("Error clearing file at: '%s' error: '%s.\n", args[0], err.Error())
		os.Exit(ErrInternal)
	}

	err = os.Rename(tempFilename, args[0])
	if err != nil {
		fmt.Printf("Could not restore original file name: '%s.\n", err.Error())
		os.Exit(ErrInternal)
	}
}
