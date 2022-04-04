// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: protos/chains.proto

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

type Option struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol       string           `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	DateString   string           `protobuf:"bytes,2,opt,name=date_string,json=dateString,proto3" json:"date_string,omitempty"`
	PriceString  string           `protobuf:"bytes,3,opt,name=price_string,json=priceString,proto3" json:"price_string,omitempty"`
	IsCall       bool             `protobuf:"varint,4,opt,name=is_call,json=isCall,proto3" json:"is_call,omitempty"`
	StrikePrice  float64          `protobuf:"fixed64,5,opt,name=strike_price,json=strikePrice,proto3" json:"strike_price,omitempty"`
	AskPrice     float64          `protobuf:"fixed64,6,opt,name=ask_price,json=askPrice,proto3" json:"ask_price,omitempty"`
	BidPrice     float64          `protobuf:"fixed64,7,opt,name=bid_price,json=bidPrice,proto3" json:"bid_price,omitempty"`
	MarkPrice    float64          `protobuf:"fixed64,8,opt,name=mark_price,json=markPrice,proto3" json:"mark_price,omitempty"`
	OpenInterest int32            `protobuf:"varint,9,opt,name=open_interest,json=openInterest,proto3" json:"open_interest,omitempty"`
	Delta        float64          `protobuf:"fixed64,10,opt,name=delta,proto3" json:"delta,omitempty"`
	Multiplier   float64          `protobuf:"fixed64,11,opt,name=multiplier,proto3" json:"multiplier,omitempty"`
	Info         *structpb.Struct `protobuf:"bytes,12,opt,name=info,proto3" json:"info,omitempty"`
}

func (x *Option) Reset() {
	*x = Option{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chains_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Option) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Option) ProtoMessage() {}

func (x *Option) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chains_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Option.ProtoReflect.Descriptor instead.
func (*Option) Descriptor() ([]byte, []int) {
	return file_protos_chains_proto_rawDescGZIP(), []int{0}
}

func (x *Option) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Option) GetDateString() string {
	if x != nil {
		return x.DateString
	}
	return ""
}

func (x *Option) GetPriceString() string {
	if x != nil {
		return x.PriceString
	}
	return ""
}

func (x *Option) GetIsCall() bool {
	if x != nil {
		return x.IsCall
	}
	return false
}

func (x *Option) GetStrikePrice() float64 {
	if x != nil {
		return x.StrikePrice
	}
	return 0
}

func (x *Option) GetAskPrice() float64 {
	if x != nil {
		return x.AskPrice
	}
	return 0
}

func (x *Option) GetBidPrice() float64 {
	if x != nil {
		return x.BidPrice
	}
	return 0
}

func (x *Option) GetMarkPrice() float64 {
	if x != nil {
		return x.MarkPrice
	}
	return 0
}

func (x *Option) GetOpenInterest() int32 {
	if x != nil {
		return x.OpenInterest
	}
	return 0
}

func (x *Option) GetDelta() float64 {
	if x != nil {
		return x.Delta
	}
	return 0
}

func (x *Option) GetMultiplier() float64 {
	if x != nil {
		return x.Multiplier
	}
	return 0
}

func (x *Option) GetInfo() *structpb.Struct {
	if x != nil {
		return x.Info
	}
	return nil
}

type Chain struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol          string             `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Date            string             `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	IsCall          bool               `protobuf:"varint,3,opt,name=is_call,json=isCall,proto3" json:"is_call,omitempty"`
	LastRefreshedAt string             `protobuf:"bytes,4,opt,name=last_refreshed_at,json=lastRefreshedAt,proto3" json:"last_refreshed_at,omitempty"`
	Options         []*Option          `protobuf:"bytes,5,rep,name=options,proto3" json:"options,omitempty"`
	OptionsByPrice  map[string]*Option `protobuf:"bytes,6,rep,name=options_by_price,json=optionsByPrice,proto3" json:"options_by_price,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Chain) Reset() {
	*x = Chain{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chains_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Chain) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Chain) ProtoMessage() {}

