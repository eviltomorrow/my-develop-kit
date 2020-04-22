package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eviltomorrow/transfer-go/pb"
	"github.com/golang/protobuf/ptypes/wrappers"
)

// UploadFile 上传文件
type UploadFile struct {
}

// GetFileInfo g
func (uf *UploadFile) GetFileInfo(ctx context.Context, text *wrappers.StringValue) (*pb.FileInfo, error) {
	if text == nil {
		return nil, fmt.Errorf("参数路径不能为空")
	}
	if text.Value == "" {
		return nil, fmt.Errorf("参数路径不能为空")
	}

	var result = &pb.FileInfo{
		Path: filepath.Join(filepath.Dir(text.Value), filepath.Base(text.Value)),
	}
	info, err := os.Stat(text.Value)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		return result, nil
	}
	if info.IsDir() {
		result.Size = info.Size()
		result.IsDir = true
		result.LastMod = info.ModTime().Unix()
	} else {
		result.Size = info.Size()
		result.IsDir = false
		result.LastMod = info.ModTime().Unix()
		md5, err := calMD5(text.Value)
		if err != nil {
			return nil, fmt.Errorf("计算文件 MD5 值异常，nest error: %v", err)
		}
		result.Md5 = md5
	}
	return result, nil
}

// SetCheckPoint 获取检查记录文件
func (uf *UploadFile) SetCheckPoint(data *pb.CheckPoint, channel pb.UploadFile_SetCheckPointServer) error {
	if data == nil {
		return fmt.Errorf("参数不能为空")
	}
	if data.Path == "" {
		return fmt.Errorf("检查文件路径不能为空")
	}

	var cppath = filepath.Join(filepath.Dir(data.Path), filepath.Base(data.Path))
	if !strings.HasSuffix(cppath, ".cpf") {
		return fmt.Errorf("检查文件名不符合规范，请以 .cpf 为后缀")
	}

	fs, err := loadFileStat(data.FileInfo)
	if err != nil {
		return fmt.Errorf("初始化文件状态失败，nest error: %v", err)
	}

	var cp = &checkpoint{
		Path: cppath,
	}
	if err := cp.load(cppath); err != nil {
		prepare(cp, fs)
	}
	if err := cp.valid(fs); err != nil {
		prepare(cp, fs)
	}

	var i int64
	for i = 0; i < int64(len(cp.FileParts)); i++ {
		var result = &pb.FilePart{
			Offset:      cp.FileParts[i].Offset,
			Num:         i,
			Size:        cp.FileParts[i].Size,
			IsCompleted: cp.FileParts[i].IsCompleted,
		}
		if !result.IsCompleted {
			err = channel.Send(result)
			if err != nil {
				return err
			}
		}
	}
	return cp.dump()
}

