package main

import (
	"fmt"
	"interpolator/embed"
	"interpolator/utils/slice_utils"
	"interpolator/utils/stdout_utils"
	"interpolator/utils/syntax_utils"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

func main() {
	args := os.Args[1:]
	if len(args) == 1 {
		switch {
		case args[0] == "version" || args[0] == "--version" || args[0] == "-v":
			fmt.Printf("Build Date: %s | Commit: %s\n", embed.Stamp_build_date, embed.Stamp_commit_hash)
			os.Exit(embed.ErrSuccess)
			return
		case args[0] == "help" || args[0] == "--help" || args[0] == "-h":
			fmt.Printf("%s\n", stdout_utils.ProcessStyle(embed.HelpMessage))
			os.Exit(embed.ErrSuccess)
		case args[0] == "-hr":
			fmt.Printf("%s\n", stdout_utils.RemoveStyle(embed.HelpMessage))
			os.Exit(embed.ErrSuccess)
		default:
			fmt.Println(embed.HelpPrompt)
			os.Exit(embed.ErrInput)
		}
	}

	recursiveMode := slice_utils.RemoveString(&args, "-r")
	replaceCount := syntax_utils.ConditionalInt(recursiveMode, -1, 1)

	if len(args) < 3 {
		fmt.Println("Define a file path, a separator and at least one key=value pair.")
		fmt.Println(embed.HelpPrompt)
		os.Exit(embed.ErrInput)
	}

	_, err := os.Stat(args[0])
	if err != nil {
		fmt.Printf("Error checking file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(embed.ErrUnkown)
	}

	if os.IsNotExist(err) {
		fmt.Printf("File does not exist: '%s'.\n", args[0])
		os.Exit(embed.ErrInput)
	}

	fileContent, err := ioutil.ReadFile(args[0])
	if err != nil {
		fmt.Printf("Error reading file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(embed.ErrUnkown)
	}
	fileContentStr := string(fileContent)
	separator := args[1]

	replaces := make(map[string]string, len(args)-3)

	for i := 2; i < len(args); i++ {
		split := strings.Split(args[i], separator)

		if len(split) != 2 {
			fmt.Printf("Invalid key=value assignment: '%s'.\n", args[i])
			os.Exit(embed.ErrInput)
		}

		_, ok := replaces[split[0]]
		if ok {
			fmt.Printf("Key: '%s' is defined multiple times in arguments.\n", split[0])
			os.Exit(embed.ErrInput)
		}

		if split[0] == "" {
			fmt.Printf("A key can not be an empty string.\n")
			os.Exit(embed.ErrInput)
		}

		regKeyMatch, err := regexp.Compile(split[0])

		if err != nil {
			fmt.Printf("Key '%s' is not regex compliant. Error: '%s'.\n", split[0], err.Error())
			os.Exit(embed.ErrInput)
		}

		matches := regKeyMatch.FindAllString(fileContentStr, -1)
		if len(matches) == 0 {
			fmt.Printf("Key: '%s' not found in file.\n", split[0])
			os.Exit(embed.ErrInput)
		} else if len(matches) > 1 && !recursiveMode {
			fmt.Printf("Key: '%s' is defined %v times in the file: '%s'.\nRun with -r for recursive mode.\n", split[0], len(matches), args[0])
			os.Exit(embed.ErrInput)
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
			os.Exit(embed.ErrInternal)
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
		os.Exit(embed.ErrInternal)
	}

	err = os.Remove(args[0])
	if err != nil {
		fmt.Printf("Error clearing file at: '%s' error: '%s.\n", args[0], err.Error())
		os.Exit(embed.ErrInternal)
	}

	err = os.Rename(tempFilename, args[0])
	if err != nil {
		fmt.Printf("Could not restore original file name: '%s.\n", err.Error())
		os.Exit(embed.ErrInternal)
	}
}
