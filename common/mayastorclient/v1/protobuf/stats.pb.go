// gRPC protobuf for IOStats

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.8
// source: stats.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ListStatsOption struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name *string `protobuf:"bytes,1,opt,name=name,proto3,oneof" json:"name,omitempty"` // If specified, returns IO stats of only specified resource.
}

func (x *ListStatsOption) Reset() {
	*x = ListStatsOption{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stats_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListStatsOption) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListStatsOption) ProtoMessage() {}

func (x *ListStatsOption) ProtoReflect() protoreflect.Message {
	mi := &file_stats_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListStatsOption.ProtoReflect.Descriptor instead.
func (*ListStatsOption) Descriptor() ([]byte, []int) {
	return file_stats_proto_rawDescGZIP(), []int{0}
}

func (x *ListStatsOption) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

type PoolIoStatsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Stats []*IoStats `protobuf:"bytes,1,rep,name=stats,proto3" json:"stats,omitempty"`
}

func (x *PoolIoStatsResponse) Reset() {
	*x = PoolIoStatsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stats_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PoolIoStatsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PoolIoStatsResponse) ProtoMessage() {}

func (x *PoolIoStatsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_stats_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PoolIoStatsResponse.ProtoReflect.Descriptor instead.
func (*PoolIoStatsResponse) Descriptor() ([]byte, []int) {
	return file_stats_proto_rawDescGZIP(), []int{1}
}

func (x *PoolIoStatsResponse) GetStats() []*IoStats {
	if x != nil {
		return x.Stats
	}
	return nil
}

// Common IO stats struct for all resource types.
type IoStats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name                 string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`                                                                   // Name of the device/resource.
	NumReadOps           uint64 `protobuf:"varint,2,opt,name=num_read_ops,json=numReadOps,proto3" json:"num_read_ops,omitempty"`                                  // Number of read operations on the device.
	BytesRead            uint64 `protobuf:"varint,3,opt,name=bytes_read,json=bytesRead,proto3" json:"bytes_read,omitempty"`                                       // Total bytes read on the device.
	NumWriteOps          uint64 `protobuf:"varint,4,opt,name=num_write_ops,json=numWriteOps,proto3" json:"num_write_ops,omitempty"`                               // Number of write operations on the device.
	BytesWritten         uint64 `protobuf:"varint,5,opt,name=bytes_written,json=bytesWritten,proto3" json:"bytes_written,omitempty"`                              // Total bytes written on the device.
	NumUnmapOps          uint64 `protobuf:"varint,6,opt,name=num_unmap_ops,json=numUnmapOps,proto3" json:"num_unmap_ops,omitempty"`                               // Number of unmap operations on the device.
	BytesUnmapped        uint64 `protobuf:"varint,7,opt,name=bytes_unmapped,json=bytesUnmapped,proto3" json:"bytes_unmapped,omitempty"`                           // Total bytes unmapped on the device.
	ReadLatencyTicks     uint64 `protobuf:"varint,8,opt,name=read_latency_ticks,json=readLatencyTicks,proto3" json:"read_latency_ticks,omitempty"`                // Total read latency on the device in ticks.
	WriteLatencyTicks    uint64 `protobuf:"varint,9,opt,name=write_latency_ticks,json=writeLatencyTicks,proto3" json:"write_latency_ticks,omitempty"`             // Total write latency on the device in ticks.
	UnmapLatencyTicks    uint64 `protobuf:"varint,10,opt,name=unmap_latency_ticks,json=unmapLatencyTicks,proto3" json:"unmap_latency_ticks,omitempty"`            // Total unmap latency on the device in ticks.
	MaxReadLatencyTicks  uint64 `protobuf:"varint,11,opt,name=max_read_latency_ticks,json=maxReadLatencyTicks,proto3" json:"max_read_latency_ticks,omitempty"`    // Max read latency in ticks occurred over num_read_ops.
	MinReadLatencyTicks  uint64 `protobuf:"varint,12,opt,name=min_read_latency_ticks,json=minReadLatencyTicks,proto3" json:"min_read_latency_ticks,omitempty"`    // Min read latency in ticks occurred over num_read_ops.
	MaxWriteLatencyTicks uint64 `protobuf:"varint,13,opt,name=max_write_latency_ticks,json=maxWriteLatencyTicks,proto3" json:"max_write_latency_ticks,omitempty"` // Max write latency in ticks occurred over num_write_ops.
	MinWriteLatencyTicks uint64 `protobuf:"varint,14,opt,name=min_write_latency_ticks,json=minWriteLatencyTicks,proto3" json:"min_write_latency_ticks,omitempty"` // Min write latency in ticks occurred over num_write_ops.
	MaxUnmapLatencyTicks uint64 `protobuf:"varint,15,opt,name=max_unmap_latency_ticks,json=maxUnmapLatencyTicks,proto3" json:"max_unmap_latency_ticks,omitempty"` // Max unmap latency in ticks occurred over num_unmap_ops.
	MinUnmapLatencyTicks uint64 `protobuf:"varint,16,opt,name=min_unmap_latency_ticks,json=minUnmapLatencyTicks,proto3" json:"min_unmap_latency_ticks,omitempty"` // Min unmap latency in ticks occurred over num_unmap_ops.
	TickRate             uint64 `protobuf:"varint,17,opt,name=tick_rate,json=tickRate,proto3" json:"tick_rate,omitempty"`                                         // Tick rate specific to node hosting the device.
}

func (x *IoStats) Reset() {
	*x = IoStats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_stats_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IoStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IoStats) ProtoMessage() {}

func (x *IoStats) ProtoReflect() protoreflect.Message {
	mi := &file_stats_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IoStats.ProtoReflect.Descriptor instead.
func (*IoStats) Descriptor() ([]byte, []int) {
	return file_stats_proto_rawDescGZIP(), []int{2}
}

func (x *IoStats) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *IoStats) GetNumReadOps() uint64 {
	if x != nil {
		return x.NumReadOps
	}
	return 0
}

func (x *IoStats) GetBytesRead() uint64 {
	if x != nil {
		return x.BytesRead
	}
	return 0
}

func (x *IoStats) GetNumWriteOps() uint64 {
	if x != nil {
		return x.NumWriteOps
	}
	return 0
}

func (x *IoStats) GetBytesWritten() uint64 {
	if x != nil {
		return x.BytesWritten
	}
	return 0
}

func (x *IoStats) GetNumUnmapOps() uint64 {
	if x != nil {
		return x.NumUnmapOps
	}
	return 0
}

func (x *IoStats) GetBytesUnmapped() uint64 {
	if x != nil {
		return x.BytesUnmapped
	}
	return 0
}

func (x *IoStats) GetReadLatencyTicks() uint64 {
	if x != nil {
		return x.ReadLatencyTicks
	}
	return 0
}

func (x *IoStats) GetWriteLatencyTicks() uint64 {
	if x != nil {
		return x.WriteLatencyTicks
	}
	return 0
}

func (x *IoStats) GetUnmapLatencyTicks() uint64 {
	if x != nil {
		return x.UnmapLatencyTicks
	}
	return 0
}

func (x *IoStats) GetMaxReadLatencyTicks() uint64 {
	if x != nil {
		return x.MaxReadLatencyTicks
	}
	return 0
}

func (x *IoStats) GetMinReadLatencyTicks() uint64 {
	if x != nil {
		return x.MinReadLatencyTicks
	}
	return 0
}

func (x *IoStats) GetMaxWriteLatencyTicks() uint64 {
	if x != nil {
		return x.MaxWriteLatencyTicks
	}
	return 0
}

func (x *IoStats) GetMinWriteLatencyTicks() uint64 {
	if x != nil {
		return x.MinWriteLatencyTicks
	}
	return 0
}

func (x *IoStats) GetMaxUnmapLatencyTicks() uint64 {
	if x != nil {
		return x.MaxUnmapLatencyTicks
	}
	return 0
}

func (x *IoStats) GetMinUnmapLatencyTicks() uint64 {
	if x != nil {
		return x.MinUnmapLatencyTicks
	}
	return 0
}

func (x *IoStats) GetTickRate() uint64 {
	if x != nil {
		return x.TickRate
	}
	return 0
}

var File_stats_proto protoreflect.FileDescriptor

var file_stats_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x73, 0x74, 0x61, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x6d,
	0x61, 0x79, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74,
	0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x33, 0x0a, 0x0f, 0x4c, 0x69, 0x73, 0x74, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x4f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x41, 0x0a, 0x13,
	0x50, 0x6f, 0x6f, 0x6c, 0x49, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x2a, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x61, 0x79, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x49, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x22,
	0xe3, 0x05, 0x0a, 0x07, 0x49, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0c, 0x6e, 0x75, 0x6d, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6f, 0x70, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0a, 0x6e, 0x75, 0x6d, 0x52, 0x65, 0x61, 0x64, 0x4f, 0x70,
	0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x62, 0x79, 0x74, 0x65, 0x73, 0x52, 0x65, 0x61, 0x64,
	0x12, 0x22, 0x0a, 0x0d, 0x6e, 0x75, 0x6d, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6f, 0x70,
	0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x57, 0x72, 0x69, 0x74,
	0x65, 0x4f, 0x70, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x77, 0x72,
	0x69, 0x74, 0x74, 0x65, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x62, 0x79, 0x74,
	0x65, 0x73, 0x57, 0x72, 0x69, 0x74, 0x74, 0x65, 0x6e, 0x12, 0x22, 0x0a, 0x0d, 0x6e, 0x75, 0x6d,
	0x5f, 0x75, 0x6e, 0x6d, 0x61, 0x70, 0x5f, 0x6f, 0x70, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x55, 0x6e, 0x6d, 0x61, 0x70, 0x4f, 0x70, 0x73, 0x12, 0x25, 0x0a,
	0x0e, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x75, 0x6e, 0x6d, 0x61, 0x70, 0x70, 0x65, 0x64, 0x18,
	0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0d, 0x62, 0x79, 0x74, 0x65, 0x73, 0x55, 0x6e, 0x6d, 0x61,
	0x70, 0x70, 0x65, 0x64, 0x12, 0x2c, 0x0a, 0x12, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6c, 0x61, 0x74,
	0x65, 0x6e, 0x63, 0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x10, 0x72, 0x65, 0x61, 0x64, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63,
	0x6b, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6c, 0x61, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x11, 0x77, 0x72, 0x69, 0x74, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63,
	0x6b, 0x73, 0x12, 0x2e, 0x0a, 0x13, 0x75, 0x6e, 0x6d, 0x61, 0x70, 0x5f, 0x6c, 0x61, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x11, 0x75, 0x6e, 0x6d, 0x61, 0x70, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63,
	0x6b, 0x73, 0x12, 0x33, 0x0a, 0x16, 0x6d, 0x61, 0x78, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x6c,
	0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x13, 0x6d, 0x61, 0x78, 0x52, 0x65, 0x61, 0x64, 0x4c, 0x61, 0x74, 0x65, 0x6e,
	0x63, 0x79, 0x54, 0x69, 0x63, 0x6b, 0x73, 0x12, 0x33, 0x0a, 0x16, 0x6d, 0x69, 0x6e, 0x5f, 0x72,
	0x65, 0x61, 0x64, 0x5f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b,
	0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x04, 0x52, 0x13, 0x6d, 0x69, 0x6e, 0x52, 0x65, 0x61, 0x64,
	0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63, 0x6b, 0x73, 0x12, 0x35, 0x0a, 0x17,
	0x6d, 0x61, 0x78, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65, 0x5f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63,
	0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x04, 0x52, 0x14, 0x6d,
	0x61, 0x78, 0x57, 0x72, 0x69, 0x74, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69,
	0x63, 0x6b, 0x73, 0x12, 0x35, 0x0a, 0x17, 0x6d, 0x69, 0x6e, 0x5f, 0x77, 0x72, 0x69, 0x74, 0x65,
	0x5f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x0e,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x14, 0x6d, 0x69, 0x6e, 0x57, 0x72, 0x69, 0x74, 0x65, 0x4c, 0x61,
	0x74, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63, 0x6b, 0x73, 0x12, 0x35, 0x0a, 0x17, 0x6d, 0x61,
	0x78, 0x5f, 0x75, 0x6e, 0x6d, 0x61, 0x70, 0x5f, 0x6c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f,
	0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x04, 0x52, 0x14, 0x6d, 0x61, 0x78,
	0x55, 0x6e, 0x6d, 0x61, 0x70, 0x4c, 0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x54, 0x69, 0x63, 0x6b,
	0x73, 0x12, 0x35, 0x0a, 0x17, 0x6d, 0x69, 0x6e, 0x5f, 0x75, 0x6e, 0x6d, 0x61, 0x70, 0x5f, 0x6c,
	0x61, 0x74, 0x65, 0x6e, 0x63, 0x79, 0x5f, 0x74, 0x69, 0x63, 0x6b, 0x73, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x14, 0x6d, 0x69, 0x6e, 0x55, 0x6e, 0x6d, 0x61, 0x70, 0x4c, 0x61, 0x74, 0x65,
	0x6e, 0x63, 0x79, 0x54, 0x69, 0x63, 0x6b, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x69, 0x63, 0x6b,
	0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x11, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x74, 0x69, 0x63,
	0x6b, 0x52, 0x61, 0x74, 0x65, 0x32, 0x9c, 0x01, 0x0a, 0x08, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52,
	0x70, 0x63, 0x12, 0x50, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x50, 0x6f, 0x6f, 0x6c, 0x49, 0x6f, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x12, 0x1c, 0x2e, 0x6d, 0x61, 0x79, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x74, 0x61, 0x74, 0x73, 0x4f, 0x70, 0x74, 0x69,
	0x6f, 0x6e, 0x1a, 0x20, 0x2e, 0x6d, 0x61, 0x79, 0x61, 0x73, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x50, 0x6f, 0x6f, 0x6c, 0x49, 0x6f, 0x53, 0x74, 0x61, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x65, 0x74, 0x49, 0x6f, 0x53,
	0x74, 0x61, 0x74, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45,
	0x6d, 0x70, 0x74, 0x79, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x6f, 0x70, 0x65, 0x6e, 0x65, 0x62, 0x73, 0x2f, 0x6d, 0x61, 0x79, 0x61, 0x73,
	0x74, 0x6f, 0x72, 0x2d, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_stats_proto_rawDescOnce sync.Once
	file_stats_proto_rawDescData = file_stats_proto_rawDesc
)

func file_stats_proto_rawDescGZIP() []byte {
	file_stats_proto_rawDescOnce.Do(func() {
		file_stats_proto_rawDescData = protoimpl.X.CompressGZIP(file_stats_proto_rawDescData)
	})
	return file_stats_proto_rawDescData
}

var file_stats_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_stats_proto_goTypes = []interface{}{
	(*ListStatsOption)(nil),     // 0: mayastor.v1.ListStatsOption
	(*PoolIoStatsResponse)(nil), // 1: mayastor.v1.PoolIoStatsResponse
	(*IoStats)(nil),             // 2: mayastor.v1.IoStats
	(*emptypb.Empty)(nil),       // 3: google.protobuf.Empty
}
var file_stats_proto_depIdxs = []int32{
	2, // 0: mayastor.v1.PoolIoStatsResponse.stats:type_name -> mayastor.v1.IoStats
	0, // 1: mayastor.v1.StatsRpc.GetPoolIoStats:input_type -> mayastor.v1.ListStatsOption
	3, // 2: mayastor.v1.StatsRpc.ResetIoStats:input_type -> google.protobuf.Empty
	1, // 3: mayastor.v1.StatsRpc.GetPoolIoStats:output_type -> mayastor.v1.PoolIoStatsResponse
	3, // 4: mayastor.v1.StatsRpc.ResetIoStats:output_type -> google.protobuf.Empty
	3, // [3:5] is the sub-list for method output_type
	1, // [1:3] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_stats_proto_init() }
func file_stats_proto_init() {
	if File_stats_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_stats_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListStatsOption); i {
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
		file_stats_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PoolIoStatsResponse); i {
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
		file_stats_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IoStats); i {
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
	file_stats_proto_msgTypes[0].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_stats_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_stats_proto_goTypes,
		DependencyIndexes: file_stats_proto_depIdxs,
		MessageInfos:      file_stats_proto_msgTypes,
	}.Build()
	File_stats_proto = out.File
	file_stats_proto_rawDesc = nil
	file_stats_proto_goTypes = nil
	file_stats_proto_depIdxs = nil
}
