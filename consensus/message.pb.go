// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: message.proto

package consensus

import (
	bytes "bytes"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Message defines the data needed by spos to communicate between nodes over network in all subrounds
type Message struct {
	HeaderHash          []byte `protobuf:"bytes,1,opt,name=HeaderHash,proto3" json:"HeaderHash,omitempty"`
	SignatureShare      []byte `protobuf:"bytes,2,opt,name=SignatureShare,proto3" json:"SignatureShare,omitempty"`
	Body                []byte `protobuf:"bytes,3,opt,name=Body,proto3" json:"Body,omitempty"`
	Header              []byte `protobuf:"bytes,4,opt,name=Header,proto3" json:"Header,omitempty"`
	PubKey              []byte `protobuf:"bytes,5,opt,name=PubKey,proto3" json:"PubKey,omitempty"`
	Signature           []byte `protobuf:"bytes,6,opt,name=Signature,proto3" json:"Signature,omitempty"`
	MsgType             int64  `protobuf:"varint,7,opt,name=MsgType,proto3" json:"MsgType,omitempty"`
	RoundIndex          int64  `protobuf:"varint,8,opt,name=RoundIndex,proto3" json:"RoundIndex,omitempty"`
	ChainID             []byte `protobuf:"bytes,9,opt,name=ChainID,proto3" json:"ChainID,omitempty"`
	PubKeysBitmap       []byte `protobuf:"bytes,10,opt,name=PubKeysBitmap,proto3" json:"PubKeysBitmap,omitempty"`
	AggregateSignature  []byte `protobuf:"bytes,11,opt,name=AggregateSignature,proto3" json:"AggregateSignature,omitempty"`
	LeaderSignature     []byte `protobuf:"bytes,12,opt,name=LeaderSignature,proto3" json:"LeaderSignature,omitempty"`
	OriginatorPid       []byte `protobuf:"bytes,13,opt,name=OriginatorPid,proto3" json:"OriginatorPid,omitempty"`
	InvalidSigners      []byte `protobuf:"bytes,14,opt,name=InvalidSigners,proto3" json:"InvalidSigners,omitempty"`
	ProcessedHeaderHash []byte `protobuf:"bytes,15,opt,name=ProcessedHeaderHash,proto3" json:"ProcessedHeaderHash,omitempty"`
}

func (m *Message) Reset()      { *m = Message{} }
func (*Message) ProtoMessage() {}
func (*Message) Descriptor() ([]byte, []int) {
	return fileDescriptor_33c57e4bae7b9afd, []int{0}
}
func (m *Message) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Message) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Message) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Message.Merge(m, src)
}
func (m *Message) XXX_Size() int {
	return m.Size()
}
func (m *Message) XXX_DiscardUnknown() {
	xxx_messageInfo_Message.DiscardUnknown(m)
}

var xxx_messageInfo_Message proto.InternalMessageInfo

func (m *Message) GetHeaderHash() []byte {
	if m != nil {
		return m.HeaderHash
	}
	return nil
}

func (m *Message) GetSignatureShare() []byte {
	if m != nil {
		return m.SignatureShare
	}
	return nil
}

func (m *Message) GetBody() []byte {
	if m != nil {
		return m.Body
	}
	return nil
}

func (m *Message) GetHeader() []byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Message) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *Message) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func (m *Message) GetMsgType() int64 {
	if m != nil {
		return m.MsgType
	}
	return 0
}

func (m *Message) GetRoundIndex() int64 {
	if m != nil {
		return m.RoundIndex
	}
	return 0
}

func (m *Message) GetChainID() []byte {
	if m != nil {
		return m.ChainID
	}
	return nil
}

func (m *Message) GetPubKeysBitmap() []byte {
	if m != nil {
		return m.PubKeysBitmap
	}
	return nil
}

func (m *Message) GetAggregateSignature() []byte {
	if m != nil {
		return m.AggregateSignature
	}
	return nil
}

func (m *Message) GetLeaderSignature() []byte {
	if m != nil {
		return m.LeaderSignature
	}
	return nil
}

func (m *Message) GetOriginatorPid() []byte {
	if m != nil {
		return m.OriginatorPid
	}
	return nil
}

