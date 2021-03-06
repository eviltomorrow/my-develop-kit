// Code generated by protoc-gen-go. DO NOT EDIT.
// source: upload.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type FileChannel_UploadStrategy int32

const (
	FileChannel_EXIST_FAILURE FileChannel_UploadStrategy = 0
	FileChannel_EXIST_COVER   FileChannel_UploadStrategy = 1
)

var FileChannel_UploadStrategy_name = map[int32]string{
	0: "EXIST_FAILURE",
	1: "EXIST_COVER",
}

var FileChannel_UploadStrategy_value = map[string]int32{
	"EXIST_FAILURE": 0,
	"EXIST_COVER":   1,
}

func (x FileChannel_UploadStrategy) String() string {
	return proto.EnumName(FileChannel_UploadStrategy_name, int32(x))
}

func (FileChannel_UploadStrategy) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_91b94b655bd2a7e5, []int{3, 0}
}

type FileInfo struct {
	Path                 string   `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	Size                 int64    `protobuf:"varint,2,opt,name=size,proto3" json:"size,omitempty"`
	IsDir                bool     `protobuf:"varint,3,opt,name=isDir,proto3" json:"isDir,omitempty"`
	Md5                  string   `protobuf:"bytes,4,opt,name=md5,proto3" json:"md5,omitempty"`
	LastMod              int64    `protobuf:"varint,5,opt,name=lastMod,proto3" json:"lastMod,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileInfo) Reset()         { *m = FileInfo{} }
func (m *FileInfo) String() string { return proto.CompactTextString(m) }
func (*FileInfo) ProtoMessage()    {}
func (*FileInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_91b94b655bd2a7e5, []int{0}
}

func (m *FileInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileInfo.Unmarshal(m, b)
}
func (m *FileInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileInfo.Marshal(b, m, deterministic)
}
func (m *FileInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileInfo.Merge(m, src)
}
func (m *FileInfo) XXX_Size() int {
	return xxx_messageInfo_FileInfo.Size(m)
}
func (m *FileInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_FileInfo.DiscardUnknown(m)
}

var xxx_messageInfo_FileInfo proto.InternalMessageInfo

func (m *FileInfo) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *FileInfo) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *FileInfo) GetIsDir() bool {
	if m != nil {
		return m.IsDir
	}
	return false
}

func (m *FileInfo) GetMd5() string {
	if m != nil {
		return m.Md5
	}
	return ""
}

func (m *FileInfo) GetLastMod() int64 {
	if m != nil {
		return m.LastMod
	}
	return 0
}

