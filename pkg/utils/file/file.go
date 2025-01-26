/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/26 21:22
@Name: file.go
@Description: file operation
*/

package file

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func Unzip(src, dst string) error {
	reader, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() { _ = reader.Close() }()

	if err := os.MkdirAll(dst, os.ModeDir|0775); err != nil {
		return err
	}

	for _, file := range reader.File {

		if strings.HasPrefix(file.Name, "__MACOSX") {
			continue
		}

		var decodeName string
		if file.Flags == 0 {
			//如果标致位是0,则是默认的本地编码,默认为gbk
			i := bytes.NewReader([]byte(file.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			//如果标志为是1<<11也就是2048,则是utf-8编码
			decodeName = file.Name
		}

		path := filepath.Join(dst, decodeName)
		if file.FileInfo().IsDir() {
			_ = os.MkdirAll(path, file.Mode())
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		n, err := io.Copy(targetFile, fileReader)
		if err != nil {
			return err
		}
		fmt.Printf("成功解压 %s ，共写入了 %d 个字符的数据\n", path, n)
	}
	return nil
}

func Remove(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return err
	}
	return nil
}

func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		fmt.Println("err1: ", err)
		return false, nil
	}
	fmt.Println("err: ", err)
	return false, err
}

func Mkdir(path string) error {
	return os.MkdirAll(path, os.ModeDir|0775)
}

func CheckAndCreatePath(path string) (err error) {
	ok, err := PathExist(path)
	if err != nil {
		return
	}
	if ok {
		return nil
	}
	err = Mkdir(path)
	return
}

func ListDir(path string) ([]fs.FileInfo, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	infos := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		infos = append(infos, info)
	}
	return infos, nil
}

func Copy(src string, dst string) error {
	fileStat, err := os.Stat(src)
	if err != nil {
		return err
	}
	if fileStat.IsDir() {
		err = copyTree(src, dst, nil)
	} else {
		_, err = copyData(src, dst, false)
	}
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return nil
}

func RenameFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return err
	}
	return nil
}

func Zip(src, dst string) error {
	os.RemoveAll(dst)
	zipfile, _ := os.Create(dst)
	defer zipfile.Close()
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	temp := strings.Split(src, "/")
	srcPrefix := strings.Join(temp[0:len(temp)-1], "/")
	filepath.Walk(src, func(path string, info os.FileInfo, _ error) error {
		if path == src {
			return nil
		}

		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, srcPrefix)

		if info.IsDir() {
			header.Name += `/`
		} else {
			header.Method = zip.Deflate
		}

		var writer io.Writer
		if !strings.HasSuffix(info.Name(), ".fiber.gz") {
			writer, _ = archive.CreateHeader(header)
		}

		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".fiber.gz") {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
	return nil
}

func CopyBigFile(src string, dst string) error {
	start := time.Now()
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	//关闭文件
	defer srcFile.Close()
	defer dstFile.Close()
	buf := make([]byte, 1024000) //切片缓冲区
	for {
		n, err := srcFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		tmp := buf[:n]
		dstFile.Write(tmp)
	}
	elapsed := time.Since(start)
	fmt.Printf("成功拷贝`%s`至`%s`，耗时：%s\n", src, dst, elapsed)
	return nil
}
