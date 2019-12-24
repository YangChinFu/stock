package path

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Exec exec return the path of executable
func Exec() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))

	return path[:index+1]
}
