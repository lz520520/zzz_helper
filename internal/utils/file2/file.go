package file2

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileContentSplit(content string) (resultSlice []string) {
	resultSlice = make([]string, 0)
	resultSlice = strings.Split(content, "\n")
	for i, c := range resultSlice {
		resultSlice[i] = strings.TrimSpace(c)
	}
	return
}

// 获取文件名，如C:/test.exe 返回 test
func GetBaseName(name string) string {
	filenameWithSuffix := filepath.Base(name)
	fileSuffix := filepath.Ext(filenameWithSuffix)
	return strings.TrimSuffix(filenameWithSuffix, fileSuffix)
}

// 获取文件名，如C:/test.exe 返回 test
func GetBaseNameWithSuffix(name string) string {
	filenameWithSuffix := filepath.Base(name)
	return filenameWithSuffix
}
func ReadFile(filename string) (resultSlice []string, err error) {
	resultSlice = make([]string, 0)
	resultBytes, err := ReadFileBytes(filename)
	if err != nil {
		return
	}
	resultSlice = strings.Split(string(resultBytes), "\n")
	for i, c := range resultSlice {
		resultSlice[i] = strings.TrimSpace(c)
	}
	return resultSlice, nil

	/*	file, err := os.Open(filename)
		if err != nil {
			return
		}
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			b, err := reader.ReadString('\n')
			if err != nil && err != io.EOF {
				return resultSlice, err
			}
			if b == "" {
				break
			}

			resultSlice = append(resultSlice, strings.TrimSpace(b))
		}*/

}

func ReadFileBytes(filename string) (resultBytes []byte, err error) {
	resultBytes = make([]byte, 0)
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return resultBytes, err
		}
		if n == 0 {
			break
		}

		resultBytes = append(resultBytes, buf[:n]...)
	}
	return resultBytes, nil
}

func WriteFile(filename string, writeBytes []byte) (err error) {
	MkdirFromFile(filename)
	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(writeBytes)
	return
}

func AppendFile(filename string, writeBytes []byte) (err error) {
	var file *os.File
	file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(writeBytes)
	return
}
func CreateFile(filename string) (f *os.File, err error) {
	MkdirFromFile(filename)
	f, err = os.Create(filename)
	return

}
func MkdirFromFile(src string) error {
	dstDir := filepath.Dir(src)
	_, err := os.Stat(dstDir)
	if err != nil {
		err = os.MkdirAll(dstDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}
func Mkdir(src string) error {
	_, err := os.Stat(src)
	if err != nil {
		err = os.MkdirAll(src, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyDir(dest, src string) error {
	// 获取源目录下的所有文件和子目录
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// 创建目标目录
	if err := os.MkdirAll(dest, os.ModePerm); err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		if entry.IsDir() {
			// 递归拷贝子目录
			if err = CopyDir(destPath, srcPath); err != nil {
				return err
			}
		} else {
			// 拷贝文件
			if _, err = CopyFile(destPath, srcPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.Create(dstName)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func MoveFile(dstName, srcName string) error {
	// 尝试将源文件移动到目标位置
	err := os.Rename(srcName, dstName)
	if err != nil {
		return err
	}
	return nil
}

func CopyFileAttachBytes(dstName, srcName string, attach []byte) (err error) {
	input, err := os.ReadFile(srcName)
	if err != nil {
		return
	}
	input = append(input, attach...)

	err = os.WriteFile(dstName, input, 0644)
	return
}

func FileHasSuffix(filename string, subBytes []byte) bool {
	input, err := os.ReadFile(filename)
	if err != nil {
		return false
	}
	return bytes.HasSuffix(input, subBytes)

}

func GetAllFile(pathname string, suffix string) ([]string, error) {
	rd, err := os.ReadDir(pathname)
	if err != nil {
		return nil, err
	}
	result := make([]string, 0)
	for _, fi := range rd {
		if !fi.IsDir() {
			//fullName := pathname + "/" +
			if strings.HasSuffix(fi.Name(), suffix) {
				result = append(result, fi.Name())
			}
		}
	}
	return result, nil

}

func GetAllSubFile(pathname string, suffix []string) ([]string, error) {
	result := make([]string, 0)
	rd, err := os.ReadDir(pathname)
	if err != nil {
		return result, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			//fmt.Printf("[%s]\n", pathname+"\\"+fi.Name())
			tmpResult, err := GetAllSubFile(filepath.Join(pathname, fi.Name()), suffix)
			if err == nil {
				result = append(result, tmpResult...)
			}
		} else {
			filterFlags := false
			for _, ext := range suffix {
				if strings.TrimLeft(filepath.Ext(fi.Name()), ".") == ext {
					filterFlags = true
				}
			}
			if !filterFlags {
				result = append(result, filepath.Join(pathname, fi.Name()))
			}
			//fmt.Println(filepath.Join(pathname,fi.Name()))
		}
	}
	return result, nil
}

func GetCurrentDir() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	return exPath
}

func GetAbsPath(src string) string {
	dst := ""
	if !filepath.IsAbs(src) {
		if absPath, err := filepath.Abs(src); err == nil {
			dst = absPath
		} else {
			dst = src
		}
	} else {
		dst = src
	}
	return dst
}
