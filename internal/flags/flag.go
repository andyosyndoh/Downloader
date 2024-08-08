package flags

import (
	"strings"
)

func OutputFileFlag(args []string) bool {
	return strings.HasPrefix(args[0], "-O=")
}

func OutputAndPath(args []string) bool {
	return (strings.HasPrefix(args[0], "-O=") && strings.HasPrefix(args[1], "-P=")) ||
		(strings.HasPrefix(args[1], "-O=") && strings.HasPrefix(args[2], "-P="))
}

func GetFlagInput(flagInput string) string {
	return flagInput[3:]
}

func FlagType(args []string) []string {
	var flagtype []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-O=") {
			flagtype = append(flagtype, "-O=")
		} else if strings.HasPrefix(arg, "-P=") {
			flagtype = append(flagtype, "-P=")
		} else if strings.HasPrefix(arg, "--rate-limit=") {
			flagtype = append(flagtype, "--rate-limit=")
		}
	}
	return flagtype
}
