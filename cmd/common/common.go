package common

import (
	"fmt"
	"os"
)

// Get the command line argument by index.
func GetCmdArg(index int, userArgs []string) string {

	// Access arguments excluding the program name
	if len(os.Args[1:]) < 1 {
		fmt.Println("Missing command line args")
		return ""
	}

	//userArgs := os.Args[1:]
	//fmt.Println("User-provided arguments:", userArgs)

	// Access individual arguments by index
	//if len(os.Args) > 1 {
	//	fmt.Println("First user argument:", os.Args[1])
	//}

	// Iterate through user-provided arguments
	fmt.Println("Iterating through user arguments:")
	for i, arg := range userArgs {
		//fmt.Printf("Argument %d: %s\n", i+1, arg)
		if i == index {
			return arg
		}
	}
	return ""
}
