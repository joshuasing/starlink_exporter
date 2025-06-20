// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: spacex_api/device/transceiver.proto

package device

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

type TransceiverModulatorState int32

const (
	TransceiverModulatorState_MODSTATE_UNKNOWN  TransceiverModulatorState = 0
	TransceiverModulatorState_MODSTATE_ENABLED  TransceiverModulatorState = 1
	TransceiverModulatorState_MODSTATE_DISABLED TransceiverModulatorState = 2
)

// Enum value maps for TransceiverModulatorState.
var (
	TransceiverModulatorState_name = map[int32]string{
		0: "MODSTATE_UNKNOWN",
		1: "MODSTATE_ENABLED",
		2: "MODSTATE_DISABLED",
	}
	TransceiverModulatorState_value = map[string]int32{
		"MODSTATE_UNKNOWN":  0,
		"MODSTATE_ENABLED":  1,
		"MODSTATE_DISABLED": 2,
	}
)

func (x TransceiverModulatorState) Enum() *TransceiverModulatorState {
	p := new(TransceiverModulatorState)
	*p = x
	return p
}

func (x TransceiverModulatorState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransceiverModulatorState) Descriptor() protoreflect.EnumDescriptor {
	return file_spacex_api_device_transceiver_proto_enumTypes[0].Descriptor()
}

func (TransceiverModulatorState) Type() protoreflect.EnumType {
	return &file_spacex_api_device_transceiver_proto_enumTypes[0]
}

func (x TransceiverModulatorState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransceiverModulatorState.Descriptor instead.
func (TransceiverModulatorState) EnumDescriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{0}
}

type TransceiverTxRxState int32

const (
	TransceiverTxRxState_TXRX_UNKNOWN  TransceiverTxRxState = 0
	TransceiverTxRxState_TXRX_ENABLED  TransceiverTxRxState = 1
	TransceiverTxRxState_TXRX_DISABLED TransceiverTxRxState = 2
)

// Enum value maps for TransceiverTxRxState.
var (
	TransceiverTxRxState_name = map[int32]string{
		0: "TXRX_UNKNOWN",
		1: "TXRX_ENABLED",
		2: "TXRX_DISABLED",
	}
	TransceiverTxRxState_value = map[string]int32{
		"TXRX_UNKNOWN":  0,
		"TXRX_ENABLED":  1,
		"TXRX_DISABLED": 2,
	}
)

func (x TransceiverTxRxState) Enum() *TransceiverTxRxState {
	p := new(TransceiverTxRxState)
	*p = x
	return p
}

func (x TransceiverTxRxState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransceiverTxRxState) Descriptor() protoreflect.EnumDescriptor {
	return file_spacex_api_device_transceiver_proto_enumTypes[1].Descriptor()
}

func (TransceiverTxRxState) Type() protoreflect.EnumType {
	return &file_spacex_api_device_transceiver_proto_enumTypes[1]
}

func (x TransceiverTxRxState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransceiverTxRxState.Descriptor instead.
func (TransceiverTxRxState) EnumDescriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{1}
}

type TransceiverTransmitBlankingState int32

const (
	TransceiverTransmitBlankingState_TB_UNKNOWN  TransceiverTransmitBlankingState = 0
	TransceiverTransmitBlankingState_TB_ENABLED  TransceiverTransmitBlankingState = 1
	TransceiverTransmitBlankingState_TB_DISABLED TransceiverTransmitBlankingState = 2
)

// Enum value maps for TransceiverTransmitBlankingState.
var (
	TransceiverTransmitBlankingState_name = map[int32]string{
		0: "TB_UNKNOWN",
		1: "TB_ENABLED",
		2: "TB_DISABLED",
	}
	TransceiverTransmitBlankingState_value = map[string]int32{
		"TB_UNKNOWN":  0,
		"TB_ENABLED":  1,
		"TB_DISABLED": 2,
	}
)

func (x TransceiverTransmitBlankingState) Enum() *TransceiverTransmitBlankingState {
	p := new(TransceiverTransmitBlankingState)
	*p = x
	return p
}

func (x TransceiverTransmitBlankingState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TransceiverTransmitBlankingState) Descriptor() protoreflect.EnumDescriptor {
	return file_spacex_api_device_transceiver_proto_enumTypes[2].Descriptor()
}

func (TransceiverTransmitBlankingState) Type() protoreflect.EnumType {
	return &file_spacex_api_device_transceiver_proto_enumTypes[2]
}

func (x TransceiverTransmitBlankingState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TransceiverTransmitBlankingState.Descriptor instead.
func (TransceiverTransmitBlankingState) EnumDescriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{2}
}