// UploadFile 上传文件分片
func (uf *UploadFile) UploadFile(channel pb.UploadFile_UploadFileServer) error {
	var cp *checkpoint
	var f *os.File
	var current int64
	var lock fileLock
	var ft = new(filestat)
	var cppath string = "nil"

	defer func() {
		if cp != nil {
			cp.dump()
		}

		md5, _ := calMD5(filepath.Join(ft.BaseDir, ft.BaseName))
		if md5 != "" && md5 == ft.MD5 {
			os.Remove(cppath)
		}

		if f != nil {
			f.Close()
		}
		if lock != nil {
			lock.release()
		}

		os.Remove(filepath.Join(ft.BaseDir, fmt.Sprintf("%s.lock", ft.BaseName)))
	}()

	var check = func(info *filestat, data *pb.FileChannel) error {
		if data == nil {
			return fmt.Errorf("上传文件数据信息不能为空")
		}
		if data.FileInfo == nil {
			return fmt.Errorf("上传文件数据信息不能为空")
		}
		if data.FilePart == nil {
			return fmt.Errorf("上传文件数据分片不能为空")
		}

		if info.BaseDir == "" {
			info.BaseDir = filepath.Dir(data.FileInfo.Path)
		} else {
			if info.BaseDir != filepath.Dir(data.FileInfo.Path) {
				return fmt.Errorf("上传文件路径发生变化")
			}
		}
		if info.BaseName == "" {
			info.BaseName = filepath.Base(data.FileInfo.Path)
		} else {
			if info.BaseName != filepath.Base(data.FileInfo.Path) {
				return fmt.Errorf("上传文件路径发生变化")
			}
		}
		if info.Size == 0 {
			info.Size = data.FileInfo.Size
		} else {
			if info.Size != data.FileInfo.Size {
				return fmt.Errorf("上传文件大小发生变化")
			}
		}

		if info.LastMod == 0 {
			info.LastMod = data.FileInfo.LastMod
		} else {
			if info.LastMod != data.FileInfo.LastMod {
				return fmt.Errorf("上传文件时间戳发生变化")
			}
		}

		if info.MD5 == "" {
			info.MD5 = data.FileInfo.Md5
		} else {
			if info.MD5 != data.FileInfo.Md5 {
				return fmt.Errorf("上传文件 MD5 值发生变化")
			}
		}

		return nil
	}

loop:
	for {
		data, err := channel.Recv()
		if err == io.EOF {
			log.Printf("客户端上传数据分片结束\r\n")
			break
		}
		if err != nil {
			log.Printf("客户端上传数据切片异常，nest error: %v\r\n", err)
			return err
		}

		log.Printf("offset: %v, size: %v, data: %v", data.FilePart.Offset, data.FilePart.Size, len(data.FilePart.Data))
		// 文件锁
		if lock == nil {
			lock, err = newFileLock(filepath.Join(ft.BaseDir, fmt.Sprintf("%s.lock", ft.BaseName)), false)
			if err != nil {
				return fmt.Errorf("获取文件锁失败，nest error: %v", err)
			}
		}

		if err := check(ft, data); err != nil {
			return err
		}

		if cppath == "nil" {
			cppath = data.Checkpoint
		} else {
			if cppath != data.Checkpoint {
				return fmt.Errorf("断点记录文件路径发生变化")
			}
		}

		switch {
		case cppath == "nil":
			if f == nil {
				if f, err = setFile(data.Strategy, ft); err != nil {
					return err
				}
			}
			if f == nil {
				md5, err := calMD5(filepath.Join(ft.BaseDir, ft.BaseName))
				if err != nil {
					return fmt.Errorf("计算文件 MD5 值异常，nest error: %v", err)
				}
				if md5 != ft.MD5 {
					return fmt.Errorf("文件已经存在, MD5 不同")
				}
				return channel.SendAndClose(&wrappers.BoolValue{
					Value: true,
				})
			}

			err = writeFile(f, data.FilePart.Offset, data.FilePart.Data, current == data.FilePart.Offset)
			if err != nil {
				return fmt.Errorf("写入文件数据失败，nest error: %v", err)
			}
			current = data.FilePart.Offset + int64(len(data.FilePart.Data))
			if current == ft.Size {
				break
			}
		default:
			if cp == nil {
				cp = &checkpoint{}
				if err := cp.load(cppath); err != nil {
					return fmt.Errorf("加载断点记录文件失败，nest error: %v", err)
				}
				if err := cp.valid(ft); err != nil {
					return fmt.Errorf("验证断点记录文件失败，nest error: %v", err)
				}
			}
			if f == nil {
				if f, err = setFile(data.Strategy, ft); err != nil {
					return err
				}
			}
			if f == nil {
				return fmt.Errorf("文件已存在 MD5 不相同")
			}

			err = writeFile(f, data.FilePart.Offset, data.FilePart.Data, current == data.FilePart.Offset)
			if err != nil {
				return fmt.Errorf("写入文件数据失败，nest error: %v", err)
			}

			if err := cp.update(data.FilePart.Num); err != nil {
				return err
			}

			current = data.FilePart.Offset + int64(len(data.FilePart.Data))
			if current == ft.Size {
				break loop
			}
		}
	}

	return channel.SendAndClose(&wrappers.BoolValue{
		Value: true,
	})
}

func setFile(strategy pb.FileChannel_UploadStrategy, ft *filestat) (*os.File, error) {
	info, err := os.Stat(filepath.Join(ft.BaseDir, ft.BaseName))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("获取文件状态失败，nest error: %v", err)
		}
		if e1 := createDIR(ft.BaseDir); e1 != nil {
			return nil, fmt.Errorf("创建目录失败，nest error: %v", e1)
		}
		if e1 := createFile(filepath.Join(ft.BaseDir, ft.BaseName), ft.Size); e1 != nil {
			return nil, fmt.Errorf("创建文件失败，nest error: %v", e1)
		}
	} else if err == nil && info.IsDir() {
		if e1 := createFile(filepath.Join(ft.BaseDir, ft.BaseName), ft.Size); e1 != nil {
			return nil, fmt.Errorf("创建文件失败，nest error: %v", e1)
		}
	} else {
		switch strategy {
		case pb.FileChannel_EXIST_COVER:
			if info.Size() == ft.Size {
				md5, err := calMD5(filepath.Join(ft.BaseDir, ft.BaseName))
				if err != nil {
					return nil, fmt.Errorf("计算文件 MD5 值异常，nest error: %v", err)
				}
				if md5 != ft.MD5 {
					if e1 := os.Remove(filepath.Join(ft.BaseDir, ft.BaseName)); e1 != nil {
						return nil, fmt.Errorf("删除文件失败，nest error: %v", e1)
					}
					if e1 := createFile(filepath.Join(ft.BaseDir, ft.BaseName), ft.Size); e1 != nil {
						return nil, fmt.Errorf("创建文件失败，nest error: %v", e1)
					}
				}
				return nil, nil
			}
			if e1 := os.Remove(filepath.Join(ft.BaseDir, ft.BaseName)); e1 != nil {
				return nil, fmt.Errorf("删除文件失败，nest error: %v", e1)
			}
			if e1 := createFile(filepath.Join(ft.BaseDir, ft.BaseName), ft.Size); e1 != nil {
				return nil, fmt.Errorf("创建文件失败，nest error: %v", e1)
			}
		case pb.FileChannel_EXIST_FAILURE:
			if info.Size() != ft.Size {
				return nil, fmt.Errorf("文件已存在 size 不同")
			}
		default:
			return nil, fmt.Errorf("未支持测策略")
		}
	}

	file, err := os.OpenFile(filepath.Join(ft.BaseDir, ft.BaseName), os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败，nest error: %v", err)
	}

	return file, nil
}
