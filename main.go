package main

import (
	"fmt"

	generator "github.com/andrewinci/rap/fieldgen"
)

func main() {
	x := generator.NewGenerator("A-Z")
	fmt.Println(x)
}
