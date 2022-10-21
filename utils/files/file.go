package files

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

import (
	"github.com/lethexixin/go-funcs/utils/times"
)

// GetAppRootWDPath 获取项目根路径目录
func GetAppRootWDPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}

// GetFileSize  获取文件大小
func GetFileSize(file multipart.File) (int, error) {
	content, err := ioutil.ReadAll(file)
	return len(content), err
}

// GetFileExt 获取文件后缀
func GetFileExt(fileName string) string {
	return path.Ext(fileName)
}

// GetFileNameInfo 获取文件的名称和后缀名信息,filename可以是test.txt,也可以是/root/test/test.txt之类的
// nameInfo代表除去后缀名的文件名,ext代表后缀名
func GetFileNameInfo(filename string) (nameInfo, ext string) {
	filename = path.Base(filename)
	ext = GetFileExt(filename)
	return strings.TrimSuffix(filename, ext), ext
}

// CheckFileIsExist  检查文件或路径是否存在
func CheckFileIsExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

// CheckPathIsDir 判断是否为路径
func CheckPathIsDir(path string) bool {
	file, err := os.Stat(path)
	if err != nil {
		return false
	}
	return file.IsDir()
}

// CheckFileIsPermission 检查文件或路径是否有权限操作
func CheckFileIsPermission(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

// IsNotExistToMkDir  目录如果不存在则新建文件夹,新建文件夹失败则返回错误,newDir代表这个目录是不是刚刚新建的目录(即这是一个空目录)
func IsNotExistToMkDir(src string) (newDir bool, err error) {
	if isExist := CheckFileIsExist(src); !isExist {
		if err := ToMkdir(src); err != nil {
			return true, err
		}
	}
	return false, nil
}

// ToMkdir 新建文件夹,新建文件夹失败则返回错误
func ToMkdir(dirName string) error {
	return os.MkdirAll(dirName, os.ModePerm)
}

// MustOpen 最大限度地尝试打开文件
func MustOpen(fileName, filePath string) (file *os.File, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %s", err.Error())
	}

	src := dir + "/" + filePath
	if CheckFileIsPermission(src) {
		return nil, fmt.Errorf("files.CheckPermission permission denied, src: %s", src)
	}
	if _, err = IsNotExistToMkDir(src); err != nil {
		return nil, fmt.Errorf("files.IsNotExistMkDir, src: %s, err: %s", src, err.Error())
	}
	if file, err = os.OpenFile(src+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644); err != nil {
		return nil, fmt.Errorf("fail to open file :%s", err.Error())
	}

	return file, nil
}

// GetFilesAndDirs 获取指定目录下的所有文件名和目录
func GetFilesAndDirs(dirPth string, extensions []string) (files []string, dirs []string, err error) {
	dirsPath, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}

	pathSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, dir := range dirsPath {
		if dir.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, dirPth+pathSep+dir.Name())
			_, _, _ = GetFilesAndDirs(dirPth+pathSep+dir.Name(), extensions)
		} else {
			if len(extensions) > 0 {
				for _, ext := range extensions {
					// 输出指定后缀格式的文件
					if ok := strings.HasSuffix(dir.Name(), "."+ext); ok {
						files = append(files, dirPth+pathSep+dir.Name())
					}
				}
			} else {
				files = append(files, dirPth+pathSep+dir.Name())
			}
		}
	}

	return files, dirs, nil
}

// GetAllFiles 获取指定目录下的所有文件;
// childDir = true 代表输出子目录下的文件;
// ext代表后缀名,如果为空则输出全部类型的文件,如果不为空,则输出指定类型的文件;
func GetAllFiles(dirPth string, childDir bool, extensions []string) (files []string, err error) {
	var dirs []string
	dirsPath, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	pathSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, dir := range dirsPath {
		if childDir {
			if dir.IsDir() { // 目录, 递归遍历
				dirs = append(dirs, dirPth+pathSep+dir.Name())
				_, _ = GetAllFiles(dirPth+pathSep+dir.Name(), childDir, extensions)
			}
		}
		if len(extensions) > 0 {
			// 输出指定后缀格式的文件
			for _, ext := range extensions {
				if ok := strings.HasSuffix(dir.Name(), "."+ext); ok {
					files = append(files, dirPth+pathSep+dir.Name())
				}
			}
		} else {
			files = append(files, dirPth+pathSep+dir.Name())
		}

	}

	if childDir {
		// 读取子目录下文件
		for _, table := range dirs {
			temp, _ := GetAllFiles(table, childDir, extensions)
			for _, tmp := range temp {
				files = append(files, tmp)
			}
		}
	}

	return files, nil
}

// WriteLines 按行写入文件,dir指的是路径,filename指的是文件名,dataLines指的是具体的数据,lineEnding指的是分割符,如\n,
// writeTime指是否需要在行首写入时间,append指是否需要追加写入
func WriteLines(dir string, filename string, dataLines []string, lineEnding string, writeTime bool, append bool) error {
	if _, err := IsNotExistToMkDir(dir); err != nil {
		return err
	}

	filename = path.Join(dir, filename)
	if !append {
		_, err := os.Stat(filename)
		if !os.IsNotExist(err) {
			if errRemove := os.Remove(filename); errRemove != nil {
				return errRemove
			}
		}
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	buffer := bufio.NewWriter(file)
	for _, line := range dataLines {
		line = fmt.Sprintf("%s%s", line, lineEnding)
		if writeTime {
			line = fmt.Sprintf("time:%s\t%s", times.CurrentDateTime(), line)
		}
		if _, err = buffer.WriteString(line); err != nil {
			return err
		}
	}

	buffer.Buffered()
	if err = buffer.Flush(); err != nil {
		return err
	}
	return nil
}

// CopyFile 封装的文件拷贝方法,srcName源文件, dstName目标文件
func CopyFile(srcName string, dstName string) (int64, error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()

	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	//拷贝文件
	return io.Copy(dst, src)
}
