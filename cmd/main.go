package main

import (
	"os"

	"github.com/valerii-smirnov/petli-test-task/application"
)

func main() {
	if err := application.New().Init().Run(os.Args); err != nil {
		panic(err)
	}
}
