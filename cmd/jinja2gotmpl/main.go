package main

import (
	"fmt"
	"os"

	"github.com/Kade-Powell/jinja2gotmpl/pkg/j2g"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: jinja2gotmpl <template-file>")
		os.Exit(1)
	}

	input, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	result, err := j2g.Transpile(string(input))
	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
