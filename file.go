package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type File struct {
	innerHandler

	source  string
	target  string
	content string
}

func (s *File) ParseArg(key, value string) {
	if key == "help" {
		s.help = true
	} else if key == "create" {
		s.command = "create"
	} else if key == "delete" {
		s.command = "delete"
	} else if key == "copy" {
		s.command = "copy"
	} else if key == "source" {
		if len(value) > 0 {
			if !filepath.IsAbs(value) {
				s.source = filepath.Join(s.folder, value)
			} else {
				s.source = value
			}
		}
	} else if key == "target" {
		if len(value) > 0 {
			if !filepath.IsAbs(value) {
				s.target = filepath.Join(s.folder, value)
			} else {
				s.target = value
			}
		}
	} else if key == "content" {
		s.content = value
	}
}

func (s *File) Handle() error {
	if s.help {
		s.ShowHelp()
		return nil
	}

	if s.command == "create" {
		return s.Create()
	} else if s.command == "delete" {
		return s.Delete()
	} else if s.command == "copy" {
		return s.Copy()
	}

	return fmt.Errorf("invalid command: %s", s.command)
}

func (s *File) ShowHelp() {
	s.ShowLine("Usage:", "", 0)
	s.ShowLine("    cmdx file <command> [arguments]", "", 0)
	s.ShowLine("The command are:", "", 0)
	labelWidth := 12
	s.ShowLine("    help", "show the command list", labelWidth)
	s.ShowLine("    create", "create file, <target=file path> [content=string value]", labelWidth)
	s.ShowLine("    delete", "delete file, <target=file path>", labelWidth)
	s.ShowLine("    copy", "copy file, <source=file path> <target=file path>", labelWidth)
}

func (s *File) Create() error {
	if len(s.target) < 1 {
		return fmt.Errorf("invalid arguments: target is empty")
	}

	filePath := s.target
	fileFolder := filepath.Dir(filePath)
	_, err := os.Stat(fileFolder)
	if os.IsNotExist(err) {
		os.MkdirAll(fileFolder, 0777)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if len(s.content) > 0 {
		_, err = fmt.Fprint(file, s.content)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *File) Delete() error {
	if len(s.target) < 1 {
		return fmt.Errorf("invalid arguments: target is empty")
	}
	info, err := os.Stat(s.target)
	if os.IsNotExist(err) {
	} else {
		if info.IsDir() {
			return fmt.Errorf("invalid arguments: target '%s' is not file", s.target)
		} else {
			return os.Remove(s.target)
		}
	}

	return nil
}

func (s *File) Copy() error {
	if len(s.source) < 1 {
		return fmt.Errorf("invalid arguments: source is empty")
	}
	info, err := os.Stat(s.source)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("invalid arguments: source '%s' is not exist", s.source)
		}
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("invalid arguments: source '%s' is not file", s.source)
	}
	if len(s.target) < 1 {
		return fmt.Errorf("invalid arguments: target is empty")
	}

	return s.copyFile(s.source, s.target, info)
}

func (s *File) copyFile(src, dest string, info os.FileInfo) error {
	folderPath := filepath.Dir(dest)
	_, err := os.Stat(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(folderPath, 0777)
		}
	}

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if err = os.Chmod(destFile.Name(), info.Mode()); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}
