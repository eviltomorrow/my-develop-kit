package service

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"

	"github.com/eviltomorrow/my-develop-kit/transfer-go/pb"
)

const (
	defaultPartSize = 1024 * 32
)

type checkpoint struct {
	Path      string
	FileStat  *filestat
	FileParts []*filepart
}

func (cp *checkpoint) load(path string) error {
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(cp); err != nil {
		return err
	}
	cp.Path = path
	return nil
}

func (cp *checkpoint) valid(fs *filestat) error {
	if filepath.Join(cp.FileStat.BaseDir, cp.FileStat.BaseName) != filepath.Join(fs.BaseDir, fs.BaseName) {
		return fmt.Errorf("上传路径变更")
	}

	if cp.FileStat.Size != fs.Size {
		return fmt.Errorf("上传文件大小发生变化")
	}

	if cp.FileStat.LastMod != fs.LastMod {
		return fmt.Errorf("上传文件修改时间发生变化")
	}

	return nil
}

func (cp *checkpoint) dump() error {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(cp); err != nil {
		return err
	}
	file, err := os.OpenFile(cp.Path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func (cp *checkpoint) update(num int64) error {
	if num >= int64(len(cp.FileParts)) {
		return fmt.Errorf("记录分片数据溢出")
	}
	cp.FileParts[int(num)].IsCompleted = true
	return nil
}

func prepare(cp *checkpoint, stat *filestat) {
	cp.FileStat = stat
	cp.FileParts = splitFile(stat.Size, defaultPartSize)
}

type filestat struct {
	BaseDir  string
	BaseName string
	Size     int64
	MD5      string
	LastMod  int64
}

func (fs *filestat) getDataPath() string {
	return filepath.Join(fs.BaseDir, fmt.Sprintf("%s.data", fs.BaseName))
}

func (fs *filestat) getCheckpointPath() string {
	return filepath.Join(fs.BaseDir, fmt.Sprintf("%s.cpf", fs.BaseName))
}

type filepart struct {
	Offset      int64
	Size        int64
	Num         int64
	IsCompleted bool
}

func loadFileStat(info *pb.FileInfo) (*filestat, error) {
	if info.Path == "" {
		return nil, fmt.Errorf("参数 Path 不能为空")
	}
	fmt.Println(info.Size)
	if info.Size <= 0 {
		return nil, fmt.Errorf("参数 Size 无效")
	}
	if info.Md5 == "" {
		return nil, fmt.Errorf("参数 MD5 不能为空")
	}
	if info.LastMod == 0 {
		return nil, fmt.Errorf("参数 LastMode 无效")
	}

	return &filestat{
		BaseDir:  filepath.Dir(info.Path),
		BaseName: filepath.Base(info.Path),
		Size:     info.Size,
		MD5:      info.Md5,
		LastMod:  info.LastMod,
	}, nil
}

func splitFile(fileSize, partSize int64) []*filepart {
	var partNum = fileSize / partSize
	if partNum > 10000 {
		partSize = fileSize / (10000 - 1)
		partNum = fileSize / partSize
	}

	var parts = make([]*filepart, 0, partNum+1)
	for i := int64(0); i < partNum; i++ {
		var p = &filepart{
			Offset:      i * partSize,
			Size:        partSize,
			Num:         i,
			IsCompleted: false,
		}
		parts = append(parts, p)
	}

	if fileSize%partSize > 0 {
		var p = &filepart{
			Offset:      int64(len(parts)) * partSize,
			Size:        fileSize % partSize,
			Num:         int64(len(parts)) + 1,
			IsCompleted: false,
		}
		parts = append(parts, p)
	}
	return parts
}
