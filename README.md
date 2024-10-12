# golang-tree
This is an implementation of the famous "tree" command for Linux/Unix.

### Supported features:
- [x] List one or more folders
- [x] Colouring of folders and files
- [x] Total number of folders and files

### To be added:
- [ ] Files first
- [ ] Only Folders
- [ ] Depth to list
- [ ] Print hidden
- [ ] Executables for the different OS's
---
### Installation instructions:

This is a Golang project, hence, to run it, please install Golang to build and run the project.
Once Golang is installed, `cd` to the files directory and run `go build`. This will output a file with the name `golang-tree` by default.

Command-line steps:
```
git clone https://github.com/rmansour/golang-tree.git && cd golang-tree && go build
```

After that, execute it by running: `./golang-tree <path(s)>` for any UNIX operating system, or `.\golang-tree <path(s)>` for Windows.
