// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.17.3
// source: node_event.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type NodeEventResp_StatusCode int32

const (
	NodeEventResp_SUCCESS NodeEventResp_StatusCode = 0
	NodeEventResp_FAILURE NodeEventResp_StatusCode = -1
)

// Enum value maps for NodeEventResp_StatusCode.
var (
	NodeEventResp_StatusCode_name = map[int32]string{
		0:  "SUCCESS",
		-1: "FAILURE",
	}
	NodeEventResp_StatusCode_value = map[string]int32{
		"SUCCESS": 0,
		"FAILURE": -1,
	}
)

func (x NodeEventResp_StatusCode) Enum() *NodeEventResp_StatusCode {
	p := new(NodeEventResp_StatusCode)
	*p = x
	return p
}

func (x NodeEventResp_StatusCode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NodeEventResp_StatusCode) Descriptor() protoreflect.EnumDescriptor {
	return file_node_event_proto_enumTypes[0].Descriptor()
}

func (NodeEventResp_StatusCode) Type() protoreflect.EnumType {
	return &file_node_event_proto_enumTypes[0]
}

func (x NodeEventResp_StatusCode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NodeEventResp_StatusCode.Descriptor instead.
func (NodeEventResp_StatusCode) EnumDescriptor() ([]byte, []int) {
	return file_node_event_proto_rawDescGZIP(), []int{2, 0}
}

// The request message
type NodeEventReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventType string            `protobuf:"bytes,1,opt,name=eventType,proto3" json:"eventType,omitempty"` // nodeInitBeginEvent,nodeInitEndEvent,complete
	Data      *NodeEventReqData `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *NodeEventReq) Reset() {
	*x = NodeEventReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_event_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeEventReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeEventReq) ProtoMessage() {}

func (x *NodeEventReq) ProtoReflect() protoreflect.Message {
	mi := &file_node_event_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeEventReq.ProtoReflect.Descriptor instead.
func (*NodeEventReq) Descriptor() ([]byte, []int) {
	return file_node_event_proto_rawDescGZIP(), []int{0}
}

func (x *NodeEventReq) GetEventType() string {
	if x != nil {
		return x.EventType
	}
	return ""
}

func (x *NodeEventReq) GetData() *NodeEventReqData {
	if x != nil {
		return x.Data
	}
	return nil
}

type NodeEventReqData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProcessID         string   `protobuf:"bytes,1,opt,name=processID,proto3" json:"processID,omitempty"`
	ProcessInstanceID string   `protobuf:"bytes,2,opt,name=processInstanceID,proto3" json:"processInstanceID,omitempty"`
	ExecutionID       string   `protobuf:"bytes,3,opt,name=executionID,proto3" json:"executionID,omitempty"`
	NodeDefKey        string   `protobuf:"bytes,4,opt,name=nodeDefKey,proto3" json:"nodeDefKey,omitempty"`
	NodeInstanceID    string   `protobuf:"bytes,5,opt,name=nodeInstanceID,proto3" json:"nodeInstanceID,omitempty"`
	TaskID            []string `protobuf:"bytes,6,rep,name=taskID,proto3" json:"taskID,omitempty"` // 节点初始化结束事件才返回：或签一个task，会签多个taskID
}

func (x *NodeEventReqData) Reset() {
	*x = NodeEventReqData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_event_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeEventReqData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeEventReqData) ProtoMessage() {}

func (x *NodeEventReqData) ProtoReflect() protoreflect.Message {
	mi := &file_node_event_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeEventReqData.ProtoReflect.Descriptor instead.
func (*NodeEventReqData) Descriptor() ([]byte, []int) {
	return file_node_event_proto_rawDescGZIP(), []int{1}
}

func (x *NodeEventReqData) GetProcessID() string {
	if x != nil {
		return x.ProcessID
	}
	return ""
}

func (x *NodeEventReqData) GetProcessInstanceID() string {
	if x != nil {
		return x.ProcessInstanceID
	}
	return ""
}

func (x *NodeEventReqData) GetExecutionID() string {
	if x != nil {
		return x.ExecutionID
	}
	return ""
}

func (x *NodeEventReqData) GetNodeDefKey() string {
	if x != nil {
		return x.NodeDefKey
	}
	return ""
}

func (x *NodeEventReqData) GetNodeInstanceID() string {
	if x != nil {
		return x.NodeInstanceID
	}
	return ""
}

func (x *NodeEventReqData) GetTaskID() []string {
	if x != nil {
		return x.TaskID
	}
	return nil
}

// The response message
type NodeEventResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code NodeEventResp_StatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=proto.NodeEventResp_StatusCode" json:"code,omitempty"`
	Msg  string                   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	Data *NodeEventRespData       `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *NodeEventResp) Reset() {
	*x = NodeEventResp{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_event_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeEventResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeEventResp) ProtoMessage() {}

func (x *NodeEventResp) ProtoReflect() protoreflect.Message {
	mi := &file_node_event_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeEventResp.ProtoReflect.Descriptor instead.
func (*NodeEventResp) Descriptor() ([]byte, []int) {
	return file_node_event_proto_rawDescGZIP(), []int{2}
}

func (x *NodeEventResp) GetCode() NodeEventResp_StatusCode {
	if x != nil {
		return x.Code
	}
	return NodeEventResp_SUCCESS
}

func (x *NodeEventResp) GetMsg() string {
	if x != nil {
		return x.Msg
	}
	return ""
}

func (x *NodeEventResp) GetData() *NodeEventRespData {
	if x != nil {
		return x.Data
	}
	return nil
}

type NodeEventRespData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Assignee    []string `protobuf:"bytes,1,rep,name=assignee,proto3" json:"assignee,omitempty"`       // 节点初始化开始事件才关注：节点处理人
	ExecuteType string   `protobuf:"bytes,2,opt,name=executeType,proto3" json:"executeType,omitempty"` // 节点初始化结束事件才关注，执行类型：pauseExecution，endExecution，endProcess
	Comments    string   `protobuf:"bytes,3,opt,name=comments,proto3" json:"comments,omitempty"`
}

func (x *NodeEventRespData) Reset() {
	*x = NodeEventRespData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_node_event_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeEventRespData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeEventRespData) ProtoMessage() {}

func (x *NodeEventRespData) ProtoReflect() protoreflect.Message {
	mi := &file_node_event_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeEventRespData.ProtoReflect.Descriptor instead.
func (*NodeEventRespData) Descriptor() ([]byte, []int) {
	return file_node_event_proto_rawDescGZIP(), []int{3}
}

func (x *NodeEventRespData) GetAssignee() []string {
	if x != nil {
		return x.Assignee
	}
	return nil
}

func (x *NodeEventRespData) GetExecuteType() string {
	if x != nil {
		return x.ExecuteType
	}
	return ""
}

func (x *NodeEventRespData) GetComments() string {
	if x != nil {
		return x.Comments
	}
	return ""
}

var File_node_event_proto protoreflect.FileDescriptor

var file_node_event_proto_rawDesc = []byte{
	0x0a, 0x10, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x59, 0x0a, 0x0c, 0x4e, 0x6f, 0x64,
	0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x12, 0x1c, 0x0a, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x2b, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x22, 0xe0, 0x01, 0x0a, 0x10, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x71, 0x44, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x0a, 0x09, 0x70, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72,
	0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x44, 0x12, 0x2c, 0x0a, 0x11, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x11, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x49, 0x6e, 0x73, 0x74, 0x61,
	0x6e, 0x63, 0x65, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x44, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x6e, 0x6f, 0x64, 0x65, 0x44,
	0x65, 0x66, 0x4b, 0x65, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6e, 0x6f, 0x64,
	0x65, 0x44, 0x65, 0x66, 0x4b, 0x65, 0x79, 0x12, 0x26, 0x0a, 0x0e, 0x6e, 0x6f, 0x64, 0x65, 0x49,
	0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x44, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x63, 0x65, 0x49, 0x44, 0x12,
	0x16, 0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x22, 0xb5, 0x01, 0x0a, 0x0d, 0x4e, 0x6f, 0x64, 0x65,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x33, 0x0a, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x2e, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x6d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6d, 0x73, 0x67,
	0x12, 0x2c, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x2f,
	0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07,
	0x53, 0x55, 0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x00, 0x12, 0x14, 0x0a, 0x07, 0x46, 0x41, 0x49,
	0x4c, 0x55, 0x52, 0x45, 0x10, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01, 0x22,
	0x6d, 0x0a, 0x11, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x1a, 0x0a, 0x08, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x61, 0x73, 0x73, 0x69, 0x67, 0x6e, 0x65, 0x65,
	0x12, 0x20, 0x0a, 0x0b, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x32, 0x41,
	0x0a, 0x09, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x34, 0x0a, 0x05, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x22,
	0x00, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2e, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_node_event_proto_rawDescOnce sync.Once
	file_node_event_proto_rawDescData = file_node_event_proto_rawDesc
)

func file_node_event_proto_rawDescGZIP() []byte {
	file_node_event_proto_rawDescOnce.Do(func() {
		file_node_event_proto_rawDescData = protoimpl.X.CompressGZIP(file_node_event_proto_rawDescData)
	})
	return file_node_event_proto_rawDescData
}

var file_node_event_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_node_event_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_node_event_proto_goTypes = []interface{}{
	(NodeEventResp_StatusCode)(0), // 0: proto.NodeEventResp.StatusCode
	(*NodeEventReq)(nil),          // 1: proto.NodeEventReq
	(*NodeEventReqData)(nil),      // 2: proto.NodeEventReqData
	(*NodeEventResp)(nil),         // 3: proto.NodeEventResp
	(*NodeEventRespData)(nil),     // 4: proto.NodeEventRespData
}
var file_node_event_proto_depIdxs = []int32{
	2, // 0: proto.NodeEventReq.data:type_name -> proto.NodeEventReqData
	0, // 1: proto.NodeEventResp.code:type_name -> proto.NodeEventResp.StatusCode
	4, // 2: proto.NodeEventResp.data:type_name -> proto.NodeEventRespData
	1, // 3: proto.NodeEvent.Event:input_type -> proto.NodeEventReq
	3, // 4: proto.NodeEvent.Event:output_type -> proto.NodeEventResp
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_node_event_proto_init() }
func file_node_event_proto_init() {
	if File_node_event_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_node_event_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeEventReq); i {
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
		file_node_event_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeEventReqData); i {
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
		file_node_event_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeEventResp); i {
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
		file_node_event_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeEventRespData); i {
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
			RawDescriptor: file_node_event_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_node_event_proto_goTypes,
		DependencyIndexes: file_node_event_proto_depIdxs,
		EnumInfos:         file_node_event_proto_enumTypes,
		MessageInfos:      file_node_event_proto_msgTypes,
	}.Build()
	File_node_event_proto = out.File
	file_node_event_proto_rawDesc = nil
	file_node_event_proto_goTypes = nil
	file_node_event_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// NodeEventClient is the client API for NodeEvent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type NodeEventClient interface {
	Event(ctx context.Context, in *NodeEventReq, opts ...grpc.CallOption) (*NodeEventResp, error)
}

type nodeEventClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeEventClient(cc grpc.ClientConnInterface) NodeEventClient {
	return &nodeEventClient{cc}
}

func (c *nodeEventClient) Event(ctx context.Context, in *NodeEventReq, opts ...grpc.CallOption) (*NodeEventResp, error) {
	out := new(NodeEventResp)
	err := c.cc.Invoke(ctx, "/proto.NodeEvent/Event", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeEventServer is the server API for NodeEvent service.
type NodeEventServer interface {
	Event(context.Context, *NodeEventReq) (*NodeEventResp, error)
}

// UnimplementedNodeEventServer can be embedded to have forward compatible implementations.
type UnimplementedNodeEventServer struct {
}

func (*UnimplementedNodeEventServer) Event(context.Context, *NodeEventReq) (*NodeEventResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Event not implemented")
}

func RegisterNodeEventServer(s *grpc.Server, srv NodeEventServer) {
	s.RegisterService(&_NodeEvent_serviceDesc, srv)
}

func _NodeEvent_Event_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeEventReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeEventServer).Event(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeEvent/Event",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeEventServer).Event(ctx, req.(*NodeEventReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _NodeEvent_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.NodeEvent",
	HandlerType: (*NodeEventServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Event",
			Handler:    _NodeEvent_Event_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "node_event.proto",
}
