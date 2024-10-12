package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Node struct {
	Name     string
	SubNodes []*Node
}

func Init(path string) (*Node, error) {
	n := Node{
		Name: path,
	}
	return &n, nil
}

// addSubNodes recursively adds sub-nodes and prints directories and files as required.
func (n *Node) addSubNodes(path string, level int, isLastDir []bool) {
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	// Sort entries: files first, then directories, both alphabetically
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].IsDir() != entries[j].IsDir() {
			return !entries[i].IsDir() // Files first, then directories
		}
		return entries[i].Name() < entries[j].Name() // Alphabetical order
	})

	for i, entry := range entries {
		subNode := &Node{Name: entry.Name()}
		if strings.HasPrefix(subNode.Name, ".") {
			continue // Skip hidden files and directories
		}

		// Determine if this is the last entry
		isLastEntry := i == len(entries)-1

		// Print directory or file
		if entry.IsDir() {
			if isLastEntry {
				fmt.Println(genIndent(level, isLastDir), "└──", subNode.Name)

			} else {
				fmt.Println(genIndent(level, isLastDir), "├──", subNode.Name)
			}
			// Add to the "isLastDir" list for the next level
			n.addSubNodes(filepath.Join(path, subNode.Name), level+1, append(isLastDir, isLastEntry))
		} else {
			if isLastEntry {
				fmt.Println(genIndent(level, isLastDir), "└──", subNode.Name)
			} else {
				fmt.Println(genIndent(level, isLastDir), "├──", subNode.Name)
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
			indent += " │   " // Use vertical lines for intermediate directories
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
	return nil
}

func main() {
	path := "." // Default to current directory
	if len(os.Args) >= 2 {
		path = os.Args[1] // If an argument is provided, use it as the path
	}

	// Start the directory tree generation
	err := Steps(path + "/test")
	if err != nil {
		fmt.Println("Error:", err)
	}
}
