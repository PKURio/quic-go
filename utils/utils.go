package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

const (
	S2CPacketSize = 1052
	C2SPacketSize = 40
)

type S2CProtocol struct {
	FID         [20]byte
	FileIdx     int16
	reserved    [2]byte
	PktIdx      int16
	PayloadSize int16
	Payload     [1024]byte
}

type C2SProtocol struct {
	FID      [20]byte
	FileIdx  int16
	reserved [2]byte
	PktIdx   [8]int16
}

func S2CMarshal(fid string, FileIdx int16, PktIdx int16, payloadSize int16, payload []byte) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	reserved := [2]byte{}
	FID, _ := hex.DecodeString(fid)
	binary.Write(bytesBuffer, binary.BigEndian, &FID)
	binary.Write(bytesBuffer, binary.BigEndian, &FileIdx)
	binary.Write(bytesBuffer, binary.BigEndian, &reserved)
	binary.Write(bytesBuffer, binary.BigEndian, &PktIdx)
	binary.Write(bytesBuffer, binary.BigEndian, &payloadSize)
	binary.Write(bytesBuffer, binary.BigEndian, &payload)
	return bytesBuffer.Bytes()
}

func C2SMarshal(fid string, FileIdx int16, PktIdx []int16) []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	reserved := [2]byte{}
	FID, _ := hex.DecodeString(fid)
	binary.Write(bytesBuffer, binary.BigEndian, &FID)
	binary.Write(bytesBuffer, binary.BigEndian, &FileIdx)
	binary.Write(bytesBuffer, binary.BigEndian, &reserved)
	for i := 0; i < 8; i++ {
		binary.Write(bytesBuffer, binary.BigEndian, &PktIdx[i])
	}
	return bytesBuffer.Bytes()
}

func S2CUnmarshal(S2CPkt []byte) (content S2CProtocol) {
	bytesBuffer := bytes.NewBuffer(S2CPkt)
	binary.Read(bytesBuffer, binary.BigEndian, &content.FID)
	binary.Read(bytesBuffer, binary.BigEndian, &content.FileIdx)
	binary.Read(bytesBuffer, binary.BigEndian, &content.reserved)
	binary.Read(bytesBuffer, binary.BigEndian, &content.PktIdx)
	binary.Read(bytesBuffer, binary.BigEndian, &content.PayloadSize)
	binary.Read(bytesBuffer, binary.BigEndian, &content.Payload)
	return
}

func C2SUnmarshal(C2SPkt []byte) (content C2SProtocol) {
	bytesBuffer := bytes.NewBuffer(C2SPkt)
	binary.Read(bytesBuffer, binary.BigEndian, &content.FID)
	binary.Read(bytesBuffer, binary.BigEndian, &content.FileIdx)
	binary.Read(bytesBuffer, binary.BigEndian, &content.reserved)
	for i := 0; i < 8; i++ {
		binary.Read(bytesBuffer, binary.BigEndian, &content.PktIdx[i])
	}
	return
}
