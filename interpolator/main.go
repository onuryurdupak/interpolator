package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

const (
	stamp_build_date  = "${build_date}"
	stamp_commit_hash = "${commit_hash}"

	errInput    = 1
	errInternal = 2
	errUnkown   = 3
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 && (args[0] == "version" || args[0] == "--version") {
		fmt.Printf("Build Date: %s | Commit: %s\n", stamp_build_date, stamp_commit_hash)
		os.Exit(0)
	}

	if len(args) < 3 {
		fmt.Println("Define a file path, a seperator and at least one key=value pair.")
		os.Exit(errInput)
	}

	_, err := os.Stat(args[0])
	if err != nil {
		fmt.Printf("Error checking file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(errUnkown)
	}

	if os.IsNotExist(err) {
		fmt.Printf("File does not exist: '%s'.\n", args[0])
		os.Exit(errInput)
	}

	fileContent, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Printf("Error reading file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(errUnkown)
	}
	fileContentStr := string(fileContent)
	separator := args[1]

	replaces := make(map[string]string, len(args)-3)

	for i := 2; i < len(args); i++ {
		split := strings.Split(args[i], separator)

		if len(split) != 2 {
			fmt.Printf("Invalid key=value assignment: '%s'.\n", args[i])
			os.Exit(errInput)
		}

		_, ok := replaces[split[0]]
		if ok {
			fmt.Printf("Key: '%s' is defined multiple times in arguments.\n", split[0])
			os.Exit(errInput)
		}

		if split[0] == "" {
			fmt.Printf("A key can not be an empty string.\n")
			os.Exit(errInput)
		}

		regKeyMatch, err := regexp.Compile(split[0])

		if err != nil {
			fmt.Printf("Key '%s' is not regex compliant. Error: '%s'.\n", split[0], err.Error())
			os.Exit(errInput)
		}

		matches := regKeyMatch.FindAllString(fileContentStr, -1)
		if len(matches) == 0 {
			fmt.Printf("Key: '%s' not found in file.\n", split[0])
			os.Exit(errInput)
		} else if len(matches) > 1 {
			fmt.Printf("Key: '%s' is defined %v times in the file: '%s'.\n", split[0], len(matches), args[0])
			os.Exit(errInput)
		}

		replaces[matches[0]] = split[1]
	}

	for k, v := range replaces {
		fileContentStr = strings.Replace(fileContentStr, k, v, -1)
	}

	const maxTries = 10
	currentTries := maxTries
	tempFilename := ""

	for {
		if currentTries == 0 {
			fmt.Printf("Could not create a unique temp file name after %d tries.\n", maxTries)
			os.Exit(errInternal)
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

	err = ioutil.WriteFile(tempFilename, []byte(fileContentStr), 0644)
	if err != nil {
		fmt.Printf("Error writing file at: '%s' error: '%s.\n", args[0], err.Error())
		os.Exit(errInternal)
	}

	err = os.Remove(args[0])
	if err != nil {
		fmt.Printf("Error clearing file at: '%s' error: '%s.\n", args[0], err.Error())
		os.Exit(errInternal)
	}

	err = os.Rename(tempFilename, args[0])
	if err != nil {
		fmt.Printf("Could not restore original file name: '%s.\n", err.Error())
		os.Exit(errInternal)
	}
}
