package file2

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"zzz_helper/internal/utils/string2"
)

const (
	linuxSplit = "/"
	winSplit   = "\\"
)

var (
	doublePattern = regexp.MustCompile(`/{2,}`)
	winPattern    = regexp.MustCompile(`^\w:\\`)
	winNetPattern = regexp.MustCompile(`^\\\\`)
)

func NewFileHelper(isWindows bool) *FileHelper {
	fileHelper := &FileHelper{
		isWindows: isWindows,
	}
	fileHelper.pathSplit = linuxSplit
	if isWindows {
		fileHelper.pathSplit = winSplit
	}
	return fileHelper
}

type FileHelper struct {
	isWindows bool
	pathSplit string
}

func (this *FileHelper) AbsJoin(dir, filename string) string {
	if dir == "" {
		if !this.isWindows && !strings.HasPrefix(filename, this.pathSplit) {
			filename = this.pathSplit + filename
		}
		return filename
	}
	return fmt.Sprintf("%s%s%s", this.TrimPathRight(dir), this.pathSplit, filename)
}
func (this *FileHelper) GetDir(path string) string {
	if strings.HasSuffix(path, this.GetPathSplit()) {
		return this.TrimPathRight(path)
	}

	path = this.FixSplit(path, false)
	i := strings.LastIndex(path, this.GetPathSplit())
	if i == -1 {
		return path
	}

	return path[:i]
}

func (this *FileHelper) Base(path string) string {
	path = this.FixSplit(path, false)
	i := strings.LastIndex(path, this.GetPathSplit())
	if i == -1 {
		return path
	}
	return path[i+1:]
}
func (this *FileHelper) ComparePathIsInSameOS(path1, path2 string) bool {
	os1 := ""
	if strings.HasPrefix(path1, linuxSplit) {
		os1 = "linux"
	} else {
		os1 = "windows"
	}

	os2 := ""
	if strings.HasPrefix(path2, linuxSplit) {
		os2 = "linux"
	} else {
		os2 = "windows"
	}

	return os1 == os2
}
func (this *FileHelper) AbsJoinDir(dir, filename string) string {
	if dir == "" {

		absPath := this.TrimPathRight(filename) + this.pathSplit
		if !this.isWindows && !strings.HasPrefix(filename, this.pathSplit) {
			absPath = this.pathSplit + absPath
		}
		return absPath
	}

	return this.FixSplit(filepath.Join(dir, filename), true)
}

// 修复分隔符
func (this *FileHelper) FixSplit(filename string, isDir bool) string {
	if winNetPattern.MatchString(filename) {
		suffix := strings.ReplaceAll(filename[2:], winSplit, linuxSplit)
		suffix = strings.ReplaceAll(suffix, linuxSplit, this.GetPathSplit())
		if !isDir {
			suffix = this.TrimPathRight(suffix)
		} else {
			suffix = this.AbsJoin(suffix, "")
		}
		filename = filename[:2] + suffix
		return filename
	}

	filename = strings.ReplaceAll(filename, winSplit, linuxSplit)
	// BUG: 修复mac上双斜杠异常问题
	filename = doublePattern.ReplaceAllString(filename, linuxSplit)

	filename = strings.ReplaceAll(filename, linuxSplit, this.GetPathSplit())

	if !isDir {
		filename = this.TrimPathRight(filename)
	} else {
		filename = this.AbsJoin(filename, "")
	}

	return filename
}

func (this *FileHelper) TrimPathRight(dir string) string {
	split := this.GetPathSplit()
	return strings.TrimRight(dir, split)
}
func (this *FileHelper) TrimPathLeft(dir string) string {
	split := this.GetPathSplit()
	return strings.TrimLeft(dir, split)
}
func (this *FileHelper) Compare(s1, s2 string) bool {
	if this.isWindows {
		return string2.CompareIgnoreCase(s1, s2)
	} else {
		return s1 == s2
	}
}
func (this *FileHelper) GetRelativePath(filename, dir string) (relativePath string) {
	if dir != "" {
		relativePath = this.TrimPrefix(filename, dir)
		relativePath = this.TrimPathLeft(relativePath)
	} else {
		relativePath = filename
	}
	return
}

func (this *FileHelper) TrimPrefix(s, prefix string) string {
	if len(s) >= len(prefix) && this.Compare(s[:len(prefix)], prefix) {
		return s[len(prefix):]
	}
	return s
}
func (this *FileHelper) HasPrefix(s, prefix string) bool {
	if len(s) >= len(prefix) && this.Compare(s[:len(prefix)], prefix) {
		return true
	}
	return false
}

func (this *FileHelper) Split(filename string) (chain []string) {
	chain = make([]string, 0)
	var tmps []string
	if winNetPattern.MatchString(filename) {
		tmps = strings.Split(filename[2:], this.GetPathSplit())
		tmps[0] = filename[:2] + tmps[0]
	} else {
		tmps = strings.Split(filename, this.GetPathSplit())

	}

	for _, tmp := range tmps {
		if tmp == "" {
			continue
		}
		chain = append(chain, tmp)
	}
	if !this.isWindows && strings.HasPrefix(filename, this.pathSplit) {
		chain = append([]string{this.pathSplit}, chain...)
	}
	return

}
func (this *FileHelper) GetPathSplit() string {

	return this.pathSplit
}

func (this *FileHelper) IsWindows() bool {

	return this.isWindows
}

// 拼接相对路径
//func (this *FileHelper) Splice(dir)  string{
//
//	filepath.Join(elem...)
//}

func (this *FileHelper) IsAbs(filename string) bool {
	status := false
	if this.isWindows {
		if winNetPattern.MatchString(filename) {
			return true
		}

		status = winPattern.MatchString(filename)
	} else {
		status = strings.HasPrefix(filename, this.GetPathSplit())
	}
	return status
}