func (m *Message) GetInvalidSigners() []byte {
	if m != nil {
		return m.InvalidSigners
	}
	return nil
}

func (m *Message) GetProcessedHeaderHash() []byte {
	if m != nil {
		return m.ProcessedHeaderHash
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "proto.Message")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor_33c57e4bae7b9afd) }

var fileDescriptor_33c57e4bae7b9afd = []byte{
	// 402 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcb, 0xee, 0x12, 0x31,
	0x14, 0xc6, 0xa7, 0x72, 0x93, 0xca, 0x25, 0xa9, 0x89, 0x69, 0x8c, 0x69, 0x88, 0x31, 0x86, 0x8d,
	0x60, 0xe2, 0x13, 0x08, 0x2e, 0x20, 0x4a, 0x24, 0xe0, 0xca, 0x5d, 0x87, 0xa9, 0x9d, 0x26, 0xd2,
	0x92, 0x76, 0xc6, 0xc8, 0xce, 0x9d, 0x5b, 0x1f, 0xc3, 0x47, 0x71, 0xc9, 0x92, 0xa5, 0x94, 0x8d,
	0x4b, 0x1e, 0xc1, 0xcc, 0xe9, 0x9f, 0x3f, 0x97, 0xb0, 0x9a, 0x7e, 0xbf, 0xef, 0x3b, 0xe7, 0x4c,
	0x73, 0x8a, 0x9b, 0x4b, 0xe1, 0x1c, 0x97, 0xa2, 0xb7, 0xb2, 0x26, 0x33, 0xa4, 0x02, 0x9f, 0xa7,
	0xaf, 0xa4, 0xca, 0xd2, 0x3c, 0xee, 0x2d, 0xcc, 0xb2, 0x2f, 0x8d, 0x34, 0x7d, 0xc0, 0x71, 0xfe,
	0x05, 0x14, 0x08, 0x38, 0x85, 0xaa, 0xe7, 0x3f, 0xcb, 0xb8, 0x36, 0x09, 0x7d, 0x08, 0xc3, 0x78,
	0x24, 0x78, 0x22, 0xec, 0x88, 0xbb, 0x94, 0xa2, 0x0e, 0xea, 0x36, 0x66, 0x67, 0x84, 0xbc, 0xc4,
	0xad, 0xb9, 0x92, 0x9a, 0x67, 0xb9, 0x15, 0xf3, 0x94, 0x5b, 0x41, 0x1f, 0x40, 0xe6, 0x8a, 0x12,
	0x82, 0xcb, 0x03, 0x93, 0xac, 0x69, 0x09, 0x5c, 0x38, 0x93, 0x27, 0xb8, 0x1a, 0x3a, 0xd1, 0x32,
	0xd0, 0x3b, 0x55, 0xf0, 0x69, 0x1e, 0xbf, 0x17, 0x6b, 0x5a, 0x09, 0x3c, 0x28, 0xf2, 0x0c, 0xd7,
	0xef, 0xbb, 0xd2, 0x2a, 0x58, 0x27, 0x40, 0x28, 0xae, 0x4d, 0x9c, 0xfc, 0xb4, 0x5e, 0x09, 0x5a,
	0xeb, 0xa0, 0x6e, 0x69, 0x76, 0x94, 0xc5, 0x1d, 0x66, 0x26, 0xd7, 0xc9, 0x58, 0x27, 0xe2, 0x3b,
	0x7d, 0x08, 0xe6, 0x19, 0x29, 0x2a, 0x87, 0x29, 0x57, 0x7a, 0xfc, 0x8e, 0xd6, 0xa1, 0xeb, 0x51,
	0x92, 0x17, 0xb8, 0x19, 0x66, 0xbb, 0x81, 0xca, 0x96, 0x7c, 0x45, 0x31, 0xf8, 0x97, 0x90, 0xf4,
	0x30, 0x79, 0x2b, 0xa5, 0x15, 0x92, 0x67, 0xe2, 0xf4, 0x83, 0x8f, 0x20, 0x7a, 0xc3, 0x21, 0x5d,
	0xdc, 0xfe, 0x00, 0x37, 0x3d, 0x85, 0x1b, 0x10, 0xbe, 0xc6, 0xc5, 0xfc, 0x8f, 0x56, 0x49, 0xa5,
	0x79, 0x66, 0xec, 0x54, 0x25, 0xb4, 0x19, 0xe6, 0x5f, 0xc0, 0x62, 0x07, 0x63, 0xfd, 0x8d, 0x7f,
	0x55, 0x49, 0x51, 0x29, 0xac, 0xa3, 0xad, 0xb0, 0x83, 0x4b, 0x4a, 0x5e, 0xe3, 0xc7, 0x53, 0x6b,
	0x16, 0xc2, 0x39, 0x91, 0x9c, 0x2d, 0xb5, 0x0d, 0xe1, 0x5b, 0xd6, 0x60, 0xb8, 0xd9, 0xb1, 0x68,
	0xbb, 0x63, 0xd1, 0x61, 0xc7, 0xd0, 0x0f, 0xcf, 0xd0, 0x6f, 0xcf, 0xd0, 0x1f, 0xcf, 0xd0, 0xc6,
	0x33, 0xb4, 0xf5, 0x0c, 0xfd, 0xf5, 0x0c, 0xfd, 0xf3, 0x2c, 0x3a, 0x78, 0x86, 0x7e, 0xed, 0x59,
	0xb4, 0xd9, 0xb3, 0x68, 0xbb, 0x67, 0xd1, 0xe7, 0xfa, 0xc2, 0x68, 0x27, 0xb4, 0xcb, 0x5d, 0x5c,
	0x85, 0x57, 0xf5, 0xe6, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x1f, 0x85, 0xc6, 0x21, 0x9c, 0x02,
	0x00, 0x00,
}

