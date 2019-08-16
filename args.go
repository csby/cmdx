package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Args struct {
	innerHandler

	handler Handler
}

func (s *Args) Parse(args []string) {
	argsLength := len(args)
	if argsLength < 1 {
		return
	}
	path, err := filepath.Abs(args[0])
	if err != nil {
		return
	}
	s.folder = filepath.Dir(path)

	if argsLength < 2 {
		s.handler = s
		return
	}

	s.handler = s.generateHandler(s.folder, strings.ToLower(args[1]))
	if s.handler == nil {
		return
	}

	for argsIndex := 2; argsIndex < argsLength; argsIndex++ {
		arg := args[argsIndex]
		splitIndex := strings.Index(arg, "=")
		if splitIndex < 1 {
			k := strings.ToLower(arg)
			if k == "help" || k == "-help" || k == "-h" {
				s.handler.ParseArg("help", "")
			} else {
				s.handler.ParseArg(k, "")
			}
			continue
		}

		if splitIndex >= len(arg)-1 {
			continue
		}
		key := strings.ToLower(arg[0:splitIndex])
		val := strings.Trim(arg[splitIndex+1:], "\"")
		s.handler.ParseArg(key, val)
	}
}

func (s *Args) Handler() Handler {
	if s.handler != nil {
		return s.handler
	}

	return s
}

func (s *Args) ParseArg(key, value string) {

}

func (s *Args) Handle() error {
	s.ShowLine("cmdx is a tool for command line extension.", "", 0)
	s.ShowLine("Usage:", "", 0)
	s.ShowLine("    cmdx <command> [arguments]", "", 0)
	s.ShowLine("The command are:", "", 0)
	labelWidth := 12
	s.ShowLine("    help", "show the command list", labelWidth)
	s.ShowLine("    folder", "create, delete, clear, copy folder", labelWidth)

	if s.help {
		return nil
	} else {
		return fmt.Errorf("invalid commad")
	}
}

func (s *Args) generateHandler(folder, name string) Handler {
	if name == "h" || name == "help" {
		s.help = true
		return s
	}

	if name == "folder" {
		return &Folder{innerHandler: innerHandler{folder: folder}, ignores: make(map[string]bool)}
	}

	return s
}
