// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: poktroll/application/application.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	types "github.com/cosmos/cosmos-sdk/types"
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

type Application struct {
	Address string      `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Stake   *types.Coin `protobuf:"bytes,2,opt,name=stake,proto3" json:"stake,omitempty"`
}

func (m *Application) Reset()         { *m = Application{} }
func (m *Application) String() string { return proto.CompactTextString(m) }
func (*Application) ProtoMessage()    {}
func (*Application) Descriptor() ([]byte, []int) {
	return fileDescriptor_5933c748ae9780b4, []int{0}
}
func (m *Application) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Application) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Application.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Application) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Application.Merge(m, src)
}
func (m *Application) XXX_Size() int {
	return m.Size()
}
func (m *Application) XXX_DiscardUnknown() {
	xxx_messageInfo_Application.DiscardUnknown(m)
}

var xxx_messageInfo_Application proto.InternalMessageInfo

func (m *Application) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *Application) GetStake() *types.Coin {
	if m != nil {
		return m.Stake
	}
	return nil
}

func init() {
	proto.RegisterType((*Application)(nil), "poktroll.application.Application")
}

func init() {
	proto.RegisterFile("poktroll/application/application.proto", fileDescriptor_5933c748ae9780b4)
}

var fileDescriptor_5933c748ae9780b4 = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2b, 0xc8, 0xcf, 0x2e,
	0x29, 0xca, 0xcf, 0xc9, 0xd1, 0x4f, 0x2c, 0x28, 0xc8, 0xc9, 0x4c, 0x4e, 0x2c, 0xc9, 0xcc, 0xcf,
	0x43, 0x66, 0xeb, 0x15, 0x14, 0xe5, 0x97, 0xe4, 0x0b, 0x89, 0xc0, 0xd4, 0xe9, 0x21, 0xc9, 0x49,
	0x49, 0x26, 0xe7, 0x17, 0xe7, 0xe6, 0x17, 0xc7, 0x83, 0xd5, 0xe8, 0x43, 0x38, 0x10, 0x0d, 0x52,
	0x72, 0x10, 0x9e, 0x7e, 0x52, 0x62, 0x71, 0xaa, 0x7e, 0x99, 0x61, 0x52, 0x6a, 0x49, 0xa2, 0xa1,
	0x7e, 0x72, 0x7e, 0x26, 0xd4, 0x40, 0xa5, 0x22, 0x2e, 0x6e, 0x47, 0x84, 0x49, 0x42, 0x46, 0x5c,
	0xec, 0x89, 0x29, 0x29, 0x45, 0xa9, 0xc5, 0xc5, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x9c, 0x4e, 0x12,
	0x97, 0xb6, 0xe8, 0x8a, 0x40, 0x4d, 0x74, 0x84, 0xc8, 0x04, 0x97, 0x14, 0x65, 0xe6, 0xa5, 0x07,
	0xc1, 0x14, 0x0a, 0xe9, 0x73, 0xb1, 0x16, 0x97, 0x24, 0x66, 0xa7, 0x4a, 0x30, 0x29, 0x30, 0x6a,
	0x70, 0x1b, 0x49, 0xea, 0x41, 0x95, 0x83, 0xac, 0xd4, 0x83, 0x5a, 0xa9, 0xe7, 0x9c, 0x9f, 0x99,
	0x17, 0x04, 0x51, 0xe7, 0x64, 0x76, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72, 0x8c, 0x0f, 0x1e,
	0xc9, 0x31, 0x4e, 0x78, 0x2c, 0xc7, 0x70, 0xe1, 0xb1, 0x1c, 0xc3, 0x8d, 0xc7, 0x72, 0x0c, 0x51,
	0x32, 0xf0, 0x60, 0xa8, 0x40, 0x09, 0x88, 0x92, 0xca, 0x82, 0xd4, 0xe2, 0x24, 0x36, 0xb0, 0x93,
	0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x8f, 0x84, 0x17, 0xcf, 0x2d, 0x01, 0x00, 0x00,
}

func (m *Application) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Application) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Application) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Stake != nil {
		{
			size, err := m.Stake.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintApplication(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.Address) > 0 {
		i -= len(m.Address)
		copy(dAtA[i:], m.Address)
		i = encodeVarintApplication(dAtA, i, uint64(len(m.Address)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintApplication(dAtA []byte, offset int, v uint64) int {
	offset -= sovApplication(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Application) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Address)
	if l > 0 {
		n += 1 + l + sovApplication(uint64(l))
	}
	if m.Stake != nil {
		l = m.Stake.Size()
		n += 1 + l + sovApplication(uint64(l))
	}
	return n
}

func sovApplication(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozApplication(x uint64) (n int) {
	return sovApplication(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Application) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowApplication
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
			return fmt.Errorf("proto: Application: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Application: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Address", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
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
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthApplication
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Address = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Stake", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowApplication
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthApplication
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthApplication
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Stake == nil {
				m.Stake = &types.Coin{}
			}
			if err := m.Stake.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipApplication(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthApplication
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
func skipApplication(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowApplication
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
					return 0, ErrIntOverflowApplication
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
					return 0, ErrIntOverflowApplication
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
				return 0, ErrInvalidLengthApplication
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupApplication
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthApplication
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthApplication        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowApplication          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupApplication = fmt.Errorf("proto: unexpected end of group")
)
