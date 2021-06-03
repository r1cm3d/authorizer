package main

import (
	"bufio"
	"fmt"
	"github.com/r1cm3d/authorizer/internal"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	timeline := internal.NewTimeline()
	fmt.Println()
	for scanner.Scan() {
		event := internal.Parse(scanner.Text())
		timeline.Process(event)
		fmt.Println(timeline.Last())
	}
	fmt.Println()
}
