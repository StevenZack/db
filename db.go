package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	dbDir string
)

type Cmd struct {
	path       string
	resultData []byte
	resultErr  error
}

func InitDB(dir string) error {
	fi, e := os.Stat(dir)
	if e != nil {
		if os.IsNotExist(e) {
			e = os.MkdirAll(dir, 0755)
			if e != nil {
				return e
			}
			dbDir = getrpath(dir)
			return nil
		}
		return e
	}
	if !fi.IsDir() {
		return errors.New("db : initDB() " + dir + " is not a dir ")
	}
	dbDir = getrpath(dir)
	return nil
}
func getrpath(s string) string {
	if len(s) < 1 || s[len(s)-1:] == string(os.PathSeparator) {
		return s
	}
	return s + string(os.PathSeparator)
}
func Get(s string) *Cmd {
	cmd := &Cmd{}
	cmd.path = dbDir + s
	fi, e := os.Stat(cmd.path)
	if e != nil {
		if os.IsNotExist(e) {
			return cmd
		}
		cmd.resultErr = e
		return cmd
	}
	if fi.IsDir() {
		cmd.resultErr = errors.New("db.Get() : " + cmd.path + "is dir")
		return cmd
	}
	f, e := os.OpenFile(cmd.path, os.O_RDONLY, 0644)
	if e != nil {
		cmd.resultErr = e
		return cmd
	}
	defer f.Close()
	cmd.resultData, e = ioutil.ReadAll(f)
	if e != nil {
		cmd.resultErr = e
		return cmd
	}
	return cmd
}
func Set(key string, data interface{}) error {
	value, e := json.Marshal(data)
	if e != nil {
		return e
	}
	path := dbDir + key
	fi, e := os.Stat(path)
	if e != nil {
		if os.IsNotExist(e) {
			return writeFileBytes(path, value)
		}
		return e
	}
	if fi.IsDir() {
		return errors.New("db.Set() : " + path + " is a dir")
	}
	return writeFileBytes(path, value)
}
func getDirOfFilePath(path string) (string, error) {
	sep := string(os.PathSeparator)
	if path[len(path)-1:] == sep {
		return "", errors.New("db.getDirOfFilePath(" + path + ") failed : path end with / (or \\)")
	}
	for i := len(path) - 1; i > -1; i-- {
		if path[i:i+1] == sep {
			return path[:i+1], nil
		}
	}
	return "", errors.New("db.getDirOfFilePath(" + path + ") failed : path incorrect")
}
func writeFileBytes(path string, value []byte) error {
	dir, e := getDirOfFilePath(path)
	if e != nil {
		return e
	}
	e = os.MkdirAll(dir, 0755)
	if e != nil {
		return e
	}
	f, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if e != nil {
		return e
	}
	defer f.Close()
	_, e = f.Write(value)
	return e
}
func (cmd *Cmd) Scan(i interface{}) error {
	if cmd.resultErr != nil {
		return cmd.resultErr
	}
	e := json.Unmarshal(cmd.resultData, i)
	return e
}