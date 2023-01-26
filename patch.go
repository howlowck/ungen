package main

import (
	"fmt"
	"os"
	"strings"
)

type PatchType int

const (
	PatchReplace = iota
	PatchDelete
	PatchInsert
)

func (t PatchType) String() string {
	return [...]string{"PatchReplace", "PatchDelete", "PatchInsert"}[t]
}

type FileSystemOp int

const (
	FileCreate = iota
	FileDelete
	FileRename
	DirectoryDelete
	DirectoryRename
)

type ContentPatch struct {
	PatchType     PatchType
	OldLineNumber int
	OldLineCount  int
	NewContent    []string
}

func (t FileSystemOp) String() string {
	return [...]string{"FileCreate", "FileDelete", "FileRename", "DirectoryDelete", "DirectoryRename"}[t]
}

type Patch struct {
	Content *ContentPatch
	File    *FilePatch
}

type FilePatch struct {
	FileOp     FileSystemOp
	SourcePath string
	TargetPath string
	Lines      *[]string
}

func (p *ContentPatch) Apply(lines []string) []string {
	var oldLIneIndex int = p.OldLineNumber - 1
	if p.PatchType == PatchDelete {
		return append(lines[:oldLIneIndex], lines[(oldLIneIndex+p.OldLineNumber):]...)
	}
	if p.PatchType == PatchReplace {
		// Delete the line
		temp := append(lines[:oldLIneIndex], lines[(oldLIneIndex+p.OldLineCount):]...)
		// Insert new content
		return append(temp[:oldLIneIndex], append(p.NewContent, temp[oldLIneIndex:]...)...)
	}
	// TODO: add test for insert
	if p.PatchType == PatchInsert {
		// Insert new content
		return append(lines[:oldLIneIndex], append(p.NewContent, lines[oldLIneIndex:]...)...)
	}
	return lines
}

func (p *FilePatch) Apply() {
	if p.FileOp == FileDelete {
		err := os.Remove(p.SourcePath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if p.FileOp == FileCreate {
		file, err := os.Create(p.SourcePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		str := strings.Join(*p.Lines, "\n")
		content := []byte(str)
		_, err = file.Write(content)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}
	if p.FileOp == FileRename || p.FileOp == DirectoryRename {
		oldPath := p.SourcePath
		newPath := p.TargetPath
		err := os.Rename(oldPath, newPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	if p.FileOp == DirectoryDelete {
		dirPath := p.TargetPath
		err := os.RemoveAll(dirPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
