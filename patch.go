package main

type ContentPatch struct {
	PatchType     PatchType
	OldLineNumber int
	OldLineCount  int
	NewContent    []string
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
