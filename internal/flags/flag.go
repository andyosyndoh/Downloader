package flags

import (
	"strings"
)

func OutputFileFlag(args []string) bool {
	return strings.HasPrefix(args[0], "-O=")
}

func OutputAndPath(args []string) bool {

	if (strings.HasPrefix(args[0], "-O=") &&
		strings.HasPrefix(args[1], "-P=")) ||
		(strings.HasPrefix(args[1], "-O=") &&
			strings.HasPrefix(args[2], "-P=")) {
		return true
	}
	return false
}
func GetFlagInput(flagInput string) string {

	var input string

	switch flagInput[:3] {
	case "-O=":
		input = flagInput[3:]
	case "-P=":
		input = flagInput[3:]
	}
	return input
}
