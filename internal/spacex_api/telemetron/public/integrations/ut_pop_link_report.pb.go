// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: spacex_api/telemetron/public/integrations/ut_pop_link_report.proto

package utpoplink

import (
	common "github.com/joshuasing/starlink_exporter/internal/spacex_api/telemetron/public/common"
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

type RateLimitReason int32

const (
	RateLimitReason_UNKNOWN           RateLimitReason = 0
	RateLimitReason_NO_LIMIT          RateLimitReason = 1
	RateLimitReason_POLICY_LIMIT      RateLimitReason = 2
	RateLimitReason_USER_CUSTOM_LIMIT RateLimitReason = 3
	RateLimitReason_OVERAGE_LIMIT     RateLimitReason = 5
)

// Enum value maps for RateLimitReason.
var (
	RateLimitReason_name = map[int32]string{
		0: "UNKNOWN",
		1: "NO_LIMIT",
		2: "POLICY_LIMIT",
		3: "USER_CUSTOM_LIMIT",
		5: "OVERAGE_LIMIT",
	}
	RateLimitReason_value = map[string]int32{
		"UNKNOWN":           0,
		"NO_LIMIT":          1,
		"POLICY_LIMIT":      2,
		"USER_CUSTOM_LIMIT": 3,
		"OVERAGE_LIMIT":     5,
	}
)

func (x RateLimitReason) Enum() *RateLimitReason {
	p := new(RateLimitReason)
	*p = x
	return p
}

func (x RateLimitReason) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RateLimitReason) Descriptor() protoreflect.EnumDescriptor {
	return file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_enumTypes[0].Descriptor()
}

func (RateLimitReason) Type() protoreflect.EnumType {
	return &file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_enumTypes[0]
}

func (x RateLimitReason) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RateLimitReason.Descriptor instead.
func (RateLimitReason) EnumDescriptor() ([]byte, []int) {
	return file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescGZIP(), []int{0}
}