type TransceiverIFLoopbackTestRequest struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	EnableIfLoopback bool                   `protobuf:"varint,1,opt,name=enable_if_loopback,json=enableIfLoopback,proto3" json:"enable_if_loopback,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *TransceiverIFLoopbackTestRequest) Reset() {
	*x = TransceiverIFLoopbackTestRequest{}
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransceiverIFLoopbackTestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransceiverIFLoopbackTestRequest) ProtoMessage() {}

func (x *TransceiverIFLoopbackTestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransceiverIFLoopbackTestRequest.ProtoReflect.Descriptor instead.
func (*TransceiverIFLoopbackTestRequest) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{0}
}

func (x *TransceiverIFLoopbackTestRequest) GetEnableIfLoopback() bool {
	if x != nil {
		return x.EnableIfLoopback
	}
	return false
}

type TransceiverIFLoopbackTestResponse struct {
	state            protoimpl.MessageState `protogen:"open.v1"`
	BerLoopbackTest  float32                `protobuf:"fixed32,1,opt,name=ber_loopback_test,json=berLoopbackTest,proto3" json:"ber_loopback_test,omitempty"`
	SnrLoopbackTest  float32                `protobuf:"fixed32,2,opt,name=snr_loopback_test,json=snrLoopbackTest,proto3" json:"snr_loopback_test,omitempty"`
	RssiLoopbackTest float32                `protobuf:"fixed32,3,opt,name=rssi_loopback_test,json=rssiLoopbackTest,proto3" json:"rssi_loopback_test,omitempty"`
	PllLock          bool                   `protobuf:"varint,4,opt,name=pll_lock,json=pllLock,proto3" json:"pll_lock,omitempty"`
	unknownFields    protoimpl.UnknownFields
	sizeCache        protoimpl.SizeCache
}

func (x *TransceiverIFLoopbackTestResponse) Reset() {
	*x = TransceiverIFLoopbackTestResponse{}
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransceiverIFLoopbackTestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransceiverIFLoopbackTestResponse) ProtoMessage() {}

func (x *TransceiverIFLoopbackTestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransceiverIFLoopbackTestResponse.ProtoReflect.Descriptor instead.
func (*TransceiverIFLoopbackTestResponse) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{1}
}

func (x *TransceiverIFLoopbackTestResponse) GetBerLoopbackTest() float32 {
	if x != nil {
		return x.BerLoopbackTest
	}
	return 0
}

func (x *TransceiverIFLoopbackTestResponse) GetSnrLoopbackTest() float32 {
	if x != nil {
		return x.SnrLoopbackTest
	}
	return 0
}

func (x *TransceiverIFLoopbackTestResponse) GetRssiLoopbackTest() float32 {
	if x != nil {
		return x.RssiLoopbackTest
	}
	return 0
}

func (x *TransceiverIFLoopbackTestResponse) GetPllLock() bool {
	if x != nil {
		return x.PllLock
	}
	return false
}

type TransceiverGetStatusRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransceiverGetStatusRequest) Reset() {
	*x = TransceiverGetStatusRequest{}
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransceiverGetStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransceiverGetStatusRequest) ProtoMessage() {}

func (x *TransceiverGetStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransceiverGetStatusRequest.ProtoReflect.Descriptor instead.
func (*TransceiverGetStatusRequest) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{2}
}

type TransceiverGetStatusResponse struct {
	state                 protoimpl.MessageState           `protogen:"open.v1"`
	ModState              TransceiverModulatorState        `protobuf:"varint,1,opt,name=mod_state,json=modState,proto3,enum=SpaceX.API.Device.TransceiverModulatorState" json:"mod_state,omitempty"`
	DemodState            TransceiverModulatorState        `protobuf:"varint,2,opt,name=demod_state,json=demodState,proto3,enum=SpaceX.API.Device.TransceiverModulatorState" json:"demod_state,omitempty"`
	TxState               TransceiverTxRxState             `protobuf:"varint,3,opt,name=tx_state,json=txState,proto3,enum=SpaceX.API.Device.TransceiverTxRxState" json:"tx_state,omitempty"`
	RxState               TransceiverTxRxState             `protobuf:"varint,4,opt,name=rx_state,json=rxState,proto3,enum=SpaceX.API.Device.TransceiverTxRxState" json:"rx_state,omitempty"`
	State                 DishState                        `protobuf:"varint,1006,opt,name=state,proto3,enum=SpaceX.API.Device.DishState" json:"state,omitempty"`
	Faults                *TransceiverFaults               `protobuf:"bytes,1007,opt,name=faults,proto3" json:"faults,omitempty"`
	TransmitBlankingState TransceiverTransmitBlankingState `protobuf:"varint,1008,opt,name=transmit_blanking_state,json=transmitBlankingState,proto3,enum=SpaceX.API.Device.TransceiverTransmitBlankingState" json:"transmit_blanking_state,omitempty"`
	ModemAsicTemp         float32                          `protobuf:"fixed32,1009,opt,name=modem_asic_temp,json=modemAsicTemp,proto3" json:"modem_asic_temp,omitempty"`
	TxIfTemp              float32                          `protobuf:"fixed32,1010,opt,name=tx_if_temp,json=txIfTemp,proto3" json:"tx_if_temp,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *TransceiverGetStatusResponse) Reset() {
	*x = TransceiverGetStatusResponse{}
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransceiverGetStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransceiverGetStatusResponse) ProtoMessage() {}

