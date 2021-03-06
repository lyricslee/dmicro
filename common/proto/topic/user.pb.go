// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common/proto/topic/user.proto

package dmicro_topic

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type UserInfo struct {
	Uid                  int64    `protobuf:"varint,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Mobile               string   `protobuf:"bytes,2,opt,name=mobile,proto3" json:"mobile,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UserInfo) Reset()         { *m = UserInfo{} }
func (m *UserInfo) String() string { return proto.CompactTextString(m) }
func (*UserInfo) ProtoMessage()    {}
func (*UserInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_85bec14ef9159631, []int{0}
}

func (m *UserInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserInfo.Unmarshal(m, b)
}
func (m *UserInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserInfo.Marshal(b, m, deterministic)
}
func (m *UserInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserInfo.Merge(m, src)
}
func (m *UserInfo) XXX_Size() int {
	return xxx_messageInfo_UserInfo.Size(m)
}
func (m *UserInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_UserInfo.DiscardUnknown(m)
}

var xxx_messageInfo_UserInfo proto.InternalMessageInfo

func (m *UserInfo) GetUid() int64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *UserInfo) GetMobile() string {
	if m != nil {
		return m.Mobile
	}
	return ""
}

type UserCreated struct {
	Id                   int64     `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Topic                string    `protobuf:"bytes,2,opt,name=topic,proto3" json:"topic,omitempty"`
	Info                 *UserInfo `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *UserCreated) Reset()         { *m = UserCreated{} }
func (m *UserCreated) String() string { return proto.CompactTextString(m) }
func (*UserCreated) ProtoMessage()    {}
func (*UserCreated) Descriptor() ([]byte, []int) {
	return fileDescriptor_85bec14ef9159631, []int{1}
}

func (m *UserCreated) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UserCreated.Unmarshal(m, b)
}
func (m *UserCreated) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UserCreated.Marshal(b, m, deterministic)
}
func (m *UserCreated) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserCreated.Merge(m, src)
}
func (m *UserCreated) XXX_Size() int {
	return xxx_messageInfo_UserCreated.Size(m)
}
func (m *UserCreated) XXX_DiscardUnknown() {
	xxx_messageInfo_UserCreated.DiscardUnknown(m)
}

var xxx_messageInfo_UserCreated proto.InternalMessageInfo

func (m *UserCreated) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UserCreated) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *UserCreated) GetInfo() *UserInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func init() {
	proto.RegisterType((*UserInfo)(nil), "dmicro.topic.UserInfo")
	proto.RegisterType((*UserCreated)(nil), "dmicro.topic.UserCreated")
}

func init() { proto.RegisterFile("common/proto/topic/user.proto", fileDescriptor_85bec14ef9159631) }

var fileDescriptor_85bec14ef9159631 = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4d, 0xce, 0xcf, 0xcd,
	0xcd, 0xcf, 0xd3, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0xd7, 0x2f, 0xc9, 0x2f, 0xc8, 0x4c, 0xd6, 0x2f,
	0x2d, 0x4e, 0x2d, 0xd2, 0x03, 0x0b, 0x08, 0xf1, 0xa4, 0xe4, 0x66, 0x26, 0x17, 0xe5, 0xeb, 0x81,
	0x25, 0x94, 0x4c, 0xb8, 0x38, 0x42, 0x8b, 0x53, 0x8b, 0x3c, 0xf3, 0xd2, 0xf2, 0x85, 0x04, 0xb8,
	0x98, 0x4b, 0x33, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x98, 0x83, 0x40, 0x4c, 0x21, 0x31, 0x2e,
	0xb6, 0xdc, 0xfc, 0xa4, 0xcc, 0x9c, 0x54, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x28, 0x4f,
	0x29, 0x9e, 0x8b, 0x1b, 0xa4, 0xcb, 0xb9, 0x28, 0x35, 0xb1, 0x24, 0x35, 0x45, 0x88, 0x8f, 0x8b,
	0x09, 0xae, 0x8f, 0x29, 0x33, 0x45, 0x48, 0x84, 0x8b, 0x15, 0x6c, 0x3a, 0x54, 0x17, 0x84, 0x23,
	0xa4, 0xc5, 0xc5, 0x92, 0x99, 0x97, 0x96, 0x2f, 0xc1, 0xac, 0xc0, 0xa8, 0xc1, 0x6d, 0x24, 0xa6,
	0x87, 0xec, 0x0e, 0x3d, 0x98, 0x23, 0x82, 0xc0, 0x6a, 0x92, 0xd8, 0xc0, 0x6e, 0x35, 0x06, 0x04,
	0x00, 0x00, 0xff, 0xff, 0xc0, 0xa8, 0xbd, 0x65, 0xcc, 0x00, 0x00, 0x00,
}
