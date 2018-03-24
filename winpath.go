package main

import (
	"fmt"
	"log"
	"os"
	s "strings"

	"golang.org/x/sys/windows/registry"
)

func main() {
	const PathEnvVariableKey = "FOO"

	if len(os.Args) != 2 {
		fmt.Println(`Usage: winpath "<your\new\path>"`)
		return
	}

	newEnvVal := os.Args[1]

	k, err := registry.OpenKey(registry.CURRENT_USER, "Environment", registry.ALL_ACCESS)
	if err != nil {
		log.Fatal(err)
	}
	defer k.Close()

	envVal, _, err := k.GetStringValue(PathEnvVariableKey)
	if err != nil {
		envVal = ""
	}

	fmt.Println("Current Value: ", envVal)
	var envVals []string
	if len(envVal) > 0 {
		envVals = s.Split(envVal, ";")
		envVals = append(envVals, newEnvVal)
		envVals = removeDuplicates(envVals)
		envVal = s.Join(envVals, ";")
	} else {
		envVal += newEnvVal
	}

	err = k.SetExpandStringValue(PathEnvVariableKey, envVal)
	if err != nil {
		log.Fatal(err)
	}

	envVal, _, err = k.GetStringValue(PathEnvVariableKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done: ", envVal)
}

func removeDuplicates(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