func (x *TransceiverGetStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransceiverGetStatusResponse.ProtoReflect.Descriptor instead.
func (*TransceiverGetStatusResponse) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{3}
}

func (x *TransceiverGetStatusResponse) GetModState() TransceiverModulatorState {
	if x != nil {
		return x.ModState
	}
	return TransceiverModulatorState_MODSTATE_UNKNOWN
}

func (x *TransceiverGetStatusResponse) GetDemodState() TransceiverModulatorState {
	if x != nil {
		return x.DemodState
	}
	return TransceiverModulatorState_MODSTATE_UNKNOWN
}

func (x *TransceiverGetStatusResponse) GetTxState() TransceiverTxRxState {
	if x != nil {
		return x.TxState
	}
	return TransceiverTxRxState_TXRX_UNKNOWN
}

func (x *TransceiverGetStatusResponse) GetRxState() TransceiverTxRxState {
	if x != nil {
		return x.RxState
	}
	return TransceiverTxRxState_TXRX_UNKNOWN
}

func (x *TransceiverGetStatusResponse) GetState() DishState {
	if x != nil {
		return x.State
	}
	return DishState_UNKNOWN
}

func (x *TransceiverGetStatusResponse) GetFaults() *TransceiverFaults {
	if x != nil {
		return x.Faults
	}
	return nil
}

func (x *TransceiverGetStatusResponse) GetTransmitBlankingState() TransceiverTransmitBlankingState {
	if x != nil {
		return x.TransmitBlankingState
	}
	return TransceiverTransmitBlankingState_TB_UNKNOWN
}

func (x *TransceiverGetStatusResponse) GetModemAsicTemp() float32 {
	if x != nil {
		return x.ModemAsicTemp
	}
	return 0
}

func (x *TransceiverGetStatusResponse) GetTxIfTemp() float32 {
	if x != nil {
		return x.TxIfTemp
	}
	return 0
}

