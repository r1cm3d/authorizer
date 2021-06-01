package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	i := 0
	for scanner.Scan() {
		//TODO: call parser here
		fmt.Printf("\nLine: %d\n%s", i, scanner.Text())
		i++
	}

	fmt.Printf("\n\n\nEnd of file")
}
