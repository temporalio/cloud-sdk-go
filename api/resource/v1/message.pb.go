// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.4
// 	protoc        (unknown)
// source: temporal/api/cloud/resource/v1/message.proto

package resourcev1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ResourceState int32

const (
	ResourceState_RESOURCE_STATE_UNSPECIFIED       ResourceState = 0
	ResourceState_RESOURCE_STATE_ACTIVATING        ResourceState = 1  // The resource is being activated.
	ResourceState_RESOURCE_STATE_ACTIVATION_FAILED ResourceState = 2  // The resource failed to activate. This is an error state. Reach out to support for remediation.
	ResourceState_RESOURCE_STATE_ACTIVE            ResourceState = 3  // The resource is active and ready to use.
	ResourceState_RESOURCE_STATE_UPDATING          ResourceState = 4  // The resource is being updated.
	ResourceState_RESOURCE_STATE_UPDATE_FAILED     ResourceState = 5  // The resource failed to update. This is an error state. Reach out to support for remediation.
	ResourceState_RESOURCE_STATE_DELETING          ResourceState = 6  // The resource is being deleted.
	ResourceState_RESOURCE_STATE_DELETE_FAILED     ResourceState = 7  // The resource failed to delete. This is an error state. Reach out to support for remediation.
	ResourceState_RESOURCE_STATE_DELETED           ResourceState = 8  // The resource has been deleted.
	ResourceState_RESOURCE_STATE_SUSPENDED         ResourceState = 9  // The resource is suspended and not available for use. Reach out to support for remediation.
	ResourceState_RESOURCE_STATE_EXPIRED           ResourceState = 10 // The resource has expired and is no longer available for use.
)

// Enum value maps for ResourceState.
var (
	ResourceState_name = map[int32]string{
		0:  "RESOURCE_STATE_UNSPECIFIED",
		1:  "RESOURCE_STATE_ACTIVATING",
		2:  "RESOURCE_STATE_ACTIVATION_FAILED",
		3:  "RESOURCE_STATE_ACTIVE",
		4:  "RESOURCE_STATE_UPDATING",
		5:  "RESOURCE_STATE_UPDATE_FAILED",
		6:  "RESOURCE_STATE_DELETING",
		7:  "RESOURCE_STATE_DELETE_FAILED",
		8:  "RESOURCE_STATE_DELETED",
		9:  "RESOURCE_STATE_SUSPENDED",
		10: "RESOURCE_STATE_EXPIRED",
	}
	ResourceState_value = map[string]int32{
		"RESOURCE_STATE_UNSPECIFIED":       0,
		"RESOURCE_STATE_ACTIVATING":        1,
		"RESOURCE_STATE_ACTIVATION_FAILED": 2,
		"RESOURCE_STATE_ACTIVE":            3,
		"RESOURCE_STATE_UPDATING":          4,
		"RESOURCE_STATE_UPDATE_FAILED":     5,
		"RESOURCE_STATE_DELETING":          6,
		"RESOURCE_STATE_DELETE_FAILED":     7,
		"RESOURCE_STATE_DELETED":           8,
		"RESOURCE_STATE_SUSPENDED":         9,
		"RESOURCE_STATE_EXPIRED":           10,
	}
)

func (x ResourceState) Enum() *ResourceState {
	p := new(ResourceState)
	*p = x
	return p
}

func (x ResourceState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResourceState) Descriptor() protoreflect.EnumDescriptor {
	return file_temporal_api_cloud_resource_v1_message_proto_enumTypes[0].Descriptor()
}

func (ResourceState) Type() protoreflect.EnumType {
	return &file_temporal_api_cloud_resource_v1_message_proto_enumTypes[0]
}