type TransceiverFaults struct {
	state                  protoimpl.MessageState `protogen:"open.v1"`
	OverTempModemAsicFault bool                   `protobuf:"varint,1,opt,name=over_temp_modem_asic_fault,json=overTempModemAsicFault,proto3" json:"over_temp_modem_asic_fault,omitempty"`
	OverTempPcbaFault      bool                   `protobuf:"varint,2,opt,name=over_temp_pcba_fault,json=overTempPcbaFault,proto3" json:"over_temp_pcba_fault,omitempty"`
	DcVoltageFault         bool                   `protobuf:"varint,3,opt,name=dc_voltage_fault,json=dcVoltageFault,proto3" json:"dc_voltage_fault,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *TransceiverFaults) Reset() {
	*x = TransceiverFaults{}
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransceiverFaults) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransceiverFaults) ProtoMessage() {}

func (x *TransceiverFaults) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransceiverFaults.ProtoReflect.Descriptor instead.
func (*TransceiverFaults) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{4}
}

func (x *TransceiverFaults) GetOverTempModemAsicFault() bool {
	if x != nil {
		return x.OverTempModemAsicFault
	}
	return false
}

func (x *TransceiverFaults) GetOverTempPcbaFault() bool {
	if x != nil {
		return x.OverTempPcbaFault
	}
	return false
}

func (x *TransceiverFaults) GetDcVoltageFault() bool {
	if x != nil {
		return x.DcVoltageFault
	}
	return false
}

type TransceiverGetTelemetryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransceiverGetTelemetryRequest) Reset() {
	*x = TransceiverGetTelemetryRequest{}
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransceiverGetTelemetryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransceiverGetTelemetryRequest) ProtoMessage() {}

func (x *TransceiverGetTelemetryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransceiverGetTelemetryRequest.ProtoReflect.Descriptor instead.
func (*TransceiverGetTelemetryRequest) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{5}
}

type TransceiverGetTelemetryResponse struct {
	state                              protoimpl.MessageState `protogen:"open.v1"`
	AntennaPointingMode                uint32                 `protobuf:"varint,1001,opt,name=antenna_pointing_mode,json=antennaPointingMode,proto3" json:"antenna_pointing_mode,omitempty"`
	AntennaPitch                       float32                `protobuf:"fixed32,1002,opt,name=antenna_pitch,json=antennaPitch,proto3" json:"antenna_pitch,omitempty"`
	AntennaRoll                        float32                `protobuf:"fixed32,1003,opt,name=antenna_roll,json=antennaRoll,proto3" json:"antenna_roll,omitempty"`
	AntennaRxTheta                     float32                `protobuf:"fixed32,1004,opt,name=antenna_rx_theta,json=antennaRxTheta,proto3" json:"antenna_rx_theta,omitempty"`
	AntennaTrueHeading                 float32                `protobuf:"fixed32,1005,opt,name=antenna_true_heading,json=antennaTrueHeading,proto3" json:"antenna_true_heading,omitempty"`
	RxChannel                          uint32                 `protobuf:"varint,1006,opt,name=rx_channel,json=rxChannel,proto3" json:"rx_channel,omitempty"`
	CurrentCellId                      uint32                 `protobuf:"varint,1007,opt,name=current_cell_id,json=currentCellId,proto3" json:"current_cell_id,omitempty"`
	SecondsUntilSlotEnd                float32                `protobuf:"fixed32,1008,opt,name=seconds_until_slot_end,json=secondsUntilSlotEnd,proto3" json:"seconds_until_slot_end,omitempty"`
	WbRssiPeakMagDb                    float32                `protobuf:"fixed32,1009,opt,name=wb_rssi_peak_mag_db,json=wbRssiPeakMagDb,proto3" json:"wb_rssi_peak_mag_db,omitempty"`
	PopPingDropRate                    float32                `protobuf:"fixed32,1010,opt,name=pop_ping_drop_rate,json=popPingDropRate,proto3" json:"pop_ping_drop_rate,omitempty"`
	SnrDb                              float32                `protobuf:"fixed32,1011,opt,name=snr_db,json=snrDb,proto3" json:"snr_db,omitempty"`
	L1SnrAvgDb                         float32                `protobuf:"fixed32,1012,opt,name=l1_snr_avg_db,json=l1SnrAvgDb,proto3" json:"l1_snr_avg_db,omitempty"`
	L1SnrMinDb                         float32                `protobuf:"fixed32,1013,opt,name=l1_snr_min_db,json=l1SnrMinDb,proto3" json:"l1_snr_min_db,omitempty"`
	L1SnrMaxDb                         float32                `protobuf:"fixed32,1014,opt,name=l1_snr_max_db,json=l1SnrMaxDb,proto3" json:"l1_snr_max_db,omitempty"`
	LmacSatelliteId                    uint32                 `protobuf:"varint,1015,opt,name=lmac_satellite_id,json=lmacSatelliteId,proto3" json:"lmac_satellite_id,omitempty"`
	TargetSatelliteId                  uint32                 `protobuf:"varint,1016,opt,name=target_satellite_id,json=targetSatelliteId,proto3" json:"target_satellite_id,omitempty"`
	GrantMcs                           uint32                 `protobuf:"varint,1017,opt,name=grant_mcs,json=grantMcs,proto3" json:"grant_mcs,omitempty"`
	GrantSymbolsAvg                    float32                `protobuf:"fixed32,1018,opt,name=grant_symbols_avg,json=grantSymbolsAvg,proto3" json:"grant_symbols_avg,omitempty"`
	DedGrant                           uint32                 `protobuf:"varint,1019,opt,name=ded_grant,json=dedGrant,proto3" json:"ded_grant,omitempty"`
	MobilityProactiveSlotChange        uint32                 `protobuf:"varint,1020,opt,name=mobility_proactive_slot_change,json=mobilityProactiveSlotChange,proto3" json:"mobility_proactive_slot_change,omitempty"`
	MobilityReactiveSlotChange         uint32                 `protobuf:"varint,1021,opt,name=mobility_reactive_slot_change,json=mobilityReactiveSlotChange,proto3" json:"mobility_reactive_slot_change,omitempty"`
	RfpTotalSynFailed                  uint32                 `protobuf:"varint,1022,opt,name=rfp_total_syn_failed,json=rfpTotalSynFailed,proto3" json:"rfp_total_syn_failed,omitempty"`
	NumOutOfSeq                        uint32                 `protobuf:"varint,1023,opt,name=num_out_of_seq,json=numOutOfSeq,proto3" json:"num_out_of_seq,omitempty"`
	NumUlmapDrop                       uint32                 `protobuf:"varint,1024,opt,name=num_ulmap_drop,json=numUlmapDrop,proto3" json:"num_ulmap_drop,omitempty"`
	CurrentSecondsOfSchedule           float32                `protobuf:"fixed32,1025,opt,name=current_seconds_of_schedule,json=currentSecondsOfSchedule,proto3" json:"current_seconds_of_schedule,omitempty"`
	SendLabelSwitchToGroundFailedCalls uint32                 `protobuf:"varint,1026,opt,name=send_label_switch_to_ground_failed_calls,json=sendLabelSwitchToGroundFailedCalls,proto3" json:"send_label_switch_to_ground_failed_calls,omitempty"`
	EmaVelocityX                       float64                `protobuf:"fixed64,1027,opt,name=ema_velocity_x,json=emaVelocityX,proto3" json:"ema_velocity_x,omitempty"`
	EmaVelocityY                       float64                `protobuf:"fixed64,1028,opt,name=ema_velocity_y,json=emaVelocityY,proto3" json:"ema_velocity_y,omitempty"`
	EmaVelocityZ                       float64                `protobuf:"fixed64,1029,opt,name=ema_velocity_z,json=emaVelocityZ,proto3" json:"ema_velocity_z,omitempty"`
	CeRssiDb                           float32                `protobuf:"fixed32,1030,opt,name=ce_rssi_db,json=ceRssiDb,proto3" json:"ce_rssi_db,omitempty"`
	unknownFields                      protoimpl.UnknownFields
	sizeCache                          protoimpl.SizeCache
}

func (x *TransceiverGetTelemetryResponse) Reset() {
	*x = TransceiverGetTelemetryResponse{}
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransceiverGetTelemetryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransceiverGetTelemetryResponse) ProtoMessage() {}

func (x *TransceiverGetTelemetryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_spacex_api_device_transceiver_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransceiverGetTelemetryResponse.ProtoReflect.Descriptor instead.
func (*TransceiverGetTelemetryResponse) Descriptor() ([]byte, []int) {
	return file_spacex_api_device_transceiver_proto_rawDescGZIP(), []int{6}
}

func (x *TransceiverGetTelemetryResponse) GetAntennaPointingMode() uint32 {
	if x != nil {
		return x.AntennaPointingMode
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetAntennaPitch() float32 {
	if x != nil {
		return x.AntennaPitch
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetAntennaRoll() float32 {
	if x != nil {
		return x.AntennaRoll
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetAntennaRxTheta() float32 {
	if x != nil {
		return x.AntennaRxTheta
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetAntennaTrueHeading() float32 {
	if x != nil {
		return x.AntennaTrueHeading
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetRxChannel() uint32 {
	if x != nil {
		return x.RxChannel
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetCurrentCellId() uint32 {
	if x != nil {
		return x.CurrentCellId
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetSecondsUntilSlotEnd() float32 {
	if x != nil {
		return x.SecondsUntilSlotEnd
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetWbRssiPeakMagDb() float32 {
	if x != nil {
		return x.WbRssiPeakMagDb
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetPopPingDropRate() float32 {
	if x != nil {
		return x.PopPingDropRate
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetSnrDb() float32 {
	if x != nil {
		return x.SnrDb
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetL1SnrAvgDb() float32 {
	if x != nil {
		return x.L1SnrAvgDb
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetL1SnrMinDb() float32 {
	if x != nil {
		return x.L1SnrMinDb
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetL1SnrMaxDb() float32 {
	if x != nil {
		return x.L1SnrMaxDb
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetLmacSatelliteId() uint32 {
	if x != nil {
		return x.LmacSatelliteId
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetTargetSatelliteId() uint32 {
	if x != nil {
		return x.TargetSatelliteId
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetGrantMcs() uint32 {
	if x != nil {
		return x.GrantMcs
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetGrantSymbolsAvg() float32 {
	if x != nil {
		return x.GrantSymbolsAvg
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetDedGrant() uint32 {
	if x != nil {
		return x.DedGrant
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetMobilityProactiveSlotChange() uint32 {
	if x != nil {
		return x.MobilityProactiveSlotChange
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetMobilityReactiveSlotChange() uint32 {
	if x != nil {
		return x.MobilityReactiveSlotChange
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetRfpTotalSynFailed() uint32 {
	if x != nil {
		return x.RfpTotalSynFailed
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetNumOutOfSeq() uint32 {
	if x != nil {
		return x.NumOutOfSeq
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetNumUlmapDrop() uint32 {
	if x != nil {
		return x.NumUlmapDrop
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetCurrentSecondsOfSchedule() float32 {
	if x != nil {
		return x.CurrentSecondsOfSchedule
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetSendLabelSwitchToGroundFailedCalls() uint32 {
	if x != nil {
		return x.SendLabelSwitchToGroundFailedCalls
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetEmaVelocityX() float64 {
	if x != nil {
		return x.EmaVelocityX
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetEmaVelocityY() float64 {
	if x != nil {
		return x.EmaVelocityY
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetEmaVelocityZ() float64 {
	if x != nil {
		return x.EmaVelocityZ
	}
	return 0
}

func (x *TransceiverGetTelemetryResponse) GetCeRssiDb() float32 {
	if x != nil {
		return x.CeRssiDb
	}
	return 0
}

var File_spacex_api_device_transceiver_proto protoreflect.FileDescriptor

const file_spacex_api_device_transceiver_proto_rawDesc = "" +
	"\n" +
	"#spacex_api/device/transceiver.proto\x12\x11SpaceX.API.Device\x1a\x1cspacex_api/device/dish.proto\"P\n" +
	" TransceiverIFLoopbackTestRequest\x12,\n" +
	"\x12enable_if_loopback\x18\x01 \x01(\bR\x10enableIfLoopback\"\xc4\x01\n" +
	"!TransceiverIFLoopbackTestResponse\x12*\n" +
	"\x11ber_loopback_test\x18\x01 \x01(\x02R\x0fberLoopbackTest\x12*\n" +
	"\x11snr_loopback_test\x18\x02 \x01(\x02R\x0fsnrLoopbackTest\x12,\n" +
	"\x12rssi_loopback_test\x18\x03 \x01(\x02R\x10rssiLoopbackTest\x12\x19\n" +
	"\bpll_lock\x18\x04 \x01(\bR\apllLock\"\x1d\n" +
	"\x1bTransceiverGetStatusRequest\"\xea\x04\n" +
	"\x1cTransceiverGetStatusResponse\x12I\n" +
	"\tmod_state\x18\x01 \x01(\x0e2,.SpaceX.API.Device.TransceiverModulatorStateR\bmodState\x12M\n" +
	"\vdemod_state\x18\x02 \x01(\x0e2,.SpaceX.API.Device.TransceiverModulatorStateR\n" +
	"demodState\x12B\n" +
	"\btx_state\x18\x03 \x01(\x0e2'.SpaceX.API.Device.TransceiverTxRxStateR\atxState\x12B\n" +
	"\brx_state\x18\x04 \x01(\x0e2'.SpaceX.API.Device.TransceiverTxRxStateR\arxState\x123\n" +
	"\x05state\x18\xee\a \x01(\x0e2\x1c.SpaceX.API.Device.DishStateR\x05state\x12=\n" +
	"\x06faults\x18\xef\a \x01(\v2$.SpaceX.API.Device.TransceiverFaultsR\x06faults\x12l\n" +
	"\x17transmit_blanking_state\x18\xf0\a \x01(\x0e23.SpaceX.API.Device.TransceiverTransmitBlankingStateR\x15transmitBlankingState\x12'\n" +
	"\x0fmodem_asic_temp\x18\xf1\a \x01(\x02R\rmodemAsicTemp\x12\x1d\n" +
	"\n" +
	"tx_if_temp\x18\xf2\a \x01(\x02R\btxIfTemp\"\xaa\x01\n" +
	"\x11TransceiverFaults\x12:\n" +
	"\x1aover_temp_modem_asic_fault\x18\x01 \x01(\bR\x16overTempModemAsicFault\x12/\n" +
	"\x14over_temp_pcba_fault\x18\x02 \x01(\bR\x11overTempPcbaFault\x12(\n" +
	"\x10dc_voltage_fault\x18\x03 \x01(\bR\x0edcVoltageFault\" \n" +
	"\x1eTransceiverGetTelemetryRequest\"\xd9\n" +
	"\n" +
	"\x1fTransceiverGetTelemetryResponse\x123\n" +
	"\x15antenna_pointing_mode\x18\xe9\a \x01(\rR\x13antennaPointingMode\x12$\n" +
	"\rantenna_pitch\x18\xea\a \x01(\x02R\fantennaPitch\x12\"\n" +
	"\fantenna_roll\x18\xeb\a \x01(\x02R\vantennaRoll\x12)\n" +
	"\x10antenna_rx_theta\x18\xec\a \x01(\x02R\x0eantennaRxTheta\x121\n" +
	"\x14antenna_true_heading\x18\xed\a \x01(\x02R\x12antennaTrueHeading\x12\x1e\n" +
	"\n" +
	"rx_channel\x18\xee\a \x01(\rR\trxChannel\x12'\n" +
	"\x0fcurrent_cell_id\x18\xef\a \x01(\rR\rcurrentCellId\x124\n" +
	"\x16seconds_until_slot_end\x18\xf0\a \x01(\x02R\x13secondsUntilSlotEnd\x12-\n" +
	"\x13wb_rssi_peak_mag_db\x18\xf1\a \x01(\x02R\x0fwbRssiPeakMagDb\x12,\n" +
	"\x12pop_ping_drop_rate\x18\xf2\a \x01(\x02R\x0fpopPingDropRate\x12\x16\n" +
	"\x06snr_db\x18\xf3\a \x01(\x02R\x05snrDb\x12\"\n" +
	"\rl1_snr_avg_db\x18\xf4\a \x01(\x02R\n" +
	"l1SnrAvgDb\x12\"\n" +
	"\rl1_snr_min_db\x18\xf5\a \x01(\x02R\n" +
	"l1SnrMinDb\x12\"\n" +
	"\rl1_snr_max_db\x18\xf6\a \x01(\x02R\n" +
	"l1SnrMaxDb\x12+\n" +
	"\x11lmac_satellite_id\x18\xf7\a \x01(\rR\x0flmacSatelliteId\x12/\n" +
	"\x13target_satellite_id\x18\xf8\a \x01(\rR\x11targetSatelliteId\x12\x1c\n" +
	"\tgrant_mcs\x18\xf9\a \x01(\rR\bgrantMcs\x12+\n" +
	"\x11grant_symbols_avg\x18\xfa\a \x01(\x02R\x0fgrantSymbolsAvg\x12\x1c\n" +
	"\tded_grant\x18\xfb\a \x01(\rR\bdedGrant\x12D\n" +
	"\x1emobility_proactive_slot_change\x18\xfc\a \x01(\rR\x1bmobilityProactiveSlotChange\x12B\n" +
	"\x1dmobility_reactive_slot_change\x18\xfd\a \x01(\rR\x1amobilityReactiveSlotChange\x120\n" +
	"\x14rfp_total_syn_failed\x18\xfe\a \x01(\rR\x11rfpTotalSynFailed\x12$\n" +
	"\x0enum_out_of_seq\x18\xff\a \x01(\rR\vnumOutOfSeq\x12%\n" +
	"\x0enum_ulmap_drop\x18\x80\b \x01(\rR\fnumUlmapDrop\x12>\n" +
	"\x1bcurrent_seconds_of_schedule\x18\x81\b \x01(\x02R\x18currentSecondsOfSchedule\x12U\n" +
	"(send_label_switch_to_ground_failed_calls\x18\x82\b \x01(\rR\"sendLabelSwitchToGroundFailedCalls\x12%\n" +
	"\x0eema_velocity_x\x18\x83\b \x01(\x01R\femaVelocityX\x12%\n" +
	"\x0eema_velocity_y\x18\x84\b \x01(\x01R\femaVelocityY\x12%\n" +
	"\x0eema_velocity_z\x18\x85\b \x01(\x01R\femaVelocityZ\x12\x1d\n" +
	"\n" +
	"ce_rssi_db\x18\x86\b \x01(\x02R\bceRssiDb*^\n" +
	"\x19TransceiverModulatorState\x12\x14\n" +
	"\x10MODSTATE_UNKNOWN\x10\x00\x12\x14\n" +
	"\x10MODSTATE_ENABLED\x10\x01\x12\x15\n" +
	"\x11MODSTATE_DISABLED\x10\x02*M\n" +
	"\x14TransceiverTxRxState\x12\x10\n" +
	"\fTXRX_UNKNOWN\x10\x00\x12\x10\n" +
	"\fTXRX_ENABLED\x10\x01\x12\x11\n" +
	"\rTXRX_DISABLED\x10\x02*S\n" +
	" TransceiverTransmitBlankingState\x12\x0e\n" +
	"\n" +
	"TB_UNKNOWN\x10\x00\x12\x0e\n" +
	"\n" +
	"TB_ENABLED\x10\x01\x12\x0f\n" +
	"\vTB_DISABLED\x10\x02B\x17Z\x15spacex.com/api/deviceb\x06proto3"

var (
	file_spacex_api_device_transceiver_proto_rawDescOnce sync.Once
	file_spacex_api_device_transceiver_proto_rawDescData []byte
)

func file_spacex_api_device_transceiver_proto_rawDescGZIP() []byte {
	file_spacex_api_device_transceiver_proto_rawDescOnce.Do(func() {
		file_spacex_api_device_transceiver_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_spacex_api_device_transceiver_proto_rawDesc), len(file_spacex_api_device_transceiver_proto_rawDesc)))
	})
	return file_spacex_api_device_transceiver_proto_rawDescData
}

var file_spacex_api_device_transceiver_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_spacex_api_device_transceiver_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_spacex_api_device_transceiver_proto_goTypes = []any{
	(TransceiverModulatorState)(0),            // 0: SpaceX.API.Device.TransceiverModulatorState
	(TransceiverTxRxState)(0),                 // 1: SpaceX.API.Device.TransceiverTxRxState
	(TransceiverTransmitBlankingState)(0),     // 2: SpaceX.API.Device.TransceiverTransmitBlankingState
	(*TransceiverIFLoopbackTestRequest)(nil),  // 3: SpaceX.API.Device.TransceiverIFLoopbackTestRequest
	(*TransceiverIFLoopbackTestResponse)(nil), // 4: SpaceX.API.Device.TransceiverIFLoopbackTestResponse
	(*TransceiverGetStatusRequest)(nil),       // 5: SpaceX.API.Device.TransceiverGetStatusRequest
	(*TransceiverGetStatusResponse)(nil),      // 6: SpaceX.API.Device.TransceiverGetStatusResponse
	(*TransceiverFaults)(nil),                 // 7: SpaceX.API.Device.TransceiverFaults
	(*TransceiverGetTelemetryRequest)(nil),    // 8: SpaceX.API.Device.TransceiverGetTelemetryRequest
	(*TransceiverGetTelemetryResponse)(nil),   // 9: SpaceX.API.Device.TransceiverGetTelemetryResponse
	(DishState)(0),                            // 10: SpaceX.API.Device.DishState
}
var file_spacex_api_device_transceiver_proto_depIdxs = []int32{
	0,  // 0: SpaceX.API.Device.TransceiverGetStatusResponse.mod_state:type_name -> SpaceX.API.Device.TransceiverModulatorState
	0,  // 1: SpaceX.API.Device.TransceiverGetStatusResponse.demod_state:type_name -> SpaceX.API.Device.TransceiverModulatorState
	1,  // 2: SpaceX.API.Device.TransceiverGetStatusResponse.tx_state:type_name -> SpaceX.API.Device.TransceiverTxRxState
	1,  // 3: SpaceX.API.Device.TransceiverGetStatusResponse.rx_state:type_name -> SpaceX.API.Device.TransceiverTxRxState
	10, // 4: SpaceX.API.Device.TransceiverGetStatusResponse.state:type_name -> SpaceX.API.Device.DishState
	7,  // 5: SpaceX.API.Device.TransceiverGetStatusResponse.faults:type_name -> SpaceX.API.Device.TransceiverFaults
	2,  // 6: SpaceX.API.Device.TransceiverGetStatusResponse.transmit_blanking_state:type_name -> SpaceX.API.Device.TransceiverTransmitBlankingState
	7,  // [7:7] is the sub-list for method output_type
	7,  // [7:7] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_spacex_api_device_transceiver_proto_init() }
func file_spacex_api_device_transceiver_proto_init() {
	if File_spacex_api_device_transceiver_proto != nil {
		return
	}
	file_spacex_api_device_dish_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_spacex_api_device_transceiver_proto_rawDesc), len(file_spacex_api_device_transceiver_proto_rawDesc)),
			NumEnums:      3,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_spacex_api_device_transceiver_proto_goTypes,
		DependencyIndexes: file_spacex_api_device_transceiver_proto_depIdxs,
		EnumInfos:         file_spacex_api_device_transceiver_proto_enumTypes,
		MessageInfos:      file_spacex_api_device_transceiver_proto_msgTypes,
	}.Build()
	File_spacex_api_device_transceiver_proto = out.File
	file_spacex_api_device_transceiver_proto_goTypes = nil
	file_spacex_api_device_transceiver_proto_depIdxs = nil
}
