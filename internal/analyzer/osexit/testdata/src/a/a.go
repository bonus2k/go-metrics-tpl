package a

import (
	"fmt"
	"os"
)

func main() {
	var i = 2
	x := sum(i)
	fmt.Println(x)
}

func sum(i int) int {
	x := i + 2
	if x < 0 {
		os.Exit(1)
	}
	return x
}
