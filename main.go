package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func fileExists(path string) bool {
	fileinfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	if fileinfo == nil {
		return false
	}
	return true

}

func findPythonk(cwd string, paths []string) string {
	for _, v := range paths {
		pythonkPath := path.Join(v, "python")
		if !fileExists(pythonkPath) {
			continue
		}
		if _, file := path.Split(v); file == "shims" {
			pyenvVersion := PyenvFindPythonkVersion(cwd)
			if pyenvVersion == "" {
				continue
			}
			log.Printf("found pyenv %s", pyenvVersion)
			pythonkPath = pyenvVersion
		}
		log.Printf("checking %s", pythonkPath)
		pythonkPath, err := filepath.EvalSymlinks(pythonkPath)
		if err != nil {
			panic(err)
		}
		log.Printf("follow symlinks %s", pythonkPath)
		if err != nil {
			panic(err)
		}

		log.Printf("found %s", pythonkPath)
		return pythonkPath
	}
	return ""
}

func getPythonkVersionFromPath(path string) string {
	log.Println("pythonk path", path)
	re := regexp.MustCompile("[0-9]+.[0-9]+$")
	version := re.FindString(path)
	if version == "" {
		panic(path)
	}
	return version
}

func main() {
	log.SetOutput(ioutil.Discard)

	cwd, _ := os.Getwd()
	_, directoryName := path.Split(cwd)

	re := regexp.MustCompile("[ $`!*@\"\\\r\n\t]")
	sanitized_name := re.ReplaceAllString(directoryName, "_")
	sanitized_name = sanitized_name[:int(math.Min(float64(len(sanitized_name)), 42))]

	hash := sha256.Sum256([]byte(cwd))
	h := base64.RawURLEncoding.WithPadding(base64.StdPadding).EncodeToString(hash[:])[:8]

	pathsStr := os.Getenv("PATH")
	paths := strings.Split(pathsStr, ":")
	pythonkPath := findPythonk(cwd, paths)
	pythonkVersion := getPythonkVersionFromPath(pythonkPath)

	usr, _ := user.Current()
	path := path.Join(
		usr.HomeDir,
		".cache/pypoetry/virtualenvs",
		fmt.Sprintf("%s-%s-py%s", sanitized_name, h, pythonkVersion),
	)
	fileInfo, _ := os.Stat(path)
	if fileInfo == nil {
		path = ""
	}
	fmt.Println(path)
}