type UtPoPLinkReport struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SlotTimestamp *common.TimestampInfo  `protobuf:"bytes,1,opt,name=slot_timestamp,json=slotTimestamp,proto3" json:"slot_timestamp,omitempty"`
	PopId         uint32                 `protobuf:"varint,2,opt,name=pop_id,json=popId,proto3" json:"pop_id,omitempty"`
	// Deprecated: Marked as deprecated in spacex_api/telemetron/public/integrations/ut_pop_link_report.proto.
	PopRackId     uint32            `protobuf:"varint,3,opt,name=pop_rack_id,json=popRackId,proto3" json:"pop_rack_id,omitempty"`
	Stats         []*UtPoPLinkStats `protobuf:"bytes,4,rep,name=stats,proto3" json:"stats,omitempty"`
	PopVersion    string            `protobuf:"bytes,5,opt,name=pop_version,json=popVersion,proto3" json:"pop_version,omitempty"`
	InstanceIndex uint32            `protobuf:"varint,6,opt,name=instance_index,json=instanceIndex,proto3" json:"instance_index,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UtPoPLinkReport) Reset() {
	*x = UtPoPLinkReport{}
	mi := &file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UtPoPLinkReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UtPoPLinkReport) ProtoMessage() {}

func (x *UtPoPLinkReport) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UtPoPLinkReport.ProtoReflect.Descriptor instead.
func (*UtPoPLinkReport) Descriptor() ([]byte, []int) {
	return file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescGZIP(), []int{0}
}

func (x *UtPoPLinkReport) GetSlotTimestamp() *common.TimestampInfo {
	if x != nil {
		return x.SlotTimestamp
	}
	return nil
}

func (x *UtPoPLinkReport) GetPopId() uint32 {
	if x != nil {
		return x.PopId
	}
	return 0
}

// Deprecated: Marked as deprecated in spacex_api/telemetron/public/integrations/ut_pop_link_report.proto.
func (x *UtPoPLinkReport) GetPopRackId() uint32 {
	if x != nil {
		return x.PopRackId
	}
	return 0
}

func (x *UtPoPLinkReport) GetStats() []*UtPoPLinkStats {
	if x != nil {
		return x.Stats
	}
	return nil
}

func (x *UtPoPLinkReport) GetPopVersion() string {
	if x != nil {
		return x.PopVersion
	}
	return ""
}

func (x *UtPoPLinkReport) GetInstanceIndex() uint32 {
	if x != nil {
		return x.InstanceIndex
	}
	return 0
}

type UtPoPLinkStats struct {
	state                                  protoimpl.MessageState `protogen:"open.v1"`
	MeasurementTimestamp                   *common.TimestampInfo  `protobuf:"bytes,1,opt,name=measurement_timestamp,json=measurementTimestamp,proto3" json:"measurement_timestamp,omitempty"`
	UtId                                   string                 `protobuf:"bytes,2,opt,name=ut_id,json=utId,proto3" json:"ut_id,omitempty"`
	PopRxSduCnt                            int64                  `protobuf:"varint,3,opt,name=pop_rx_sdu_cnt,json=popRxSduCnt,proto3" json:"pop_rx_sdu_cnt,omitempty"`
	SduLossCnt                             int64                  `protobuf:"varint,4,opt,name=sdu_loss_cnt,json=sduLossCnt,proto3" json:"sdu_loss_cnt,omitempty"`
	UplinkBytesLast_15S                    uint64                 `protobuf:"varint,5,opt,name=uplink_bytes_last_15s,json=uplinkBytesLast15s,proto3" json:"uplink_bytes_last_15s,omitempty"`
	DownlinkBytesLast_15S                  uint64                 `protobuf:"varint,6,opt,name=downlink_bytes_last_15s,json=downlinkBytesLast15s,proto3" json:"downlink_bytes_last_15s,omitempty"`
	UplinkCplaneAclOtherViolationsLast_15S uint64                 `protobuf:"varint,7,opt,name=uplink_cplane_acl_other_violations_last_15s,json=uplinkCplaneAclOtherViolationsLast15s,proto3" json:"uplink_cplane_acl_other_violations_last_15s,omitempty"`
	unknownFields                          protoimpl.UnknownFields
	sizeCache                              protoimpl.SizeCache
}

func (x *UtPoPLinkStats) Reset() {
	*x = UtPoPLinkStats{}
	mi := &file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UtPoPLinkStats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UtPoPLinkStats) ProtoMessage() {}

func (x *UtPoPLinkStats) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UtPoPLinkStats.ProtoReflect.Descriptor instead.
func (*UtPoPLinkStats) Descriptor() ([]byte, []int) {
	return file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescGZIP(), []int{1}
}

func (x *UtPoPLinkStats) GetMeasurementTimestamp() *common.TimestampInfo {
	if x != nil {
		return x.MeasurementTimestamp
	}
	return nil
}

func (x *UtPoPLinkStats) GetUtId() string {
	if x != nil {
		return x.UtId
	}
	return ""
}

func (x *UtPoPLinkStats) GetPopRxSduCnt() int64 {
	if x != nil {
		return x.PopRxSduCnt
	}
	return 0
}

func (x *UtPoPLinkStats) GetSduLossCnt() int64 {
	if x != nil {
		return x.SduLossCnt
	}
	return 0
}

func (x *UtPoPLinkStats) GetUplinkBytesLast_15S() uint64 {
	if x != nil {
		return x.UplinkBytesLast_15S
	}
	return 0
}

func (x *UtPoPLinkStats) GetDownlinkBytesLast_15S() uint64 {
	if x != nil {
		return x.DownlinkBytesLast_15S
	}
	return 0
}

func (x *UtPoPLinkStats) GetUplinkCplaneAclOtherViolationsLast_15S() uint64 {
	if x != nil {
		return x.UplinkCplaneAclOtherViolationsLast_15S
	}
	return 0
}

var File_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto protoreflect.FileDescriptor

const file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDesc = "" +
	"\n" +
	"Bspacex_api/telemetron/public/integrations/ut_pop_link_report.proto\x12)SpaceX.API.Telemetron.Public.Integrations\x1a.spacex_api/telemetron/public/common/time.proto\"\xc0\x02\n" +
	"\x0fUtPoPLinkReport\x12Y\n" +
	"\x0eslot_timestamp\x18\x01 \x01(\v22.SpaceX.API.Telemetron.Public.Common.TimestampInfoR\rslotTimestamp\x12\x15\n" +
	"\x06pop_id\x18\x02 \x01(\rR\x05popId\x12\"\n" +
	"\vpop_rack_id\x18\x03 \x01(\rB\x02\x18\x01R\tpopRackId\x12O\n" +
	"\x05stats\x18\x04 \x03(\v29.SpaceX.API.Telemetron.Public.Integrations.UtPoPLinkStatsR\x05stats\x12\x1f\n" +
	"\vpop_version\x18\x05 \x01(\tR\n" +
	"popVersion\x12%\n" +
	"\x0einstance_index\x18\x06 \x01(\rR\rinstanceIndex\"\x9b\x03\n" +
	"\x0eUtPoPLinkStats\x12g\n" +
	"\x15measurement_timestamp\x18\x01 \x01(\v22.SpaceX.API.Telemetron.Public.Common.TimestampInfoR\x14measurementTimestamp\x12\x13\n" +
	"\x05ut_id\x18\x02 \x01(\tR\x04utId\x12#\n" +
	"\x0epop_rx_sdu_cnt\x18\x03 \x01(\x03R\vpopRxSduCnt\x12 \n" +
	"\fsdu_loss_cnt\x18\x04 \x01(\x03R\n" +
	"sduLossCnt\x121\n" +
	"\x15uplink_bytes_last_15s\x18\x05 \x01(\x04R\x12uplinkBytesLast15s\x125\n" +
	"\x17downlink_bytes_last_15s\x18\x06 \x01(\x04R\x14downlinkBytesLast15s\x12Z\n" +
	"+uplink_cplane_acl_other_violations_last_15s\x18\a \x01(\x04R%uplinkCplaneAclOtherViolationsLast15s*\x85\x01\n" +
	"\x0fRateLimitReason\x12\v\n" +
	"\aUNKNOWN\x10\x00\x12\f\n" +
	"\bNO_LIMIT\x10\x01\x12\x10\n" +
	"\fPOLICY_LIMIT\x10\x02\x12\x15\n" +
	"\x11USER_CUSTOM_LIMIT\x10\x03\x12\x11\n" +
	"\rOVERAGE_LIMIT\x10\x05\"\x04\b\x04\x10\x04*\x15HIGH_HOURLY_AVG_LIMITB9Z7spacex.com/api/telemetron/public/integrations/utpoplinkb\x06proto3"

var (
	file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescOnce sync.Once
	file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescData []byte
)

func file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescGZIP() []byte {
	file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescOnce.Do(func() {
		file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDesc), len(file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDesc)))
	})
	return file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDescData
}

var file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_goTypes = []any{
	(RateLimitReason)(0),         // 0: SpaceX.API.Telemetron.Public.Integrations.RateLimitReason
	(*UtPoPLinkReport)(nil),      // 1: SpaceX.API.Telemetron.Public.Integrations.UtPoPLinkReport
	(*UtPoPLinkStats)(nil),       // 2: SpaceX.API.Telemetron.Public.Integrations.UtPoPLinkStats
	(*common.TimestampInfo)(nil), // 3: SpaceX.API.Telemetron.Public.Common.TimestampInfo
}
var file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_depIdxs = []int32{
	3, // 0: SpaceX.API.Telemetron.Public.Integrations.UtPoPLinkReport.slot_timestamp:type_name -> SpaceX.API.Telemetron.Public.Common.TimestampInfo
	2, // 1: SpaceX.API.Telemetron.Public.Integrations.UtPoPLinkReport.stats:type_name -> SpaceX.API.Telemetron.Public.Integrations.UtPoPLinkStats
	3, // 2: SpaceX.API.Telemetron.Public.Integrations.UtPoPLinkStats.measurement_timestamp:type_name -> SpaceX.API.Telemetron.Public.Common.TimestampInfo
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_init() }
func file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_init() {
	if File_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDesc), len(file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_rawDesc)),
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_goTypes,
		DependencyIndexes: file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_depIdxs,
		EnumInfos:         file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_enumTypes,
		MessageInfos:      file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_msgTypes,
	}.Build()
	File_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto = out.File
	file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_goTypes = nil
	file_spacex_api_telemetron_public_integrations_ut_pop_link_report_proto_depIdxs = nil
}
