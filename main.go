package main

import (
	"github.com/aerokube/images/cmd"
)

//go:generate pkger -include /static -include /selenium/base -o build

func main() {
	cmd.Execute()
}
