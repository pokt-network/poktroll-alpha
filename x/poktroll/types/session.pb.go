// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: poktroll/poktroll/session.proto

package types

import (
	fmt "fmt"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
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

type Session struct {
	Id               string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	SessionNumber    int64  `protobuf:"varint,2,opt,name=session_number,json=sessionNumber,proto3" json:"session_number,omitempty"`
	SessionHeight    int64  `protobuf:"varint,3,opt,name=session_height,json=sessionHeight,proto3" json:"session_height,omitempty"`
	NumSessionBlocks int64  `protobuf:"varint,4,opt,name=num_session_blocks,json=numSessionBlocks,proto3" json:"num_session_blocks,omitempty"`
	// CONSIDERATION: Should we add a `RelayChain` enum and use it across the board?
	// CONSIDERATION: Should a single session support multiple relay chains?
	// TECHDEBT: Do we need backwards with v0? https://docs.pokt.network/supported-blockchains/
	RelayChain string `protobuf:"bytes,5,opt,name=relay_chain,json=relayChain,proto3" json:"relay_chain,omitempty"`
	// CONSIDERATION: Should a single session support multiple geo zones?
	GeoZone string `protobuf:"bytes,6,opt,name=geo_zone,json=geoZone,proto3" json:"geo_zone,omitempty"`
}

func (m *Session) Reset()         { *m = Session{} }
func (m *Session) String() string { return proto.CompactTextString(m) }
func (*Session) ProtoMessage()    {}
func (*Session) Descriptor() ([]byte, []int) {
	return fileDescriptor_7de4a7d078a82fb8, []int{0}
}
func (m *Session) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Session) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Session.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Session) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Session.Merge(m, src)
}
func (m *Session) XXX_Size() int {
	return m.Size()
}
func (m *Session) XXX_DiscardUnknown() {
	xxx_messageInfo_Session.DiscardUnknown(m)
}

var xxx_messageInfo_Session proto.InternalMessageInfo

func (m *Session) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Session) GetSessionNumber() int64 {
	if m != nil {
		return m.SessionNumber
	}
	return 0
}

func (m *Session) GetSessionHeight() int64 {
	if m != nil {
		return m.SessionHeight
	}
	return 0
}

func (m *Session) GetNumSessionBlocks() int64 {
	if m != nil {
		return m.NumSessionBlocks
	}
	return 0
}

func (m *Session) GetRelayChain() string {
	if m != nil {
		return m.RelayChain
	}
	return ""
}

func (m *Session) GetGeoZone() string {
	if m != nil {
		return m.GeoZone
	}
	return ""
}

func init() {
	proto.RegisterType((*Session)(nil), "poktroll.poktroll.Session")
}

func init() { proto.RegisterFile("poktroll/poktroll/session.proto", fileDescriptor_7de4a7d078a82fb8) }

var fileDescriptor_7de4a7d078a82fb8 = []byte{
	// 245 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x2f, 0xc8, 0xcf, 0x2e,
	0x29, 0xca, 0xcf, 0xc9, 0xd1, 0x87, 0x33, 0x8a, 0x53, 0x8b, 0x8b, 0x33, 0xf3, 0xf3, 0xf4, 0x0a,
	0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x04, 0x61, 0xe2, 0x7a, 0x30, 0x86, 0xd2, 0x45, 0x46, 0x2e, 0xf6,
	0x60, 0x88, 0x22, 0x21, 0x3e, 0x2e, 0xa6, 0xcc, 0x14, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xce, 0x20,
	0xa6, 0xcc, 0x14, 0x21, 0x55, 0x2e, 0x3e, 0xa8, 0xfe, 0xf8, 0xbc, 0xd2, 0xdc, 0xa4, 0xd4, 0x22,
	0x09, 0x26, 0x05, 0x46, 0x0d, 0xe6, 0x20, 0x5e, 0xa8, 0xa8, 0x1f, 0x58, 0x10, 0x59, 0x59, 0x46,
	0x6a, 0x66, 0x7a, 0x46, 0x89, 0x04, 0x33, 0x8a, 0x32, 0x0f, 0xb0, 0xa0, 0x90, 0x0e, 0x97, 0x50,
	0x5e, 0x69, 0x6e, 0x3c, 0x4c, 0x69, 0x52, 0x4e, 0x7e, 0x72, 0x76, 0xb1, 0x04, 0x0b, 0x58, 0xa9,
	0x40, 0x5e, 0x69, 0x2e, 0xd4, 0x15, 0x4e, 0x60, 0x71, 0x21, 0x79, 0x2e, 0xee, 0xa2, 0xd4, 0x9c,
	0xc4, 0xca, 0xf8, 0xe4, 0x8c, 0xc4, 0xcc, 0x3c, 0x09, 0x56, 0xb0, 0xa3, 0xb8, 0xc0, 0x42, 0xce,
	0x20, 0x11, 0x21, 0x49, 0x2e, 0x8e, 0xf4, 0xd4, 0xfc, 0xf8, 0xaa, 0xfc, 0xbc, 0x54, 0x09, 0x36,
	0xb0, 0x2c, 0x7b, 0x7a, 0x6a, 0x7e, 0x54, 0x7e, 0x5e, 0xaa, 0x93, 0xf1, 0x89, 0x47, 0x72, 0x8c,
	0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24, 0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72,
	0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0x49, 0xc2, 0x03, 0xa6, 0x02, 0x11, 0x46, 0x25, 0x95, 0x05,
	0xa9, 0xc5, 0x49, 0x6c, 0xe0, 0x20, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xec, 0xd3, 0x17,
	0x6b, 0x45, 0x01, 0x00, 0x00,
}

