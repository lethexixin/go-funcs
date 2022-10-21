package zips

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Unzip 解压zip, archive: 源文件路径(非目录), target: 解压目标目录
func Unzip(archive, target string) error {
	reader, err := zip.OpenReader(archive)
	if err != nil {
		return fmt.Errorf("zip.OpenReader(%s) err: %s", archive, err.Error())
	}
	if err = os.MkdirAll(target, 0755); err != nil {
		return fmt.Errorf("os.MkdirAll(%s) err: %s", target, err.Error())
	}
	for _, file := range reader.File {
		filePath := path.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(filePath, file.Mode())
			continue
		}

		fileReader, errOpen := file.Open()
		if errOpen != nil {
			return fmt.Errorf("file.Open err: %s", errOpen.Error())
		}

		targetFile, errOpenFile := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if errOpenFile != nil {
			return fmt.Errorf("os.OpenFile(%s) err: %s", filePath, errOpenFile.Error())
		}

		_, err = io.Copy(targetFile, fileReader)
		_ = fileReader.Close()
		_ = targetFile.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Zip 多个文件压缩为zip, filePaths: 需要压缩的文件的路径, zipPath: 最终压缩后的文件路径
func Zip(filePaths []string, zipPath string) (err error) {
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	if len(filePaths) > 0 {
		for _, filepath := range filePaths {
			filepathList := strings.Split(filepath, `\`)
			filepathTemp := filepathList[len(filepathList)-1]
			zipFile, err := w.Create(filepathTemp)
			if err != nil {
				return fmt.Errorf("w.Create(%s) err: %s", filepathTemp, err.Error())
			}
			file, err := os.Open(filepath)
			if err != nil {
				return fmt.Errorf("os.Open(%s) err: %s", filepath, err.Error())
			}
			bs, err := ioutil.ReadAll(file)
			if err != nil {
				return fmt.Errorf("ioutil.ReadAll() err: %s", err.Error())
			}
			_, _ = zipFile.Write(bs)
		}
	}
	_ = w.Close()

	f, err := os.OpenFile(zipPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("os.OpenFile(%s) err: %s", zipPath, err.Error())
	}
	defer f.Close()
	_, err = buf.WriteTo(f)
	return err
}
