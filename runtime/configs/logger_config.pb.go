// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v4.23.4
// source: logger_config.proto

package configs

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogLevel int32

const (
	LogLevel_LOG_LEVEL_UNSPECIFIED LogLevel = 0
	LogLevel_LOG_LEVEL_DEBUG       LogLevel = 1
	LogLevel_LOG_LEVEL_INFO        LogLevel = 2
	LogLevel_LOG_LEVEL_WARN        LogLevel = 3
	LogLevel_LOG_LEVEL_ERROR       LogLevel = 4
	LogLevel_LOG_LEVEL_FATAL       LogLevel = 5
	LogLevel_LOG_LEVEL_PANIC       LogLevel = 6
)

// Enum value maps for LogLevel.
var (
	LogLevel_name = map[int32]string{
		0: "LOG_LEVEL_UNSPECIFIED",
		1: "LOG_LEVEL_DEBUG",
		2: "LOG_LEVEL_INFO",
		3: "LOG_LEVEL_WARN",
		4: "LOG_LEVEL_ERROR",
		5: "LOG_LEVEL_FATAL",
		6: "LOG_LEVEL_PANIC",
	}
	LogLevel_value = map[string]int32{
		"LOG_LEVEL_UNSPECIFIED": 0,
		"LOG_LEVEL_DEBUG":       1,
		"LOG_LEVEL_INFO":        2,
		"LOG_LEVEL_WARN":        3,
		"LOG_LEVEL_ERROR":       4,
		"LOG_LEVEL_FATAL":       5,
		"LOG_LEVEL_PANIC":       6,
	}
)

func (x LogLevel) Enum() *LogLevel {
	p := new(LogLevel)
	*p = x
	return p
}

func (x LogLevel) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogLevel) Descriptor() protoreflect.EnumDescriptor {
	return file_logger_config_proto_enumTypes[0].Descriptor()
}

func (LogLevel) Type() protoreflect.EnumType {
	return &file_logger_config_proto_enumTypes[0]
}

func (x LogLevel) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogLevel.Descriptor instead.
func (LogLevel) EnumDescriptor() ([]byte, []int) {
	return file_logger_config_proto_rawDescGZIP(), []int{0}
}

type LogFormat int32

const (
	LogFormat_LOG_FORMAT_UNSPECIFIED LogFormat = 0
	LogFormat_LOG_FORMAT_JSON        LogFormat = 1
	LogFormat_LOG_FORMAT_PRETTY      LogFormat = 2
)

// Enum value maps for LogFormat.
var (
	LogFormat_name = map[int32]string{
		0: "LOG_FORMAT_UNSPECIFIED",
		1: "LOG_FORMAT_JSON",
		2: "LOG_FORMAT_PRETTY",
	}
	LogFormat_value = map[string]int32{
		"LOG_FORMAT_UNSPECIFIED": 0,
		"LOG_FORMAT_JSON":        1,
		"LOG_FORMAT_PRETTY":      2,
	}
)

func (x LogFormat) Enum() *LogFormat {
	p := new(LogFormat)
	*p = x
	return p
}

func (x LogFormat) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (LogFormat) Descriptor() protoreflect.EnumDescriptor {
	return file_logger_config_proto_enumTypes[1].Descriptor()
}

func (LogFormat) Type() protoreflect.EnumType {
	return &file_logger_config_proto_enumTypes[1]
}

func (x LogFormat) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use LogFormat.Descriptor instead.
func (LogFormat) EnumDescriptor() ([]byte, []int) {
	return file_logger_config_proto_rawDescGZIP(), []int{1}
}

type LoggerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Level  string `protobuf:"bytes,1,opt,name=level,proto3" json:"level,omitempty"`
	Format string `protobuf:"bytes,2,opt,name=format,proto3" json:"format,omitempty"`
}

func (x *LoggerConfig) Reset() {
	*x = LoggerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_logger_config_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoggerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoggerConfig) ProtoMessage() {}

