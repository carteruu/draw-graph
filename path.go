package draw_graph

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func AppPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}
func AppDir() string {
	appPath := AppPath()
	index := strings.LastIndex(appPath, string(os.PathSeparator))
	appDirPath := appPath[:index]
	return appDirPath
}