func (this *Message) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Message)
	if !ok {
		that2, ok := that.(Message)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !bytes.Equal(this.HeaderHash, that1.HeaderHash) {
		return false
	}
	if !bytes.Equal(this.SignatureShare, that1.SignatureShare) {
		return false
	}
	if !bytes.Equal(this.Body, that1.Body) {
		return false
	}
	if !bytes.Equal(this.Header, that1.Header) {
		return false
	}
	if !bytes.Equal(this.PubKey, that1.PubKey) {
		return false
	}
	if !bytes.Equal(this.Signature, that1.Signature) {
		return false
	}
	if this.MsgType != that1.MsgType {
		return false
	}
	if this.RoundIndex != that1.RoundIndex {
		return false
	}
	if !bytes.Equal(this.ChainID, that1.ChainID) {
		return false
	}
	if !bytes.Equal(this.PubKeysBitmap, that1.PubKeysBitmap) {
		return false
	}
	if !bytes.Equal(this.AggregateSignature, that1.AggregateSignature) {
		return false
	}
	if !bytes.Equal(this.LeaderSignature, that1.LeaderSignature) {
		return false
	}
	if !bytes.Equal(this.OriginatorPid, that1.OriginatorPid) {
		return false
	}
	if !bytes.Equal(this.InvalidSigners, that1.InvalidSigners) {
		return false
	}
	if !bytes.Equal(this.ProcessedHeaderHash, that1.ProcessedHeaderHash) {
		return false
	}
	return true
}
func (this *Message) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 19)
	s = append(s, "&consensus.Message{")
	s = append(s, "HeaderHash: "+fmt.Sprintf("%#v", this.HeaderHash)+",\n")
	s = append(s, "SignatureShare: "+fmt.Sprintf("%#v", this.SignatureShare)+",\n")
	s = append(s, "Body: "+fmt.Sprintf("%#v", this.Body)+",\n")
	s = append(s, "Header: "+fmt.Sprintf("%#v", this.Header)+",\n")
	s = append(s, "PubKey: "+fmt.Sprintf("%#v", this.PubKey)+",\n")
	s = append(s, "Signature: "+fmt.Sprintf("%#v", this.Signature)+",\n")
	s = append(s, "MsgType: "+fmt.Sprintf("%#v", this.MsgType)+",\n")
	s = append(s, "RoundIndex: "+fmt.Sprintf("%#v", this.RoundIndex)+",\n")
	s = append(s, "ChainID: "+fmt.Sprintf("%#v", this.ChainID)+",\n")
	s = append(s, "PubKeysBitmap: "+fmt.Sprintf("%#v", this.PubKeysBitmap)+",\n")
	s = append(s, "AggregateSignature: "+fmt.Sprintf("%#v", this.AggregateSignature)+",\n")
	s = append(s, "LeaderSignature: "+fmt.Sprintf("%#v", this.LeaderSignature)+",\n")
	s = append(s, "OriginatorPid: "+fmt.Sprintf("%#v", this.OriginatorPid)+",\n")
	s = append(s, "InvalidSigners: "+fmt.Sprintf("%#v", this.InvalidSigners)+",\n")
	s = append(s, "ProcessedHeaderHash: "+fmt.Sprintf("%#v", this.ProcessedHeaderHash)+",\n")
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringMessage(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}
func (m *Message) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Message) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Message) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ProcessedHeaderHash) > 0 {
		i -= len(m.ProcessedHeaderHash)
		copy(dAtA[i:], m.ProcessedHeaderHash)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.ProcessedHeaderHash)))
		i--
		dAtA[i] = 0x7a
	}
	if len(m.InvalidSigners) > 0 {
		i -= len(m.InvalidSigners)
		copy(dAtA[i:], m.InvalidSigners)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.InvalidSigners)))
		i--
		dAtA[i] = 0x72
	}
	if len(m.OriginatorPid) > 0 {
		i -= len(m.OriginatorPid)
		copy(dAtA[i:], m.OriginatorPid)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.OriginatorPid)))
		i--
		dAtA[i] = 0x6a
	}
	if len(m.LeaderSignature) > 0 {
		i -= len(m.LeaderSignature)
		copy(dAtA[i:], m.LeaderSignature)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.LeaderSignature)))
		i--
		dAtA[i] = 0x62
	}
	if len(m.AggregateSignature) > 0 {
		i -= len(m.AggregateSignature)
		copy(dAtA[i:], m.AggregateSignature)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.AggregateSignature)))
		i--
		dAtA[i] = 0x5a
	}
	if len(m.PubKeysBitmap) > 0 {
		i -= len(m.PubKeysBitmap)
		copy(dAtA[i:], m.PubKeysBitmap)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.PubKeysBitmap)))
		i--
		dAtA[i] = 0x52
	}
	if len(m.ChainID) > 0 {
		i -= len(m.ChainID)
		copy(dAtA[i:], m.ChainID)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.ChainID)))
		i--
		dAtA[i] = 0x4a
	}
	if m.RoundIndex != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.RoundIndex))
		i--
		dAtA[i] = 0x40
	}
	if m.MsgType != 0 {
		i = encodeVarintMessage(dAtA, i, uint64(m.MsgType))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Header) > 0 {
		i -= len(m.Header)
		copy(dAtA[i:], m.Header)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Header)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Body) > 0 {
		i -= len(m.Body)
		copy(dAtA[i:], m.Body)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.Body)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SignatureShare) > 0 {
		i -= len(m.SignatureShare)
		copy(dAtA[i:], m.SignatureShare)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.SignatureShare)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.HeaderHash) > 0 {
		i -= len(m.HeaderHash)
		copy(dAtA[i:], m.HeaderHash)
		i = encodeVarintMessage(dAtA, i, uint64(len(m.HeaderHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMessage(dAtA []byte, offset int, v uint64) int {
	offset -= sovMessage(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Message) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.HeaderHash)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.SignatureShare)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.Body)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.Header)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	if m.MsgType != 0 {
		n += 1 + sovMessage(uint64(m.MsgType))
	}
	if m.RoundIndex != 0 {
		n += 1 + sovMessage(uint64(m.RoundIndex))
	}
	l = len(m.ChainID)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.PubKeysBitmap)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.AggregateSignature)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.LeaderSignature)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.OriginatorPid)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.InvalidSigners)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	l = len(m.ProcessedHeaderHash)
	if l > 0 {
		n += 1 + l + sovMessage(uint64(l))
	}
	return n
}

