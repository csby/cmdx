package main

import "testing"

func TestFolder_isMatched(t *testing.T) {
	folder := &Folder{patterns: make([]string, 0)}
	result := folder.isMatched("x.dll")
	if result != true {
		t.Fatal(result)
	}
	result = folder.isMatched("x.xml")
	if result != true {
		t.Fatal(result)
	}

	folder.patterns = append(folder.patterns, "*.dll")
	result = folder.isMatched("x.dll")
	if result != true {
		t.Fatal(result)
	}
	result = folder.isMatched("x.xml")
	if result != false {
		t.Fatal(result)
	}

	folder.patterns = append(folder.patterns, "x.xml")
	result = folder.isMatched("x.xml")
	if result != true {
		t.Fatal(result)
	}

	result = folder.isMatched("abc-x.xmlvv")
	if result != false {
		t.Fatal(result)
	}

	folder.patterns = append(folder.patterns, "ab*.xmlvv")
	result = folder.isMatched("abc-x.xmlvv")
	if result != true {
		t.Fatal(result)
	}
}