type CheckPoint struct {
	Path                 string    `protobuf:"bytes,1,opt,name=path,proto3" json:"path,omitempty"`
	FileInfo             *FileInfo `protobuf:"bytes,2,opt,name=fileInfo,proto3" json:"fileInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *CheckPoint) Reset()         { *m = CheckPoint{} }
func (m *CheckPoint) String() string { return proto.CompactTextString(m) }
func (*CheckPoint) ProtoMessage()    {}
func (*CheckPoint) Descriptor() ([]byte, []int) {
	return fileDescriptor_91b94b655bd2a7e5, []int{1}
}

func (m *CheckPoint) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CheckPoint.Unmarshal(m, b)
}
func (m *CheckPoint) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CheckPoint.Marshal(b, m, deterministic)
}
func (m *CheckPoint) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CheckPoint.Merge(m, src)
}
func (m *CheckPoint) XXX_Size() int {
	return xxx_messageInfo_CheckPoint.Size(m)
}
func (m *CheckPoint) XXX_DiscardUnknown() {
	xxx_messageInfo_CheckPoint.DiscardUnknown(m)
}

var xxx_messageInfo_CheckPoint proto.InternalMessageInfo

func (m *CheckPoint) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *CheckPoint) GetFileInfo() *FileInfo {
	if m != nil {
		return m.FileInfo
	}
	return nil
}

type FilePart struct {
	Num                  int64    `protobuf:"varint,1,opt,name=num,proto3" json:"num,omitempty"`
	Offset               int64    `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	Size                 int64    `protobuf:"varint,4,opt,name=size,proto3" json:"size,omitempty"`
	IsCompleted          bool     `protobuf:"varint,5,opt,name=isCompleted,proto3" json:"isCompleted,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FilePart) Reset()         { *m = FilePart{} }
func (m *FilePart) String() string { return proto.CompactTextString(m) }
func (*FilePart) ProtoMessage()    {}
func (*FilePart) Descriptor() ([]byte, []int) {
	return fileDescriptor_91b94b655bd2a7e5, []int{2}
}

func (m *FilePart) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FilePart.Unmarshal(m, b)
}
func (m *FilePart) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FilePart.Marshal(b, m, deterministic)
}
func (m *FilePart) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FilePart.Merge(m, src)
}
func (m *FilePart) XXX_Size() int {
	return xxx_messageInfo_FilePart.Size(m)
}
func (m *FilePart) XXX_DiscardUnknown() {
	xxx_messageInfo_FilePart.DiscardUnknown(m)
}

var xxx_messageInfo_FilePart proto.InternalMessageInfo

func (m *FilePart) GetNum() int64 {
	if m != nil {
		return m.Num
	}
	return 0
}

func (m *FilePart) GetOffset() int64 {
	if m != nil {
		return m.Offset
	}
	return 0
}

func (m *FilePart) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func (m *FilePart) GetSize() int64 {
	if m != nil {
		return m.Size
	}
	return 0
}

func (m *FilePart) GetIsCompleted() bool {
	if m != nil {
		return m.IsCompleted
	}
	return false
}

type FileChannel struct {
	Checkpoint           string                     `protobuf:"bytes,1,opt,name=checkpoint,proto3" json:"checkpoint,omitempty"`
	FilePart             *FilePart                  `protobuf:"bytes,2,opt,name=filePart,proto3" json:"filePart,omitempty"`
	FileInfo             *FileInfo                  `protobuf:"bytes,3,opt,name=fileInfo,proto3" json:"fileInfo,omitempty"`
	Strategy             FileChannel_UploadStrategy `protobuf:"varint,4,opt,name=strategy,proto3,enum=pb.FileChannel_UploadStrategy" json:"strategy,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *FileChannel) Reset()         { *m = FileChannel{} }
func (m *FileChannel) String() string { return proto.CompactTextString(m) }
func (*FileChannel) ProtoMessage()    {}
func (*FileChannel) Descriptor() ([]byte, []int) {
	return fileDescriptor_91b94b655bd2a7e5, []int{3}
}

func (m *FileChannel) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileChannel.Unmarshal(m, b)
}
func (m *FileChannel) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileChannel.Marshal(b, m, deterministic)
}
func (m *FileChannel) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileChannel.Merge(m, src)
}
func (m *FileChannel) XXX_Size() int {
	return xxx_messageInfo_FileChannel.Size(m)
}
func (m *FileChannel) XXX_DiscardUnknown() {
	xxx_messageInfo_FileChannel.DiscardUnknown(m)
}

var xxx_messageInfo_FileChannel proto.InternalMessageInfo

func (m *FileChannel) GetCheckpoint() string {
	if m != nil {
		return m.Checkpoint
	}
	return ""
}

func (m *FileChannel) GetFilePart() *FilePart {
	if m != nil {
		return m.FilePart
	}
	return nil
}

func (m *FileChannel) GetFileInfo() *FileInfo {
	if m != nil {
		return m.FileInfo
	}
	return nil
}

func (m *FileChannel) GetStrategy() FileChannel_UploadStrategy {
	if m != nil {
		return m.Strategy
	}
	return FileChannel_EXIST_FAILURE
}