func sovMessage(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMessage(x uint64) (n int) {
	return sovMessage(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Message) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Message{`,
		`HeaderHash:` + fmt.Sprintf("%v", this.HeaderHash) + `,`,
		`SignatureShare:` + fmt.Sprintf("%v", this.SignatureShare) + `,`,
		`Body:` + fmt.Sprintf("%v", this.Body) + `,`,
		`Header:` + fmt.Sprintf("%v", this.Header) + `,`,
		`PubKey:` + fmt.Sprintf("%v", this.PubKey) + `,`,
		`Signature:` + fmt.Sprintf("%v", this.Signature) + `,`,
		`MsgType:` + fmt.Sprintf("%v", this.MsgType) + `,`,
		`RoundIndex:` + fmt.Sprintf("%v", this.RoundIndex) + `,`,
		`ChainID:` + fmt.Sprintf("%v", this.ChainID) + `,`,
		`PubKeysBitmap:` + fmt.Sprintf("%v", this.PubKeysBitmap) + `,`,
		`AggregateSignature:` + fmt.Sprintf("%v", this.AggregateSignature) + `,`,
		`LeaderSignature:` + fmt.Sprintf("%v", this.LeaderSignature) + `,`,
		`OriginatorPid:` + fmt.Sprintf("%v", this.OriginatorPid) + `,`,
		`InvalidSigners:` + fmt.Sprintf("%v", this.InvalidSigners) + `,`,
		`ProcessedHeaderHash:` + fmt.Sprintf("%v", this.ProcessedHeaderHash) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringMessage(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Message) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Message: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Message: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HeaderHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.HeaderHash = append(m.HeaderHash[:0], dAtA[iNdEx:postIndex]...)
			if m.HeaderHash == nil {
				m.HeaderHash = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignatureShare", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignatureShare = append(m.SignatureShare[:0], dAtA[iNdEx:postIndex]...)
			if m.SignatureShare == nil {
				m.SignatureShare = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Body", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Body = append(m.Body[:0], dAtA[iNdEx:postIndex]...)
			if m.Body == nil {
				m.Body = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Header", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Header = append(m.Header[:0], dAtA[iNdEx:postIndex]...)
			if m.Header == nil {
				m.Header = []byte{}
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgType", wireType)
			}
			m.MsgType = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MsgType |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field RoundIndex", wireType)
			}
			m.RoundIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.RoundIndex |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChainID", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChainID = append(m.ChainID[:0], dAtA[iNdEx:postIndex]...)
			if m.ChainID == nil {
				m.ChainID = []byte{}
			}
			iNdEx = postIndex
		case 10:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKeysBitmap", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKeysBitmap = append(m.PubKeysBitmap[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKeysBitmap == nil {
				m.PubKeysBitmap = []byte{}
			}
			iNdEx = postIndex
		case 11:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AggregateSignature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AggregateSignature = append(m.AggregateSignature[:0], dAtA[iNdEx:postIndex]...)
			if m.AggregateSignature == nil {
				m.AggregateSignature = []byte{}
			}
			iNdEx = postIndex
		case 12:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field LeaderSignature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.LeaderSignature = append(m.LeaderSignature[:0], dAtA[iNdEx:postIndex]...)
			if m.LeaderSignature == nil {
				m.LeaderSignature = []byte{}
			}
			iNdEx = postIndex
		case 13:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginatorPid", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OriginatorPid = append(m.OriginatorPid[:0], dAtA[iNdEx:postIndex]...)
			if m.OriginatorPid == nil {
				m.OriginatorPid = []byte{}
			}
			iNdEx = postIndex
		case 14:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field InvalidSigners", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.InvalidSigners = append(m.InvalidSigners[:0], dAtA[iNdEx:postIndex]...)
			if m.InvalidSigners == nil {
				m.InvalidSigners = []byte{}
			}
			iNdEx = postIndex
		case 15:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProcessedHeaderHash", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMessage
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMessage
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ProcessedHeaderHash = append(m.ProcessedHeaderHash[:0], dAtA[iNdEx:postIndex]...)
			if m.ProcessedHeaderHash == nil {
				m.ProcessedHeaderHash = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMessage(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthMessage
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipMessage(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMessage
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMessage
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthMessage
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMessage
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMessage
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMessage        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMessage          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMessage = fmt.Errorf("proto: unexpected end of group")
)
