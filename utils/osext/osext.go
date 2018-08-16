package osext

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var cx, ce = executableClean()

func executableClean() (string, error) {
	p, err := executable()
	return filepath.Clean(p), err
}

// Executable returns an absolute path that can be used to
// re-invoke the current program.
// It may not be valid after the current program exits.
func Executable() (string, error) {
	return cx, ce
}

// Returns same path as Executable, returns just the folder
// path. Excludes the executable name and any trailing slash.
func ExecutableFolder() (string, error) {
	p, err := Executable()
	if err != nil {
		return "", err
	}

	p = filepath.Dir(p) + string(os.PathSeparator)
	return p, nil
}

//Check file exits
func Exits(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

//check Read Permission
func HaveReadPermission(filename string) bool {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)

	if err != nil {
		return false
	}
	defer file.Close()

	return !os.IsPermission(err)
}

//
func HaveWritePermission(filename string) bool {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)

	if err != nil {
		return false
	}
	defer file.Close()

	return !os.IsPermission(err)
}

func HaveRWPermission(filename string) bool {
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)

	if err != nil {
		return false
	}
	defer file.Close()

	return !os.IsPermission(err)
}

//Search dir file
func SearchDir(dir string) (fi []string) {
	d, err := os.Open(dir)

	if err != nil {
		return nil
	}
	defer d.Close()

	f, err := d.Readdir(-1)
	if err != nil {
		return nil
	}

	fi = make([]string, len(f))
	for _, finfo := range f {
		if finfo.Mode().IsRegular() {
			fi = append(fi, finfo.Name())
		}
	}

	return fi
}

//获取当前可执行文件的所在目录
//对于非 Windows 系统，以 / 作路径分隔符，对于 Windows 系统，以 \ 作路径分隔符。
func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

//获取当前文件的所在目录
func GetWd() string {
	wd, _ := os.Getwd()
	return wd + string(os.PathSeparator)
}

//获取当前执行的文件名称
func GetCurrentFileName(bContainExt bool) string {
	_, filename, _, _ := runtime.Caller(0)
	if bContainExt {
		return path.Base(filename)
	} else {
		return strings.Replace(path.Base(filename), filepath.Ext(filename), "", -1)
	}
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