func (x *Chain) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chains_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Chain.ProtoReflect.Descriptor instead.
func (*Chain) Descriptor() ([]byte, []int) {
	return file_protos_chains_proto_rawDescGZIP(), []int{1}
}

func (x *Chain) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Chain) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *Chain) GetIsCall() bool {
	if x != nil {
		return x.IsCall
	}
	return false
}

func (x *Chain) GetLastRefreshedAt() string {
	if x != nil {
		return x.LastRefreshedAt
	}
	return ""
}

func (x *Chain) GetOptions() []*Option {
	if x != nil {
		return x.Options
	}
	return nil
}

func (x *Chain) GetOptionsByPrice() map[string]*Option {
	if x != nil {
		return x.OptionsByPrice
	}
	return nil
}

type GetChainInfoRequest struct {
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

func (x *GetChainInfoRequest) Reset() {
	*x = GetChainInfoRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chains_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChainInfoRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChainInfoRequest) ProtoMessage() {}

func (x *GetChainInfoRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chains_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChainInfoRequest.ProtoReflect.Descriptor instead.
func (*GetChainInfoRequest) Descriptor() ([]byte, []int) {
	return file_protos_chains_proto_rawDescGZIP(), []int{2}
}

func (x *GetChainInfoRequest) GetSymbols() []string {
	if x != nil {
		return x.Symbols
	}
	return nil
}

func (x *GetChainInfoRequest) GetRefreshType() int32 {
	if x != nil && x.RefreshType != nil {
		return *x.RefreshType
	}
	return 0
}

type GetChainInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol          string   `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	Dates           []string `protobuf:"bytes,2,rep,name=dates,proto3" json:"dates,omitempty"`
	LastRefreshedAt string   `protobuf:"bytes,3,opt,name=last_refreshed_at,json=lastRefreshedAt,proto3" json:"last_refreshed_at,omitempty"`
	ErrorCode       int32    `protobuf:"varint,4,opt,name=error_code,json=errorCode,proto3" json:"error_code,omitempty"`
	ErrorMessage    string   `protobuf:"bytes,5,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
}

func (x *GetChainInfoResponse) Reset() {
	*x = GetChainInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chains_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChainInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChainInfoResponse) ProtoMessage() {}

func (x *GetChainInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chains_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChainInfoResponse.ProtoReflect.Descriptor instead.
func (*GetChainInfoResponse) Descriptor() ([]byte, []int) {
	return file_protos_chains_proto_rawDescGZIP(), []int{3}
}

func (x *GetChainInfoResponse) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *GetChainInfoResponse) GetDates() []string {
	if x != nil {
		return x.Dates
	}
	return nil
}

func (x *GetChainInfoResponse) GetLastRefreshedAt() string {
	if x != nil {
		return x.LastRefreshedAt
	}
	return ""
}

func (x *GetChainInfoResponse) GetErrorCode() int32 {
	if x != nil {
		return x.ErrorCode
	}
	return 0
}

func (x *GetChainInfoResponse) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

type GetChainRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Symbol string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	//*
	// Return the chain on a particular date.
	Date string `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	//*
	// Get the call or put chain
	IsCall bool `protobuf:"varint,3,opt,name=is_call,json=isCall,proto3" json:"is_call,omitempty"`
	//*
	// MAX_INT effective means never refresh
	// <= 0 => Always refresh
	// Any other value indicates to only refresh if last refresh
	// was before this threshold.
	RefreshType *int32 `protobuf:"varint,4,opt,name=refresh_type,json=refreshType,proto3,oneof" json:"refresh_type,omitempty"`
}

func (x *GetChainRequest) Reset() {
	*x = GetChainRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chains_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChainRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChainRequest) ProtoMessage() {}

func (x *GetChainRequest) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chains_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChainRequest.ProtoReflect.Descriptor instead.
func (*GetChainRequest) Descriptor() ([]byte, []int) {
	return file_protos_chains_proto_rawDescGZIP(), []int{4}
}

func (x *GetChainRequest) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *GetChainRequest) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *GetChainRequest) GetIsCall() bool {
	if x != nil {
		return x.IsCall
	}
	return false
}

func (x *GetChainRequest) GetRefreshType() int32 {
	if x != nil && x.RefreshType != nil {
		return *x.RefreshType
	}
	return 0
}

type GetChainResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ErrorMessage *string `protobuf:"bytes,1,opt,name=error_message,json=errorMessage,proto3,oneof" json:"error_message,omitempty"`
	Chain        *Chain  `protobuf:"bytes,2,opt,name=chain,proto3" json:"chain,omitempty"`
}

func (x *GetChainResponse) Reset() {
	*x = GetChainResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_chains_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetChainResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetChainResponse) ProtoMessage() {}

func (x *GetChainResponse) ProtoReflect() protoreflect.Message {
	mi := &file_protos_chains_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetChainResponse.ProtoReflect.Descriptor instead.
func (*GetChainResponse) Descriptor() ([]byte, []int) {
	return file_protos_chains_proto_rawDescGZIP(), []int{5}
}

func (x *GetChainResponse) GetErrorMessage() string {
	if x != nil && x.ErrorMessage != nil {
		return *x.ErrorMessage
	}
	return ""
}

func (x *GetChainResponse) GetChain() *Chain {
	if x != nil {
		return x.Chain
	}
	return nil
}

var File_protos_chains_proto protoreflect.FileDescriptor

var file_protos_chains_proto_rawDesc = []byte{
	0x0a, 0x13, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x1a, 0x1c, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73,
	0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x81, 0x03, 0x0a, 0x06,
	0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12, 0x1f,
	0x0a, 0x0b, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12,
	0x21, 0x0a, 0x0c, 0x70, 0x72, 0x69, 0x63, 0x65, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x69, 0x63, 0x65, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x63, 0x61, 0x6c, 0x6c, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x73,
	0x74, 0x72, 0x69, 0x6b, 0x65, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6b, 0x65, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x61, 0x73, 0x6b, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x08, 0x61, 0x73, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x62,
	0x69, 0x64, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x08,
	0x62, 0x69, 0x64, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x72, 0x6b,
	0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x09, 0x6d, 0x61,
	0x72, 0x6b, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6f, 0x70, 0x65, 0x6e, 0x5f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x65, 0x73, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c,
	0x6f, 0x70, 0x65, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05,
	0x64, 0x65, 0x6c, 0x74, 0x61, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x64, 0x65, 0x6c,
	0x74, 0x61, 0x12, 0x1e, 0x0a, 0x0a, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69, 0x65, 0x72,
	0x18, 0x0b, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0a, 0x6d, 0x75, 0x6c, 0x74, 0x69, 0x70, 0x6c, 0x69,
	0x65, 0x72, 0x12, 0x2b, 0x0a, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x52, 0x04, 0x69, 0x6e, 0x66, 0x6f, 0x22,
	0xc2, 0x02, 0x0a, 0x05, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x63, 0x61, 0x6c, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x2a,
	0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x65, 0x64,
	0x5f, 0x61, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x52,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x65, 0x64, 0x41, 0x74, 0x12, 0x28, 0x0a, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x07, 0x6f, 0x70, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x12, 0x4b, 0x0a, 0x10, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f,
	0x62, 0x79, 0x5f, 0x70, 0x72, 0x69, 0x63, 0x65, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x4f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x50, 0x72, 0x69, 0x63, 0x65, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x52, 0x0e, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x50, 0x72, 0x69, 0x63,
	0x65, 0x1a, 0x51, 0x0a, 0x13, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x42, 0x79, 0x50, 0x72,
	0x69, 0x63, 0x65, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x24, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2e, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0x68, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e,
	0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x73, 0x79,
	0x6d, 0x62, 0x6f, 0x6c, 0x73, 0x12, 0x26, 0x0a, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x0b, 0x72,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a,
	0x0d, 0x5f, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0xb4,
	0x01, 0x0a, 0x14, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f, 0x6c, 0x12,
	0x14, 0x0a, 0x05, 0x64, 0x61, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x05,
	0x64, 0x61, 0x74, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x11, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x72, 0x65,
	0x66, 0x72, 0x65, 0x73, 0x68, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x65, 0x64, 0x41,
	0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x43, 0x6f, 0x64, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x8f, 0x01, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61,
	0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x79, 0x6d,
	0x62, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x79, 0x6d, 0x62, 0x6f,
	0x6c, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x69, 0x73, 0x5f, 0x63, 0x61, 0x6c, 0x6c,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x69, 0x73, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x26,
	0x0a, 0x0c, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x48, 0x00, 0x52, 0x0b, 0x72, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54,
	0x79, 0x70, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0f, 0x0a, 0x0d, 0x5f, 0x72, 0x65, 0x66, 0x72, 0x65,
	0x73, 0x68, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0x73, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x28, 0x0a, 0x0d, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x88, 0x01, 0x01, 0x12, 0x23, 0x0a, 0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x52, 0x05, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x9e, 0x01, 0x0a,
	0x0c, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a,
	0x0c, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x1b, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12, 0x3f, 0x0a, 0x08,
	0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x12, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x18, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x10, 0x5a,
	0x0e, 0x74, 0x64, 0x70, 0x72, 0x6f, 0x78, 0x79, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_chains_proto_rawDescOnce sync.Once
	file_protos_chains_proto_rawDescData = file_protos_chains_proto_rawDesc
)

func file_protos_chains_proto_rawDescGZIP() []byte {
	file_protos_chains_proto_rawDescOnce.Do(func() {
		file_protos_chains_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_chains_proto_rawDescData)
	})
	return file_protos_chains_proto_rawDescData
}

var file_protos_chains_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_protos_chains_proto_goTypes = []interface{}{
	(*Option)(nil),               // 0: protos.Option
	(*Chain)(nil),                // 1: protos.Chain
	(*GetChainInfoRequest)(nil),  // 2: protos.GetChainInfoRequest
	(*GetChainInfoResponse)(nil), // 3: protos.GetChainInfoResponse
	(*GetChainRequest)(nil),      // 4: protos.GetChainRequest
	(*GetChainResponse)(nil),     // 5: protos.GetChainResponse
	nil,                          // 6: protos.Chain.OptionsByPriceEntry
	(*structpb.Struct)(nil),      // 7: google.protobuf.Struct
}
var file_protos_chains_proto_depIdxs = []int32{
	7, // 0: protos.Option.info:type_name -> google.protobuf.Struct
	0, // 1: protos.Chain.options:type_name -> protos.Option
	6, // 2: protos.Chain.options_by_price:type_name -> protos.Chain.OptionsByPriceEntry
	1, // 3: protos.GetChainResponse.chain:type_name -> protos.Chain
	0, // 4: protos.Chain.OptionsByPriceEntry.value:type_name -> protos.Option
	2, // 5: protos.ChainService.GetChainInfo:input_type -> protos.GetChainInfoRequest
	4, // 6: protos.ChainService.GetChain:input_type -> protos.GetChainRequest
	3, // 7: protos.ChainService.GetChainInfo:output_type -> protos.GetChainInfoResponse
	5, // 8: protos.ChainService.GetChain:output_type -> protos.GetChainResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_protos_chains_proto_init() }
func file_protos_chains_proto_init() {
	if File_protos_chains_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_chains_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Option); i {
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
		file_protos_chains_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Chain); i {
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
		file_protos_chains_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChainInfoRequest); i {
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
		file_protos_chains_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChainInfoResponse); i {
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
		file_protos_chains_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChainRequest); i {
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
		file_protos_chains_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetChainResponse); i {
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
	file_protos_chains_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_protos_chains_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_protos_chains_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_protos_chains_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_protos_chains_proto_goTypes,
		DependencyIndexes: file_protos_chains_proto_depIdxs,
		MessageInfos:      file_protos_chains_proto_msgTypes,
	}.Build()
	File_protos_chains_proto = out.File
	file_protos_chains_proto_rawDesc = nil
	file_protos_chains_proto_goTypes = nil
	file_protos_chains_proto_depIdxs = nil
}