func init() {
	proto.RegisterEnum("pb.FileChannel_UploadStrategy", FileChannel_UploadStrategy_name, FileChannel_UploadStrategy_value)
	proto.RegisterType((*FileInfo)(nil), "pb.FileInfo")
	proto.RegisterType((*CheckPoint)(nil), "pb.CheckPoint")
	proto.RegisterType((*FilePart)(nil), "pb.FilePart")
	proto.RegisterType((*FileChannel)(nil), "pb.FileChannel")
}

func init() {
	proto.RegisterFile("upload.proto", fileDescriptor_91b94b655bd2a7e5)
}

var fileDescriptor_91b94b655bd2a7e5 = []byte{
	// 447 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x92, 0xdf, 0x6e, 0xd3, 0x30,
	0x14, 0xc6, 0xeb, 0xa5, 0x1b, 0xe1, 0xa4, 0xeb, 0x8a, 0x85, 0x50, 0x14, 0xa1, 0xaa, 0xf2, 0x55,
	0xae, 0x32, 0x28, 0x70, 0x03, 0xe2, 0x02, 0x4a, 0x27, 0x15, 0x81, 0x98, 0x5c, 0x36, 0x71, 0x87,
	0xdc, 0xd5, 0x49, 0x23, 0xdc, 0xd8, 0x8a, 0x1d, 0x21, 0xb8, 0xe0, 0xc5, 0x78, 0x2c, 0x5e, 0x00,
	0xd9, 0xf9, 0xb3, 0xb6, 0xa0, 0xdd, 0x9d, 0xf3, 0xf9, 0xd8, 0xdf, 0xf7, 0xb3, 0x0e, 0x0c, 0x2a,
	0x25, 0x24, 0x5b, 0x27, 0xaa, 0x94, 0x46, 0xe2, 0x23, 0xb5, 0x8a, 0xc6, 0x99, 0x94, 0x99, 0xe0,
	0xe7, 0x4e, 0x59, 0x55, 0xe9, 0xf9, 0xf7, 0x92, 0x29, 0xc5, 0x4b, 0x5d, 0xcf, 0x10, 0x03, 0xfe,
	0x45, 0x2e, 0xf8, 0xa2, 0x48, 0x25, 0xc6, 0xd0, 0x57, 0xcc, 0x6c, 0x42, 0x34, 0x41, 0xf1, 0x7d,
	0xea, 0x6a, 0xab, 0xe9, 0xfc, 0x27, 0x0f, 0x8f, 0x26, 0x28, 0xf6, 0xa8, 0xab, 0xf1, 0x43, 0x38,
	0xce, 0xf5, 0xbb, 0xbc, 0x0c, 0xbd, 0x09, 0x8a, 0x7d, 0x5a, 0x37, 0x78, 0x04, 0xde, 0x76, 0xfd,
	0x22, 0xec, 0xbb, 0xcb, 0xb6, 0xc4, 0x21, 0xdc, 0x13, 0x4c, 0x9b, 0x8f, 0x72, 0x1d, 0x1e, 0xbb,
	0xeb, 0x6d, 0x4b, 0xde, 0x03, 0xcc, 0x36, 0xfc, 0xe6, 0xdb, 0xa5, 0xcc, 0x0b, 0xf3, 0x5f, 0xdf,
	0x18, 0xfc, 0xb4, 0xc9, 0xe5, 0xbc, 0x83, 0xe9, 0x20, 0x51, 0xab, 0xa4, 0xcd, 0x4a, 0xbb, 0x53,
	0xf2, 0xab, 0x26, 0xb8, 0x64, 0xa5, 0xb1, 0x19, 0x8a, 0x6a, 0xeb, 0x1e, 0xf2, 0xa8, 0x2d, 0xf1,
	0x23, 0x38, 0x91, 0x69, 0xaa, 0xb9, 0x69, 0x08, 0x9a, 0xce, 0x7a, 0xae, 0x99, 0x61, 0x0e, 0x61,
	0x40, 0x5d, 0xdd, 0xb1, 0xf6, 0x77, 0x58, 0x27, 0x10, 0xe4, 0x7a, 0x26, 0xb7, 0x4a, 0x70, 0xc3,
	0x6b, 0x0e, 0x9f, 0xee, 0x4a, 0xe4, 0x0f, 0x82, 0xc0, 0x06, 0x98, 0x6d, 0x58, 0x51, 0x70, 0x81,
	0xc7, 0x00, 0x37, 0x96, 0x4d, 0x59, 0xb6, 0x86, 0x69, 0x47, 0x69, 0xc9, 0x6c, 0xde, 0x43, 0x32,
	0xab, 0xd1, 0xee, 0x74, 0xef, 0x0f, 0xbc, 0xbb, 0xfe, 0x00, 0xbf, 0x04, 0x5f, 0x9b, 0x92, 0x19,
	0x9e, 0xfd, 0x70, 0xe9, 0x87, 0xd3, 0x71, 0x3b, 0xd9, 0xc4, 0x4a, 0xae, 0xdc, 0x66, 0x2c, 0x9b,
	0x29, 0xda, 0xcd, 0x93, 0xe7, 0x30, 0xdc, 0x3f, 0xc3, 0x0f, 0xe0, 0x74, 0xfe, 0x65, 0xb1, 0xfc,
	0xfc, 0xf5, 0xe2, 0xcd, 0xe2, 0xc3, 0x15, 0x9d, 0x8f, 0x7a, 0xf8, 0x0c, 0x82, 0x5a, 0x9a, 0x7d,
	0xba, 0x9e, 0xd3, 0x11, 0x9a, 0xfe, 0x46, 0x00, 0xf5, 0x35, 0x6b, 0x82, 0x5f, 0x41, 0x90, 0x71,
	0xd3, 0x6d, 0xd2, 0xe3, 0xa4, 0x5e, 0xbb, 0xa4, 0x5d, 0xbb, 0x64, 0x69, 0xca, 0xbc, 0xc8, 0xae,
	0x99, 0xa8, 0x78, 0xb4, 0x47, 0x41, 0x7a, 0xf8, 0x29, 0x9c, 0x6a, 0x6e, 0x76, 0x16, 0x62, 0x68,
	0x07, 0x6e, 0xfb, 0x68, 0xef, 0x83, 0x48, 0xef, 0x09, 0xc2, 0xaf, 0x01, 0xaa, 0x5b, 0xf7, 0xb3,
	0x03, 0xd8, 0x28, 0xfa, 0xc7, 0xff, 0xad, 0x94, 0xc2, 0xb9, 0x93, 0x5e, 0x8c, 0x56, 0x27, 0x4e,
	0x7f, 0xf6, 0x37, 0x00, 0x00, 0xff, 0xff, 0x4a, 0xa4, 0x0c, 0x9c, 0x30, 0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UploadFileClient is the client API for UploadFile service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UploadFileClient interface {
	GetFileInfo(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*FileInfo, error)
	// 获取检查记录文件
	SetCheckPoint(ctx context.Context, in *CheckPoint, opts ...grpc.CallOption) (UploadFile_SetCheckPointClient, error)
	// 上传文件分片
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (UploadFile_UploadFileClient, error)
}

type uploadFileClient struct {
	cc grpc.ClientConnInterface
}

func NewUploadFileClient(cc grpc.ClientConnInterface) UploadFileClient {
	return &uploadFileClient{cc}
}

func (c *uploadFileClient) GetFileInfo(ctx context.Context, in *wrappers.StringValue, opts ...grpc.CallOption) (*FileInfo, error) {
	out := new(FileInfo)
	err := c.cc.Invoke(ctx, "/pb.UploadFile/getFileInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uploadFileClient) SetCheckPoint(ctx context.Context, in *CheckPoint, opts ...grpc.CallOption) (UploadFile_SetCheckPointClient, error) {
	stream, err := c.cc.NewStream(ctx, &_UploadFile_serviceDesc.Streams[0], "/pb.UploadFile/setCheckPoint", opts...)
	if err != nil {
		return nil, err
	}
	x := &uploadFileSetCheckPointClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type UploadFile_SetCheckPointClient interface {
	Recv() (*FilePart, error)
	grpc.ClientStream
}

type uploadFileSetCheckPointClient struct {
	grpc.ClientStream
}

func (x *uploadFileSetCheckPointClient) Recv() (*FilePart, error) {
	m := new(FilePart)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *uploadFileClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (UploadFile_UploadFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &_UploadFile_serviceDesc.Streams[1], "/pb.UploadFile/uploadFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &uploadFileUploadFileClient{stream}
	return x, nil
}

type UploadFile_UploadFileClient interface {
	Send(*FileChannel) error
	CloseAndRecv() (*wrappers.BoolValue, error)
	grpc.ClientStream
}

type uploadFileUploadFileClient struct {
	grpc.ClientStream
}

func (x *uploadFileUploadFileClient) Send(m *FileChannel) error {
	return x.ClientStream.SendMsg(m)
}

func (x *uploadFileUploadFileClient) CloseAndRecv() (*wrappers.BoolValue, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(wrappers.BoolValue)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// UploadFileServer is the server API for UploadFile service.
type UploadFileServer interface {
	GetFileInfo(context.Context, *wrappers.StringValue) (*FileInfo, error)
	// 获取检查记录文件
	SetCheckPoint(*CheckPoint, UploadFile_SetCheckPointServer) error
	// 上传文件分片
	UploadFile(UploadFile_UploadFileServer) error
}

// UnimplementedUploadFileServer can be embedded to have forward compatible implementations.
type UnimplementedUploadFileServer struct {
}

func (*UnimplementedUploadFileServer) GetFileInfo(ctx context.Context, req *wrappers.StringValue) (*FileInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileInfo not implemented")
}
func (*UnimplementedUploadFileServer) SetCheckPoint(req *CheckPoint, srv UploadFile_SetCheckPointServer) error {
	return status.Errorf(codes.Unimplemented, "method SetCheckPoint not implemented")
}
func (*UnimplementedUploadFileServer) UploadFile(srv UploadFile_UploadFileServer) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}

func RegisterUploadFileServer(s *grpc.Server, srv UploadFileServer) {
	s.RegisterService(&_UploadFile_serviceDesc, srv)
}

func _UploadFile_GetFileInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrappers.StringValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UploadFileServer).GetFileInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.UploadFile/GetFileInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UploadFileServer).GetFileInfo(ctx, req.(*wrappers.StringValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _UploadFile_SetCheckPoint_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(CheckPoint)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(UploadFileServer).SetCheckPoint(m, &uploadFileSetCheckPointServer{stream})
}

type UploadFile_SetCheckPointServer interface {
	Send(*FilePart) error
	grpc.ServerStream
}

type uploadFileSetCheckPointServer struct {
	grpc.ServerStream
}

func (x *uploadFileSetCheckPointServer) Send(m *FilePart) error {
	return x.ServerStream.SendMsg(m)
}

func _UploadFile_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(UploadFileServer).UploadFile(&uploadFileUploadFileServer{stream})
}

type UploadFile_UploadFileServer interface {
	SendAndClose(*wrappers.BoolValue) error
	Recv() (*FileChannel, error)
	grpc.ServerStream
}

type uploadFileUploadFileServer struct {
	grpc.ServerStream
}

func (x *uploadFileUploadFileServer) SendAndClose(m *wrappers.BoolValue) error {
	return x.ServerStream.SendMsg(m)
}

func (x *uploadFileUploadFileServer) Recv() (*FileChannel, error) {
	m := new(FileChannel)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _UploadFile_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.UploadFile",
	HandlerType: (*UploadFileServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "getFileInfo",
			Handler:    _UploadFile_GetFileInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "setCheckPoint",
			Handler:       _UploadFile_SetCheckPoint_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "uploadFile",
			Handler:       _UploadFile_UploadFile_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "upload.proto",
}
