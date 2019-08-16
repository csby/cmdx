package main

import "fmt"

type Handler interface {
	ParseArg(key, value string)
	Handle() error
}

type innerHandler struct {
	folder  string
	command string
	help    bool
}

func (s *innerHandler) ShowLine(label, value string, labelWidth int) {
	format := fmt.Sprintf("%%-%ds %%s", labelWidth)
	fmt.Printf(format, label, value)
	fmt.Println("")
}
