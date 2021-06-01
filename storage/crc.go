package storage
/*
功能一：提供基础的 crc 文件读取功能
功能二：解析 crc 文件的内容
功能三：提供基于 crc 文件的校验方式
*/

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var (
	ErrorReadCRCFile = errors.New("failed to read crc file")
	ErrorParseCRC    = errors.New("failed to parse crc content")
	ErrorInvalidCRC  = errors.New("invalid crc content")
	ErrorVerifyCRC   = errors.New("error verifying crc of file")
)

type File struct {
	Files []Item `json:"files"`
}

type Item struct {
	FileLen  int64    `json:"file_len"`
	FileName string   `json:"file_name"`
	FilePath string   `json:"file_path"`
	FID      string   `json:"fid"`
	CRCArray []uint32 `json:"crc_array"`
}

// GetCRCFilePath 获取一个crc文件的路径
func GetCRCFilePath(fid string) string {
	p := storagePath + fid + CRCFileExtension
	fmt.Println("warning: ",p)
	return p
}

// LoadCRCFromFile 从磁盘加载 CRC 信息到内存索引
// 若发现已存在则直接返回
func LoadCRCFromFile(rid string) (*Item, error) {
	content, err := ioutil.ReadFile(GetCRCFilePath(rid))
	if err != nil {
		fmt.Println("failed to read CRC file of rid " + rid)
		return nil, ErrorReadCRCFile
	}
	return ParseCRC(rid, content)
}

// ParseCRC 从json格式的内容中解析出 CRC 数组并添加到索引
func ParseCRC(rid string, content []byte) (*Item, error) {
	crc := &File{}
	err := json.Unmarshal(content, crc)
	if err != nil {
		fmt.Println("failed to parse CRC content of rid " + rid)
		return nil, ErrorParseCRC
	}
	if len(crc.Files) < 1 {
		fmt.Println("invalid CRC content of rid " + rid)
		return nil, ErrorInvalidCRC
	}

	return &crc.Files[0], nil
}

// VerifyBlockCRC 校验一个2MB数据块(读出或写入)的CRC值
func VerifyBlockCRC(blockIndex int, block []byte, crcArray []uint32) bool {
	if blockIndex >= len(crcArray) {
		fmt.Println("block_index: ", blockIndex)
		fmt.Println("crc_array_len: ", len(crcArray))
		fmt.Println("error: ", ErrorVerifyCRC.Error())
		return false
	}

	blockCRC := crcArray[blockIndex]
	dataCRC := Checksum(block)
	if dataCRC != blockCRC {
		fmt.Println("crc: ", blockCRC)
		fmt.Println("data_crc: ", dataCRC)
		fmt.Println("block_index: ", blockIndex)
		fmt.Println("error: ", ErrorVerifyCRC.Error())
		return false
	}
	return true
}

// Checksum 计算标准的 IEEE802.3 CRC 结果
func Checksum(data []byte) uint32 {
	return crc32.ChecksumIEEE(data)
}

func CreateStoragePath() {
	storagePath = "./data"
	if storagePath == "" {
		fmt.Println("no storage path provided")
	}
	if storagePath[len(storagePath)-1:] != string(filepath.Separator) {
		storagePath += string(filepath.Separator)
	}
	err := os.MkdirAll(path.Dir(storagePath), 0755)
	if err != nil {
		fmt.Println("failed to create storage path ", storagePath)
	}
}
