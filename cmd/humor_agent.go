package main

import (
	"github.com/wilenceyao/humor-agent/internal"
	"os"
)

func main() {
	a := &internal.HumorAgent{}
	err := a.Start()
	if err != nil {
		os.Exit(1)
	}
}
