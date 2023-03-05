package main

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type InjectionContext struct {
	dotFilePath      string
	injectionHistory map[string][]int
	injectionContent map[string][]string
}

func (p *Program) Inject(ctx *InjectionContext) {
	for _, c := range p.Commands {
		if c.Inject == nil {
			continue
		}
		dirPath := filepath.Dir(ctx.dotFilePath)
		targetFilePath := path.Join(dirPath, c.Inject.FilePath.Name)
		lines := func() []string {
			lineContent := ctx.injectionContent[targetFilePath]
			if lineContent == nil {
				content, _ := os.ReadFile(targetFilePath)
				return strings.Split(string(content), "\n")
			} else {
				return lineContent
			}
		}()
		injectionLine := []string{"// UNGEN: " + c.Inject.CmdString}
		injectionHistory := func() []int {
			hist := ctx.injectionHistory[targetFilePath]
			if hist == nil {
				return make([]int, 0)
			} else {
				return hist
			}
		}()
		targetLine := calculateTargetLine(injectionHistory, c.Inject.TargetLine)
		oldLineIndex := targetLine - 1
		newContent := append(lines[:oldLineIndex], append(injectionLine, lines[oldLineIndex:]...)...)
		ctx.injectionContent[targetFilePath] = newContent
		// Need to append the original TargetLine
		ctx.injectionHistory[targetFilePath] = sortAndInsertHistory(injectionHistory, c.Inject.TargetLine)
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
