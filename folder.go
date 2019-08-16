package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Folder struct {
	innerHandler

	source  string
	target  string
	ignores map[string]bool
}

func (s *Folder) ParseArg(key, value string) {
	if key == "help" {
		s.help = true
	} else if key == "create" {
		s.command = "create"
	} else if key == "delete" {
		s.command = "delete"
	} else if key == "copy" {
		s.command = "copy"
	} else if key == "clear" {
		s.command = "clear"
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
	} else if key == "ignore" {
		items := strings.Split(value, "|")
		count := len(items)
		for i := 0; i < count; i++ {
			item := strings.ToLower(items[i])
			if len(item) > 0 {
				s.ignores[item] = true
			}
		}
	}
}

func (s *Folder) Handle() error {
	if s.help {
		s.ShowHelp()
		return nil
	}

	if s.command == "create" {
		return s.Create()
	} else if s.command == "delete" {
		return s.Delete()
	} else if s.command == "clear" {
		return s.Clear()
	} else if s.command == "copy" {
		return s.Copy()
	}

	return fmt.Errorf("invalid command: %s", s.command)
}

func (s *Folder) ShowHelp() {
	s.ShowLine("Usage:", "", 0)
	s.ShowLine("    cmdx folder <command> [arguments]", "", 0)
	s.ShowLine("The command are:", "", 0)
	labelWidth := 12
	s.ShowLine("    help", "show the command list", labelWidth)
	s.ShowLine("    create", "create folder, <target=folder path>", labelWidth)
	s.ShowLine("    delete", "delete folder, <target=folder path>", labelWidth)
	s.ShowLine("    clear", "clear folder, <target=folder path> [ignore=node_modules|.git]", labelWidth)
	s.ShowLine("    copy", "copy folder, <source=folder path> <target=folder path> [ignore=folder1|file1.txt|.git]", labelWidth)
}

func (s *Folder) Create() error {
	if len(s.target) < 1 {
		return fmt.Errorf("invalid arguments: target is empty")
	}

	return os.MkdirAll(s.target, 0700)
}

func (s *Folder) Delete() error {
	if len(s.target) < 1 {
		return fmt.Errorf("invalid arguments: target is empty")
	}
	info, err := os.Stat(s.target)
	if os.IsNotExist(err) {
	} else {
		if !info.IsDir() {
			return fmt.Errorf("invalid arguments: target '%s' is not folder", s.target)
		} else {
			return os.RemoveAll(s.target)
		}
	}

	return nil
}

func (s *Folder) Clear() error {
	if len(s.target) < 1 {
		return fmt.Errorf("invalid arguments: target is empty")
	}
	info, err := os.Stat(s.target)
	if os.IsNotExist(err) {
	} else {
		if !info.IsDir() {
			return fmt.Errorf("invalid arguments: target '%s' is not folder", s.target)
		} else {
			items, err := ioutil.ReadDir(s.target)
			if err != nil {
				return err
			}
			for _, item := range items {
				if s.isIgnored(item.Name()) {
					continue
				}
				path := filepath.Join(s.target, item.Name())
				err = os.RemoveAll(path)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *Folder) Copy() error {
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
	if !info.IsDir() {
		return fmt.Errorf("invalid arguments: source '%s' is not folder", s.source)
	}
	if len(s.target) < 1 {
		return fmt.Errorf("invalid arguments: target is empty")
	}

	return s.copy(s.source, s.target, info)
}

func (s *Folder) copy(src, dest string, info os.FileInfo) error {
	if s.isIgnored(info.Name()) {
		return nil
	}

	if info.IsDir() {
		return s.copyDirectory(src, dest, info)
	}
	return s.copyFile(src, dest, info)
}

func (s *Folder) copyFile(src, dest string, info os.FileInfo) error {
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

func (s *Folder) copyDirectory(src, dest string, info os.FileInfo) error {
	if err := os.MkdirAll(dest, info.Mode()); err != nil {
		return err
	}

	infos, err := ioutil.ReadDir(src)
	if err != nil {
		return err
	}

	for _, info := range infos {
		err := s.copy(filepath.Join(src, info.Name()), filepath.Join(dest, info.Name()), info)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Folder) isIgnored(name string) bool {
	if len(name) > 0 {
		_, ok := s.ignores[strings.ToLower(name)]
		if ok {
			return true
		}
	}

	return false
}
