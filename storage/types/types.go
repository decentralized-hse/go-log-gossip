package types

import (
	"os"
)

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

type AppendOnlyFile struct {
	file *os.File
}

func OpenAppendOnlyFile(filePath string) (*AppendOnlyFile, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, OS_USER_RW)

	if err != nil {
		return nil, err
	}

	return &AppendOnlyFile{file: file}, nil
}

func (appendOnlyFile *AppendOnlyFile) Size() (int64, error) {
	fileStat, err := appendOnlyFile.file.Stat()
	if err != nil {
		return 0, err
	}

	return fileStat.Size(), nil
}

func (appendOnlyFile *AppendOnlyFile) Append(data string) (int, error) {
	return appendOnlyFile.file.WriteString(data)
}

func (appendOnlyFile *AppendOnlyFile) AppendBytes(data []byte) (int, error) {
	return appendOnlyFile.file.Write(data)
}

func (appendOnlyFile *AppendOnlyFile) AppendLine(data string) (int, error) {
	return appendOnlyFile.Append(data + "\n")
}

func (appendOnlyFile *AppendOnlyFile) Close() error {
	return appendOnlyFile.file.Close()
}

func (appendOnlyFile *AppendOnlyFile) ReadAt(buff []byte, offset int64) (int, error) {
	return appendOnlyFile.file.ReadAt(buff, offset)
}
