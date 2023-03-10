package internal

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type InjectionContext struct {
	DotFilePath      string
	InjectionHistory map[string][]int
	InjectionContent map[string][]string
}

func (p *Program) Inject(ctx *InjectionContext) {
	for _, c := range p.Commands {
		if c.Inject == nil {
			continue
		}
		dirPath := filepath.Dir(ctx.DotFilePath)
		targetFilePath := path.Join(dirPath, c.Inject.FilePath.Name)
		lines := func() []string {
			lineContent := ctx.InjectionContent[targetFilePath]
			if lineContent == nil {
				content, _ := os.ReadFile(targetFilePath)
				return strings.Split(string(content), "\n")
			} else {
				return lineContent
			}
		}()
		injectionLine := []string{"// UNGEN: " + c.Inject.CmdString}
		injectionHistory := func() []int {
			hist := ctx.InjectionHistory[targetFilePath]
			if hist == nil {
				return make([]int, 0)
			} else {
				return hist
			}
		}()
		targetLine := calculateTargetLine(injectionHistory, c.Inject.TargetLine)
		oldLineIndex := targetLine - 1
		newContent := append(lines[:oldLineIndex], append(injectionLine, lines[oldLineIndex:]...)...)
		ctx.InjectionContent[targetFilePath] = newContent
		// Need to append the original TargetLine
		ctx.InjectionHistory[targetFilePath] = sortAndInsertHistory(injectionHistory, c.Inject.TargetLine)
	}
}
func sortAndInsertHistory(injHist []int, origLn int) []int {
	sort.Ints(injHist)

	insertI := sort.SearchInts(injHist, origLn)
	return append(injHist[:insertI], append([]int{origLn}, injHist[insertI:]...)...)
}

func calculateTargetLine(injectionHistory []int, originalTargetLine int) int {
	index := sort.SearchInts(injectionHistory, originalTargetLine+1)
	return originalTargetLine + index
}
