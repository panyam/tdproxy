// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: protos/tickers.proto

package protos

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Ticker struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol          string           `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	LastRefreshedAt string           `protobuf:"bytes,2,opt,name=last_refreshed_at,json=lastRefreshedAt,proto3" json:"last_refreshed_at,omitempty"`
	Info            *structpb.Struct `protobuf:"bytes,3,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *Ticker) Reset() {
	*x = Ticker{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_tickers_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Ticker) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Ticker) ProtoMessage() {}

func (x *Ticker) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tickers_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Ticker.ProtoReflect.Descriptor instead.
func (*Ticker) Descriptor() ([]byte, []int) {
	return file_protos_tickers_proto_rawDescGZIP(), []int{0}
}

func (x *Ticker) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Ticker) GetLastRefreshedAt() string {
	if x != nil {
		return x.LastRefreshedAt
	}
	return ""
}

func (x *Ticker) GetInfo() *structpb.Struct {
	if x != nil {
		return x.Info
	}
	return nil
}

type GetTickersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbols []string `protobuf:"bytes,1,rep,name=symbols,proto3" json:"symbols,omitempty"`
	//*
	// MAX_INT effective means never refresh
	// <= 0 => Always refresh
	// Any other value indicates to only refresh if last refresh
	// was before this threshold.
	RefreshType *int32 `protobuf:"varint,2,opt,name=refresh_type,json=refreshType,proto3,oneof" json:"refresh_type,omitempty"`
}

func (x *GetTickersRequest) Reset() {
	*x = GetTickersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_tickers_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTickersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTickersRequest) ProtoMessage() {}

func (x *GetTickersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tickers_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTickersRequest.ProtoReflect.Descriptor instead.
func (*GetTickersRequest) Descriptor() ([]byte, []int) {
	return file_protos_tickers_proto_rawDescGZIP(), []int{1}
}

func (x *GetTickersRequest) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

func (x *GetTickersRequest) GetRefreshType() int32 {
	if x != nil && x.RefreshType != nil {
		return *x.RefreshType
	}
	return 0
}

type GetTickersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status  bool               `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Tickers map[string]*Ticker `protobuf:"bytes,2,rep,name=tickers,proto3" json:"tickers,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Errors  map[string]string  `protobuf:"bytes,3,rep,name=errors,proto3" json:"errors,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GetTickersResponse) Reset() {
	*x = GetTickersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_tickers_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTickersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTickersResponse) ProtoMessage() {}

func (x *GetTickersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_tickers_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTickersResponse.ProtoReflect.Descriptor instead.
func (*GetTickersResponse) Descriptor() ([]byte, []int) {
	return file_protos_tickers_proto_rawDescGZIP(), []int{2}
}

func (x *GetTickersResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

func (x *GetTickersResponse) GetTickers() map[string]*Ticker {
	if x != nil {
		return x.Tickers
	}
	return nil
}

func (x *GetTickersResponse) GetErrors() map[string]string {
	if x != nil {
		return x.Errors
	}
	return nil
}

var File_protos_tickers_proto protoreflect.FileDescriptor

var file_protos_tickers_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x73,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x79, 0x0a, 0x06,
	0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x2a,
	0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x52,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x65, 0x64, 0x41, 0x74, 0x12, 0x2b, 0x0a, 0x04, 0x69, 0x6e,
	0x66, 0x6f, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22, 0x66, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x12, 0x26, 0x0a, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73,
	0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x0b,
	0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0f,
	0x0a, 0x0d, 0x5f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22,
	0xb6, 0x02, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x41,
	0x0a, 0x07, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x07, 0x74, 0x69, 0x63, 0x6b, 0x65, 0x72,
	0x73, 0x12, 0x3e, 0x0a, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x26, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x69,
	0x63, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x45, 0x72,
	0x72, 0x6f, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x65, 0x72, 0x72, 0x6f, 0x72,
	0x73, 0x1a, 0x4a, 0x0a, 0x0c, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a, 0x39, 0x0a,
	0x0b, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x56, 0x0a, 0x0d, 0x54, 0x69, 0x63, 0x6b,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x45, 0x0a, 0x0a, 0x47, 0x65, 0x74,
	0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x12, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2e, 0x47, 0x65, 0x74, 0x54, 0x69, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x54,
	0x69, 0x63, 0x6b, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x42, 0x1a, 0x5a, 0x18, 0x6c, 0x65, 0x67, 0x66, 0x69, 0x6e, 0x64, 0x65, 0x72, 0x2f, 0x74, 0x64,
	0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_tickers_proto_rawDescOnce sync.Once
	file_protos_tickers_proto_rawDescData = file_protos_tickers_proto_rawDesc
)

func file_protos_tickers_proto_rawDescGZIP() []byte {
	file_protos_tickers_proto_rawDescOnce.Do(func() {
		file_protos_tickers_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_tickers_proto_rawDescData)
	})
	return file_protos_tickers_proto_rawDescData
}

var file_protos_tickers_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_protos_tickers_proto_goTypes = []interface{}{
	(*Ticker)(nil),             // 0: protos.Ticker
	(*GetTickersRequest)(nil),  // 1: protos.GetTickersRequest
	(*GetTickersResponse)(nil), // 2: protos.GetTickersResponse
	nil,                        // 3: protos.GetTickersResponse.TickersEntry
	nil,                        // 4: protos.GetTickersResponse.ErrorsEntry
	(*structpb.Struct)(nil),    // 5: google.protobuf.Struct
}
var file_protos_tickers_proto_depIdxs = []int32{
	5, // 0: protos.Ticker.info:type_name -> google.protobuf.Struct
	3, // 1: protos.GetTickersResponse.tickers:type_name -> protos.GetTickersResponse.TickersEntry
	4, // 2: protos.GetTickersResponse.errors:type_name -> protos.GetTickersResponse.ErrorsEntry
	0, // 3: protos.GetTickersResponse.TickersEntry.value:type_name -> protos.Ticker
	1, // 4: protos.TickerService.GetTickers:input_type -> protos.GetTickersRequest
	2, // 5: protos.TickerService.GetTickers:output_type -> protos.GetTickersResponse
	5, // [5:6] is the sub-list for method output_type
	4, // [4:5] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_protos_tickers_proto_init() }
func file_protos_tickers_proto_init() {
	if File_protos_tickers_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_tickers_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Ticker); i {
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
		file_protos_tickers_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTickersRequest); i {
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
		file_protos_tickers_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTickersResponse); i {
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
	file_protos_tickers_proto_msgTypes[1].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_tickers_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_tickers_proto_goTypes,
		DependencyIndexes: file_protos_tickers_proto_depIdxs,
		MessageInfos:      file_protos_tickers_proto_msgTypes,
	}.Build()
	File_protos_tickers_proto = out.File
	file_protos_tickers_proto_rawDesc = nil
	file_protos_tickers_proto_goTypes = nil
	file_protos_tickers_proto_depIdxs = nil
}