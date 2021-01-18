package imports

import (
	"golang.org/x/tools/imports"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var modregex = regexp.MustCompile(`module ([^\s]*)`)

func GoModuleRoot(dir string) (string, bool) {
	dir, err := filepath.Abs(dir)
	if err != nil {
		panic(err)
	}
	dir = filepath.ToSlash(dir)
	modDir := dir
	assumedPart := ""
	for {
		f, err := ioutil.ReadFile(filepath.Join(modDir, "go.mod"))
		if err == nil {
			// found it, stop searching
			return string(modregex.FindSubmatch(f)[1]) + assumedPart, true
		}

		assumedPart = "/" + filepath.Base(modDir) + assumedPart
		parentDir, err := filepath.Abs(filepath.Join(modDir, ".."))
		if err != nil {
			panic(err)
		}

		if parentDir == modDir {
			// Walked all the way to the root and didnt find anything :'(
			break
		}
		modDir = parentDir
	}
	return "", false
}

func Imports(filename string, src []byte) ([]byte, error){
	return imports.Process(filename, src, &imports.Options{FormatOnly: false, Comments: true, TabIndent: true, TabWidth: 8, Fragment: true})
}


