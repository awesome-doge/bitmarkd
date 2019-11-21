// Code generated by protoc-gen-go. DO NOT EDIT.
// source: peerstore.proto

package announce

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

type PeerList struct {
	Peers                []*PeerItem `protobuf:"bytes,1,rep,name=Peers,proto3" json:"Peers,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *PeerList) Reset()         { *m = PeerList{} }
func (m *PeerList) String() string { return proto.CompactTextString(m) }
func (*PeerList) ProtoMessage()    {}
func (*PeerList) Descriptor() ([]byte, []int) {
	return fileDescriptor_37da6c6f39403d68, []int{0}
}

func (m *PeerList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerList.Unmarshal(m, b)
}
func (m *PeerList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerList.Marshal(b, m, deterministic)
}
func (m *PeerList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerList.Merge(m, src)
}
func (m *PeerList) XXX_Size() int {
	return xxx_messageInfo_PeerList.Size(m)
}
func (m *PeerList) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerList.DiscardUnknown(m)
}

var xxx_messageInfo_PeerList proto.InternalMessageInfo

func (m *PeerList) GetPeers() []*PeerItem {
	if m != nil {
		return m.Peers
	}
	return nil
}

type PeerItem struct {
	PeerID               []byte   `protobuf:"bytes,1,opt,name=PeerID,proto3" json:"PeerID,omitempty"`
	Listeners            *Addrs   `protobuf:"bytes,2,opt,name=Listeners,proto3" json:"Listeners,omitempty"`
	Timestamp            uint64   `protobuf:"varint,3,opt,name=Timestamp,proto3" json:"Timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PeerItem) Reset()         { *m = PeerItem{} }
func (m *PeerItem) String() string { return proto.CompactTextString(m) }
func (*PeerItem) ProtoMessage()    {}
func (*PeerItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_37da6c6f39403d68, []int{1}
}

func (m *PeerItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PeerItem.Unmarshal(m, b)
}
func (m *PeerItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PeerItem.Marshal(b, m, deterministic)
}
func (m *PeerItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PeerItem.Merge(m, src)
}
func (m *PeerItem) XXX_Size() int {
	return xxx_messageInfo_PeerItem.Size(m)
}
func (m *PeerItem) XXX_DiscardUnknown() {
	xxx_messageInfo_PeerItem.DiscardUnknown(m)
}

var xxx_messageInfo_PeerItem proto.InternalMessageInfo

func (m *PeerItem) GetPeerID() []byte {
	if m != nil {
		return m.PeerID
	}
	return nil
}

func (m *PeerItem) GetListeners() *Addrs {
	if m != nil {
		return m.Listeners
	}
	return nil
}

func (m *PeerItem) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

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
	return fileDescriptor_37da6c6f39403d68, []int{2}
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

func init() {
	proto.RegisterType((*PeerList)(nil), "bitmark.bitmarkd.announce.PeerList")
	proto.RegisterType((*PeerItem)(nil), "bitmark.bitmarkd.announce.PeerItem")
	proto.RegisterType((*Addrs)(nil), "bitmark.bitmarkd.announce.Addrs")
}

func init() { proto.RegisterFile("peerstore.proto", fileDescriptor_37da6c6f39403d68) }

var fileDescriptor_37da6c6f39403d68 = []byte{
	// 206 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0x48, 0x4d, 0x2d,
	0x2a, 0x2e, 0xc9, 0x2f, 0x4a, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x92, 0x4c, 0xca, 0x2c,
	0xc9, 0x4d, 0x2c, 0xca, 0xd6, 0x83, 0xd2, 0x29, 0x7a, 0x89, 0x79, 0x79, 0xf9, 0xa5, 0x79, 0xc9,
	0xa9, 0x4a, 0xae, 0x5c, 0x1c, 0x01, 0xa9, 0xa9, 0x45, 0x3e, 0x99, 0xc5, 0x25, 0x42, 0x96, 0x5c,
	0xac, 0x20, 0x76, 0xb1, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0xb7, 0x91, 0xb2, 0x1e, 0x4e, 0x6d, 0x7a,
	0x20, 0x75, 0x9e, 0x25, 0xa9, 0xb9, 0x41, 0x10, 0x1d, 0x4a, 0x0d, 0x8c, 0x10, 0x73, 0x40, 0x62,
	0x42, 0x62, 0x5c, 0x6c, 0x60, 0xb6, 0x8b, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x4f, 0x10, 0x94, 0x27,
	0x64, 0xc7, 0xc5, 0x09, 0xb2, 0x27, 0x35, 0x0f, 0x64, 0x07, 0x93, 0x02, 0xa3, 0x06, 0xb7, 0x91,
	0x02, 0x1e, 0x3b, 0x1c, 0x53, 0x52, 0x8a, 0x8a, 0x83, 0x10, 0x5a, 0x84, 0x64, 0xb8, 0x38, 0x43,
	0x32, 0x73, 0x53, 0x8b, 0x4b, 0x12, 0x73, 0x0b, 0x24, 0x98, 0x15, 0x18, 0x35, 0x58, 0x82, 0x10,
	0x02, 0x4a, 0x8a, 0x5c, 0xac, 0x60, 0x1d, 0x42, 0x12, 0x5c, 0xec, 0x20, 0x46, 0x6a, 0x31, 0xc4,
	0x23, 0x3c, 0x41, 0x30, 0xae, 0x13, 0x57, 0x14, 0x07, 0xcc, 0xf4, 0x24, 0x36, 0x70, 0xd0, 0x18,
	0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xb3, 0x10, 0x4f, 0xb7, 0x2d, 0x01, 0x00, 0x00,
}