func (x *LoggerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_logger_config_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoggerConfig.ProtoReflect.Descriptor instead.
func (*LoggerConfig) Descriptor() ([]byte, []int) {
	return file_logger_config_proto_rawDescGZIP(), []int{0}
}

func (x *LoggerConfig) GetLevel() string {
	if x != nil {
		return x.Level
	}
	return ""
}

func (x *LoggerConfig) GetFormat() string {
	if x != nil {
		return x.Format
	}
	return ""
}

var File_logger_config_proto protoreflect.FileDescriptor

var file_logger_config_proto_rawDesc = []byte{
	0x0a, 0x13, 0x6c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x22, 0x3c,
	0x0a, 0x0c, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x14,
	0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c,
	0x65, 0x76, 0x65, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x66, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x2a, 0xa1, 0x01, 0x0a,
	0x08, 0x4c, 0x6f, 0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x12, 0x19, 0x0a, 0x15, 0x4c, 0x4f, 0x47,
	0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49,
	0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45,
	0x4c, 0x5f, 0x44, 0x45, 0x42, 0x55, 0x47, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x4c, 0x4f, 0x47,
	0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x02, 0x12, 0x12, 0x0a,
	0x0e, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x57, 0x41, 0x52, 0x4e, 0x10,
	0x03, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x04, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x4f, 0x47, 0x5f, 0x4c, 0x45,
	0x56, 0x45, 0x4c, 0x5f, 0x46, 0x41, 0x54, 0x41, 0x4c, 0x10, 0x05, 0x12, 0x13, 0x0a, 0x0f, 0x4c,
	0x4f, 0x47, 0x5f, 0x4c, 0x45, 0x56, 0x45, 0x4c, 0x5f, 0x50, 0x41, 0x4e, 0x49, 0x43, 0x10, 0x06,
	0x2a, 0x53, 0x0a, 0x09, 0x4c, 0x6f, 0x67, 0x46, 0x6f, 0x72, 0x6d, 0x61, 0x74, 0x12, 0x1a, 0x0a,
	0x16, 0x4c, 0x4f, 0x47, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x55, 0x4e, 0x53, 0x50,
	0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x4c, 0x4f, 0x47,
	0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x4a, 0x53, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x15,
	0x0a, 0x11, 0x4c, 0x4f, 0x47, 0x5f, 0x46, 0x4f, 0x52, 0x4d, 0x41, 0x54, 0x5f, 0x50, 0x52, 0x45,
	0x54, 0x54, 0x59, 0x10, 0x02, 0x42, 0x32, 0x5a, 0x30, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x70, 0x6f, 0x6b, 0x74, 0x2d, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b,
	0x2f, 0x63, 0x6d, 0x74, 0x2d, 0x70, 0x6f, 0x6b, 0x74, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_logger_config_proto_rawDescOnce sync.Once
	file_logger_config_proto_rawDescData = file_logger_config_proto_rawDesc
)

func file_logger_config_proto_rawDescGZIP() []byte {
	file_logger_config_proto_rawDescOnce.Do(func() {
		file_logger_config_proto_rawDescData = protoimpl.X.CompressGZIP(file_logger_config_proto_rawDescData)
	})
	return file_logger_config_proto_rawDescData
}

var file_logger_config_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_logger_config_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_logger_config_proto_goTypes = []interface{}{
	(LogLevel)(0),        // 0: configs.LogLevel
	(LogFormat)(0),       // 1: configs.LogFormat
	(*LoggerConfig)(nil), // 2: configs.LoggerConfig
}
var file_logger_config_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_logger_config_proto_init() }
func file_logger_config_proto_init() {
	if File_logger_config_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_logger_config_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoggerConfig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_logger_config_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_logger_config_proto_goTypes,
		DependencyIndexes: file_logger_config_proto_depIdxs,
		EnumInfos:         file_logger_config_proto_enumTypes,
		MessageInfos:      file_logger_config_proto_msgTypes,
	}.Build()
	File_logger_config_proto = out.File
	file_logger_config_proto_rawDesc = nil
	file_logger_config_proto_goTypes = nil
	file_logger_config_proto_depIdxs = nil
}