func (x ResourceState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResourceState.Descriptor instead.
func (ResourceState) EnumDescriptor() ([]byte, []int) {
	return file_temporal_api_cloud_resource_v1_message_proto_rawDescGZIP(), []int{0}
}

var File_temporal_api_cloud_resource_v1_message_proto protoreflect.FileDescriptor

var file_temporal_api_cloud_resource_v1_message_proto_rawDesc = string([]byte{
	0x0a, 0x2c, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63,
	0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x76, 0x31,
	0x2f, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1e,
	0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x6f,
	0x75, 0x64, 0x2e, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x2a, 0xe3,
	0x02, 0x0a, 0x0d, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65,
	0x12, 0x1e, 0x0a, 0x1a, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00,
	0x12, 0x1d, 0x0a, 0x19, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x41, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12,
	0x24, 0x0a, 0x20, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x45, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x41, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x46, 0x41, 0x49,
	0x4c, 0x45, 0x44, 0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43,
	0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x03,
	0x12, 0x1b, 0x0a, 0x17, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x45, 0x5f, 0x55, 0x50, 0x44, 0x41, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x20, 0x0a,
	0x1c, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f,
	0x55, 0x50, 0x44, 0x41, 0x54, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x05, 0x12,
	0x1b, 0x0a, 0x17, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x45, 0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x06, 0x12, 0x20, 0x0a, 0x1c,
	0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x44,
	0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x07, 0x12, 0x1a,
	0x0a, 0x16, 0x52, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45,
	0x5f, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x44, 0x10, 0x08, 0x12, 0x1c, 0x0a, 0x18, 0x52, 0x45,
	0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x55, 0x53,
	0x50, 0x45, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x09, 0x12, 0x1a, 0x0a, 0x16, 0x52, 0x45, 0x53, 0x4f,
	0x55, 0x52, 0x43, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x45, 0x58, 0x50, 0x49, 0x52,
	0x45, 0x44, 0x10, 0x0a, 0x42, 0x97, 0x02, 0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x65, 0x6d,
	0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2e,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x76, 0x31, 0x42, 0x0c, 0x4d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x46, 0x67, 0x6f, 0x2e,
	0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x69, 0x6f, 0x2f, 0x63, 0x6c, 0x6f, 0x75,
	0x64, 0x2d, 0x73, 0x64, 0x6b, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6f, 0x72,
	0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x2f, 0x72, 0x65, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x2f, 0x76, 0x31, 0x3b, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x76, 0x31, 0xa2, 0x02, 0x04, 0x54, 0x41, 0x43, 0x52, 0xaa, 0x02, 0x1e, 0x54, 0x65, 0x6d,
	0x70, 0x6f, 0x72, 0x61, 0x6c, 0x2e, 0x41, 0x70, 0x69, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x2e,
	0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x1e, 0x54, 0x65,
	0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x5c, 0x41, 0x70, 0x69, 0x5c, 0x43, 0x6c, 0x6f, 0x75, 0x64,
	0x5c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x2a, 0x54,
	0x65, 0x6d, 0x70, 0x6f, 0x72, 0x61, 0x6c, 0x5c, 0x41, 0x70, 0x69, 0x5c, 0x43, 0x6c, 0x6f, 0x75,
	0x64, 0x5c, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x22, 0x54, 0x65, 0x6d, 0x70,
	0x6f, 0x72, 0x61, 0x6c, 0x3a, 0x3a, 0x41, 0x70, 0x69, 0x3a, 0x3a, 0x43, 0x6c, 0x6f, 0x75, 0x64,
	0x3a, 0x3a, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_temporal_api_cloud_resource_v1_message_proto_rawDescOnce sync.Once
	file_temporal_api_cloud_resource_v1_message_proto_rawDescData []byte
)

func file_temporal_api_cloud_resource_v1_message_proto_rawDescGZIP() []byte {
	file_temporal_api_cloud_resource_v1_message_proto_rawDescOnce.Do(func() {
		file_temporal_api_cloud_resource_v1_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_temporal_api_cloud_resource_v1_message_proto_rawDesc), len(file_temporal_api_cloud_resource_v1_message_proto_rawDesc)))
	})
	return file_temporal_api_cloud_resource_v1_message_proto_rawDescData
}

var file_temporal_api_cloud_resource_v1_message_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_temporal_api_cloud_resource_v1_message_proto_goTypes = []any{
	(ResourceState)(0), // 0: temporal.api.cloud.resource.v1.ResourceState
}
var file_temporal_api_cloud_resource_v1_message_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_temporal_api_cloud_resource_v1_message_proto_init() }
func file_temporal_api_cloud_resource_v1_message_proto_init() {
	if File_temporal_api_cloud_resource_v1_message_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_temporal_api_cloud_resource_v1_message_proto_rawDesc), len(file_temporal_api_cloud_resource_v1_message_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_temporal_api_cloud_resource_v1_message_proto_goTypes,
		DependencyIndexes: file_temporal_api_cloud_resource_v1_message_proto_depIdxs,
		EnumInfos:         file_temporal_api_cloud_resource_v1_message_proto_enumTypes,
	}.Build()
	File_temporal_api_cloud_resource_v1_message_proto = out.File
	file_temporal_api_cloud_resource_v1_message_proto_goTypes = nil
	file_temporal_api_cloud_resource_v1_message_proto_depIdxs = nil
}
