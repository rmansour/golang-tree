package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	colorReset = "\033[0m"
	colorBlue  = "\033[34m" // For directories
	colorGreen = "\033[32m" // For files
)

type Node struct {
	Name      string
	DirCount  int
	FileCount int
	SubNodes  []*Node
}

func Init(path string) (*Node, error) {
	n := Node{
		Name: path,
	}
	return &n, nil
}

func (n *Node) addSubNodes(path string, level int, isLastDir []bool) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Filter out hidden files and directories (those starting with ".")
	var visibleEntries []os.DirEntry
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			visibleEntries = append(visibleEntries, entry)
		}
	}

	// Sort entries: files first, then directories, both alphabetically
	sort.SliceStable(visibleEntries, func(i, j int) bool {
		if visibleEntries[i].IsDir() != visibleEntries[j].IsDir() {
			return !visibleEntries[i].IsDir() // Files first, then directories
		}
		return visibleEntries[i].Name() < visibleEntries[j].Name() // Alphabetical order
	})

	// Process only the visible entries
	for i, entry := range visibleEntries {
		subNode := &Node{Name: entry.Name()}
		// Determine if this is the last entry
		isLastEntry := i == len(visibleEntries)-1

		// Print directory or file
		if entry.IsDir() {
			n.DirCount++
			if isLastEntry {
				fmt.Printf("%s└── %s%s%s\n", genIndent(level, isLastDir), colorBlue, subNode.Name, colorReset)
			} else {
				fmt.Printf("%s├── %s%s%s\n", genIndent(level, isLastDir), colorBlue, subNode.Name, colorReset)
			}
			// Add subdirectories recursively
			subNode.addSubNodes(filepath.Join(path, subNode.Name), level+1, append(isLastDir, isLastEntry))
		} else {
			n.FileCount++
			if isLastEntry {
				fmt.Printf("%s└── %s%s%s\n", genIndent(level, isLastDir), colorGreen, subNode.Name, colorReset)
			} else {
				fmt.Printf("%s├── %s%s%s\n", genIndent(level, isLastDir), colorGreen, subNode.Name, colorReset)
			}
		}

		// Add the sub-node to its parent
		n.SubNodes = append(n.SubNodes, subNode)
	}
}

func genIndent(level int, isLastDir []bool) string {
	indent := ""
	for i := 0; i < level-1; i++ {
		if isLastDir[i] {
			indent += "    " // Use spaces for the last directory in the parent
		} else {
			indent += "│   " // Use vertical lines for intermediate directories
		}
	}
	return indent
}

func Steps(path string) error {
	absPath, err := filepath.Abs(path) // Convert to absolute path
	if err != nil {
		return err
	}

	n, err := Init(absPath)
	if err != nil {
		return err
	}
	fmt.Println(n.Name)
	n.addSubNodes(n.Name, 1, []bool{})

	if n.DirCount >= 0 || n.FileCount >= 0 {
		fmt.Printf("\n%v directories, %v files\n", n.DirCount, n.FileCount)
	}
	return nil
}

func main() {
	path := "." // Default to current directory
	if len(os.Args) >= 2 {
		for i := 1; i < len(os.Args); i++ {
			path = os.Args[i] // If an argument is provided, use it as the path
			// Start the directory tree generation
			err := Steps(path)
			if err != nil {
				fmt.Println("Error:", err)
				os.Exit(1)
			}
		}
	} else {
		err := Steps(path)
		if err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
	}
	os.Exit(0)

}
