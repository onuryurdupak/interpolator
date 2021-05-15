package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) < 2 {
		fmt.Println("Define a file path and at least one key=value pair.")
		os.Exit(1)
	}

	_, err := os.Stat(args[0])

	if err != nil {
		fmt.Printf("Error checking file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(2)
	}

	if os.IsNotExist(err) {
		fmt.Printf("File does not exist: '%s'.\n", args[0])
		os.Exit(3)
	}

	fileContent, err := ioutil.ReadFile(args[0])
	fileContentStr := string(fileContent)

	if err != nil {
		fmt.Printf("Error reading file at: '%s' error: '%s'.\n", args[0], err.Error())
		os.Exit(4)
	}

	replaces := map[string]string{}

	for i := 1; i < len(args); i++ {
		split := strings.Split(args[i], "=")

		if len(split) != 2 {
			fmt.Printf("Invalid key=value assignment: '%s'.\n", args[i])
			os.Exit(5)
		}

		_, ok := replaces[split[0]]

		if ok {
			fmt.Printf("Key: '%s' is defined multiple times in arguments.\n", split[0])
			os.Exit(5)
		}

		if split[0] == "" {
			fmt.Printf("A key can not be an empty string.\n")
			os.Exit(6)
		}

		regKeyMatch, err := regexp.Compile(split[0])

		if err != nil {
			fmt.Printf("Key '%s' is not regex compliant. Error: '%s'.\n",
				split[0], err.Error())
			os.Exit(7)
		}

		matches := regKeyMatch.FindAllStringIndex(fileContentStr, -1)

		if len(matches) == 0 {
			fmt.Printf("Key: '%s' not found in file.\n", split[0])
			os.Exit(8)
		} else if len(matches) > 1 {
			fmt.Printf("Key: '%s' is defined %v times in the file: '%s'.\n",
				split[0], len(matches), args[0])
			os.Exit(9)
		}

		replaces[split[0]] = split[1]
	}

	for k, v := range replaces {
		fileContentStr = strings.Replace(fileContentStr, k, v, -1)
	}

	err = os.Remove(args[0])

	if err != nil {
		fmt.Printf("Error clearing file at: '%s' error: '%s.\n", args[0], err.Error())
		os.Exit(10)
	}

	err = ioutil.WriteFile(args[0], []byte(fileContentStr), 0644)

	if err != nil {
		fmt.Printf("Error writing file at: '%s' error: '%s.\n", args[0], err.Error())
		os.Exit(11)
	}

	fmt.Println("Success.")
}