func (m *Session) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Session) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Session) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.GeoZone) > 0 {
		i -= len(m.GeoZone)
		copy(dAtA[i:], m.GeoZone)
		i = encodeVarintSession(dAtA, i, uint64(len(m.GeoZone)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.RelayChain) > 0 {
		i -= len(m.RelayChain)
		copy(dAtA[i:], m.RelayChain)
		i = encodeVarintSession(dAtA, i, uint64(len(m.RelayChain)))
		i--
		dAtA[i] = 0x2a
	}
	if m.NumSessionBlocks != 0 {
		i = encodeVarintSession(dAtA, i, uint64(m.NumSessionBlocks))
		i--
		dAtA[i] = 0x20
	}
	if m.SessionHeight != 0 {
		i = encodeVarintSession(dAtA, i, uint64(m.SessionHeight))
		i--
		dAtA[i] = 0x18
	}
	if m.SessionNumber != 0 {
		i = encodeVarintSession(dAtA, i, uint64(m.SessionNumber))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Id) > 0 {
		i -= len(m.Id)
		copy(dAtA[i:], m.Id)
		i = encodeVarintSession(dAtA, i, uint64(len(m.Id)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSession(dAtA []byte, offset int, v uint64) int {
	offset -= sovSession(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Session) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Id)
	if l > 0 {
		n += 1 + l + sovSession(uint64(l))
	}
	if m.SessionNumber != 0 {
		n += 1 + sovSession(uint64(m.SessionNumber))
	}
	if m.SessionHeight != 0 {
		n += 1 + sovSession(uint64(m.SessionHeight))
	}
	if m.NumSessionBlocks != 0 {
		n += 1 + sovSession(uint64(m.NumSessionBlocks))
	}
	l = len(m.RelayChain)
	if l > 0 {
		n += 1 + l + sovSession(uint64(l))
	}
	l = len(m.GeoZone)
	if l > 0 {
		n += 1 + l + sovSession(uint64(l))
	}
	return n
}

func sovSession(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSession(x uint64) (n int) {
	return sovSession(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Session) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSession
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
			return fmt.Errorf("proto: Session: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Session: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSession
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSession
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Id = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SessionNumber", wireType)
			}
			m.SessionNumber = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SessionNumber |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SessionHeight", wireType)
			}
			m.SessionHeight = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SessionHeight |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field NumSessionBlocks", wireType)
			}
			m.NumSessionBlocks = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.NumSessionBlocks |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RelayChain", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSession
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSession
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RelayChain = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field GeoZone", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSession
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthSession
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSession
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.GeoZone = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipSession(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSession
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
func skipSession(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSession
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
					return 0, ErrIntOverflowSession
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
					return 0, ErrIntOverflowSession
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
				return 0, ErrInvalidLengthSession
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSession
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSession
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSession        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSession          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSession = fmt.Errorf("proto: unexpected end of group")
)
