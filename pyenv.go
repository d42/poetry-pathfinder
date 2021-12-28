package main

import (
	"log"
	"os"
	"path"
	"strings"
)

func pyenvReadPythonkVersionFile(dir string) string {
	for {
		log.Println("pyenv", dir)
		if dir == "/" {
			return ""
		}
		pythonVersionPath := path.Join(dir, ".python-version")
		if fileExists(pythonVersionPath) {
			dat, err := os.ReadFile(pythonVersionPath)
			if err != nil {
				panic(err)
			}
			return strings.Replace(string(dat), "\n", "", -1)
		}
		return ""

	}
}

func PyenvFindPythonkVersion(dir string) string {
	pyenvRoot := os.Getenv("PYENV_ROOT")
	if pyenvRoot == "" {
		return ""
	}
	version := pyenvReadPythonkVersionFile(dir)
	if version == "" {
		return ""
	}
	return path.Join(pyenvRoot, "versions", version, "bin", "python")
}
