package main

type PatchType int

const (
	PatchReplace = iota
	PatchDelete
	PatchInsert
)

func (t PatchType) String() string {
	return [...]string{"Replace", "Delete", "Insert"}[t]
}

type Patch struct {
	PatchType PatchType
	OldLineNumber int
	OldLineCount int
	NewContent string
}

func Eval(fileText string, p *Program, programLineNumber int) []Patch {
	patch := Patch{
		PatchType: PatchReplace,
		OldLineNumber: 1,
		OldLineCount: 2,
		NewContent: "test",
	}
	return []Patch{patch}
}