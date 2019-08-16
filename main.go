package main

import (
	"fmt"
	"os"
)

func main() {
	args := &Args{}
	args.Parse(os.Args)
	handler := args.Handler()
	err := handler.Handle()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
