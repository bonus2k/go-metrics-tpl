package b

import (
	"fmt"
	"os"
)

func main() {
	var i = 2
	x := sum(i)

	if x < 5 {
		os.Exit(1) // want "it is recommended to avoid calling the os.Exit function from the main function"
	}

	switch true {
	case true:
		os.Exit(0) // want "it is recommended to avoid calling the os.Exit function from the main function"
	}

	fmt.Println(x)
}

func sum(i int) int {
	x := i + 2
	return x
}
