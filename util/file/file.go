/**
2 * @Author: Nico
3 * @Date: 2021/1/18 0:39
4 */
package file

import (
	"github.com/zeta-io/zctl/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	return false, err
}

func GetTpls(path string) ([]string, error) {
	tpls := make([]string, 0)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasSuffix(info.Name(), "tpl") {
			tpls = append(tpls, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tpls, nil
}

func GetDir(path string) (string, error){
	index := strings.LastIndex(path, string(os.PathSeparator))
	if index == -1 {
		return "", errors.ErrOutputIsNotFile
	}
	return path[0:index], nil
}

func Read(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	return b, err
}

func Write(path string, source []byte) error {
	index := strings.LastIndex(path, string(os.PathSeparator))
	if index == -1 {
		return errors.ErrOutputIsNotFile
	}
	err := os.MkdirAll(path[0:index], 0666)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, source, 0666)
}
