package development

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	colorReset = "\033[0m"
	colorBlue  = "\033[34m" // For directories
	colorGreen = "\033[32m" // For files
)

const (
	hiddenFile = "."
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
		if !strings.HasPrefix(entry.Name(), hiddenFile) {
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

func Steps(paths []string) error {
	for path := range paths {
		absPath, err := filepath.Abs(strconv.Itoa(path)) // Convert to absolute paths
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
	}

	return nil
}

func handleFlags(args []string) ([]string, map[string]bool) {
	flags := map[string]bool{
		"fast":      false,
		"recursive": false,
		"verbose":   false,
	}

	// Default to current directory if no path is provided
	var path []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") || strings.HasPrefix(arg, "-") {
			// Handle flags
			switch arg {
			case "--fast", "-f":
				flags["fast"] = true
			case "--recursive", "-r":
				flags["recursive"] = true
			case "--verbose", "-v":
				flags["verbose"] = true
			default:
				fmt.Printf("Unknown flag: %s\n", arg)
				return []string{""}, nil
			}
		} else {
			// Assume it's the path if it's not a flag
			if arg == "." {
				path = make([]string, 0)
				path[0] = "."
				return path, nil
			}
			path = append(path, arg)
		}
	}
	return path, flags
}

func Development() {
	//args := os.Args[1:] // Skip the program name
	args := []string{"~/Downloads/Torrent\\ Downloads", "./test"}
	// Handle flags and get the path
	path, flags := handleFlags(args)
	if len(path) < 1 || flags == nil {
		os.Exit(1)
	}

	// Print the flags for debugging purposes
	//fmt.Printf("Flags: fast=%v, recursive=%v, verbose=%v\n", flags["fast"], flags["recursive"], flags["verbose"])

	// Start the directory tree generation
	err := Steps(path)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	os.Exit(0)
}
