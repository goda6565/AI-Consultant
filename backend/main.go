package main

import (
	"fmt"
	"os"

	"github.com/goda6565/ai-consultant/backend/cmd"
)

func main() {
	cmd := cmd.NewRootCommand()
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
