// Code generated by protoc-gen-go. DO NOT EDIT.
// source: p2pMessage.proto

package p2p

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

type P2PMessage struct {
	Data                 [][]byte `protobuf:"bytes,1,rep,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *P2PMessage) Reset()         { *m = P2PMessage{} }
func (m *P2PMessage) String() string { return proto.CompactTextString(m) }
func (*P2PMessage) ProtoMessage()    {}
func (*P2PMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a259f8162cba2831, []int{0}
}

func (m *P2PMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_P2PMessage.Unmarshal(m, b)
}
func (m *P2PMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_P2PMessage.Marshal(b, m, deterministic)
}
func (m *P2PMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_P2PMessage.Merge(m, src)
}
func (m *P2PMessage) XXX_Size() int {
	return xxx_messageInfo_P2PMessage.Size(m)
}
func (m *P2PMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_P2PMessage.DiscardUnknown(m)
}

var xxx_messageInfo_P2PMessage proto.InternalMessageInfo

func (m *P2PMessage) GetData() [][]byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type BusMessage struct {
	Command              string   `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	Parameters           [][]byte `protobuf:"bytes,2,rep,name=Parameters,proto3" json:"Parameters,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *BusMessage) Reset()         { *m = BusMessage{} }
func (m *BusMessage) String() string { return proto.CompactTextString(m) }
func (*BusMessage) ProtoMessage()    {}
func (*BusMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a259f8162cba2831, []int{1}
}

func (m *BusMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BusMessage.Unmarshal(m, b)
}
func (m *BusMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BusMessage.Marshal(b, m, deterministic)
}
func (m *BusMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BusMessage.Merge(m, src)
}
func (m *BusMessage) XXX_Size() int {
	return xxx_messageInfo_BusMessage.Size(m)
}
func (m *BusMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_BusMessage.DiscardUnknown(m)
}

var xxx_messageInfo_BusMessage proto.InternalMessageInfo

func (m *BusMessage) GetCommand() string {
	if m != nil {
		return m.Command
	}
	return ""
}

func (m *BusMessage) GetParameters() [][]byte {
	if m != nil {
		return m.Parameters
	}
	return nil
}

// to parse Listeners parameter from announce module
type Addrs struct {
	Address              [][]byte `protobuf:"bytes,1,rep,name=Address,proto3" json:"Address,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Addrs) Reset()         { *m = Addrs{} }
func (m *Addrs) String() string { return proto.CompactTextString(m) }
func (*Addrs) ProtoMessage()    {}
func (*Addrs) Descriptor() ([]byte, []int) {
	return fileDescriptor_a259f8162cba2831, []int{2}
}

func (m *Addrs) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Addrs.Unmarshal(m, b)
}
func (m *Addrs) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Addrs.Marshal(b, m, deterministic)
}
func (m *Addrs) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Addrs.Merge(m, src)
}
func (m *Addrs) XXX_Size() int {
	return xxx_messageInfo_Addrs.Size(m)
}
func (m *Addrs) XXX_DiscardUnknown() {
	xxx_messageInfo_Addrs.DiscardUnknown(m)
}

var xxx_messageInfo_Addrs proto.InternalMessageInfo

func (m *Addrs) GetAddress() [][]byte {
	if m != nil {
		return m.Address
	}
	return nil
}

type MockProtoMessage struct {
	Data                 string   `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MockProtoMessage) Reset()         { *m = MockProtoMessage{} }
func (m *MockProtoMessage) String() string { return proto.CompactTextString(m) }
func (*MockProtoMessage) ProtoMessage()    {}
func (*MockProtoMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_a259f8162cba2831, []int{3}
}

func (m *MockProtoMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MockProtoMessage.Unmarshal(m, b)
}
func (m *MockProtoMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MockProtoMessage.Marshal(b, m, deterministic)
}
func (m *MockProtoMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MockProtoMessage.Merge(m, src)
}
func (m *MockProtoMessage) XXX_Size() int {
	return xxx_messageInfo_MockProtoMessage.Size(m)
}
func (m *MockProtoMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_MockProtoMessage.DiscardUnknown(m)
}

var xxx_messageInfo_MockProtoMessage proto.InternalMessageInfo

func (m *MockProtoMessage) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*P2PMessage)(nil), "bitmark.bitmarkd.p2p.P2PMessage")
	proto.RegisterType((*BusMessage)(nil), "bitmark.bitmarkd.p2p.BusMessage")
	proto.RegisterType((*Addrs)(nil), "bitmark.bitmarkd.p2p.Addrs")
	proto.RegisterType((*MockProtoMessage)(nil), "bitmark.bitmarkd.p2p.MockProtoMessage")
}

func init() { proto.RegisterFile("p2pMessage.proto", fileDescriptor_a259f8162cba2831) }

var fileDescriptor_a259f8162cba2831 = []byte{
	// 182 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x64, 0x8f, 0x31, 0x0b, 0xc2, 0x30,
	0x14, 0x84, 0xa9, 0xb5, 0x8a, 0x0f, 0x87, 0x12, 0x1c, 0x32, 0x49, 0xcd, 0x20, 0x9d, 0x3a, 0xc4,
	0x5f, 0x60, 0x11, 0xb7, 0x42, 0xe8, 0xe8, 0x96, 0x36, 0x41, 0xa4, 0xc4, 0x84, 0xbc, 0xf8, 0xff,
	0x25, 0xa1, 0x05, 0xc1, 0x29, 0xdf, 0x85, 0xe3, 0xee, 0x1e, 0x94, 0x8e, 0xbb, 0x4e, 0x23, 0xca,
	0xa7, 0x6e, 0x9c, 0xb7, 0xc1, 0x92, 0xc3, 0xf0, 0x0a, 0x46, 0xfa, 0xa9, 0x99, 0x5f, 0xd5, 0x38,
	0xee, 0x58, 0x05, 0x20, 0xb8, 0x98, 0x9d, 0x84, 0xc0, 0xfa, 0x26, 0x83, 0xa4, 0x59, 0x95, 0xd7,
	0xfb, 0x3e, 0x31, 0xbb, 0x03, 0xb4, 0x1f, 0x5c, 0x1c, 0x14, 0xb6, 0xa3, 0x35, 0x46, 0xbe, 0x15,
	0xcd, 0xaa, 0xac, 0xde, 0xf5, 0x8b, 0x24, 0x47, 0x00, 0x21, 0xbd, 0x34, 0x3a, 0x68, 0x8f, 0x74,
	0x95, 0x12, 0x7e, 0x7e, 0xd8, 0x09, 0x8a, 0xab, 0x52, 0x1e, 0x63, 0x44, 0x04, 0x8d, 0x38, 0xf7,
	0x2c, 0x92, 0x9d, 0xa1, 0xec, 0xec, 0x38, 0x89, 0xb8, 0xf7, 0x7f, 0x52, 0x6c, 0x4b, 0xdc, 0x16,
	0x8f, 0xdc, 0x71, 0x37, 0x6c, 0xd2, 0x61, 0x97, 0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc1, 0xf2,
	0xad, 0xda, 0xec, 0x00, 0x00, 0x00,
}
