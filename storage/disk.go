package storage
/*
功能一：为服务端提供数据文件的读取能力
功能二：提供基于文件名的解析服务
 */

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	DataFileExtension string = ".data"
	CRCFileExtension  string = ".crc"
)

var (
	ErrorOpenDataFile = errors.New("error opening data file")
	ErrorReadDataFile = errors.New("error opening data file")
)

var (
	storagePath  = "../../data/"
)

// ReadFile 读取一个磁盘文件
func ReadFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, ErrorOpenDataFile
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if 	err != nil {
		return nil, ErrorReadDataFile
	}

	return content, nil
}

// parseDataFilename 解析数据文件的文件名+扩展名部分，并返回 FID 和 文件的块序号
func ParseDataFilename(fileBaseName string) (fid string, blockIndex int) {
	filename := strings.Split(fileBaseName, ".")[0]
	filenameSlice := strings.Split(filename, ":")
	fid = filenameSlice[0]
	blockIndex, _ = strconv.Atoi(filenameSlice[1])
	return fid, blockIndex
}


