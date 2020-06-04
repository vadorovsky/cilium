// Code generated by protoc-gen-go. DO NOT EDIT.
// source: envoy/config/trace/v3/xray.proto

package envoy_config_trace_v3

import (
	fmt "fmt"
	v3 "github.com/cilium/proxy/go/envoy/config/core/v3"
	_ "github.com/cncf/udpa/go/udpa/annotations"
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type XRayConfig struct {
	// The UDP endpoint of the X-Ray Daemon where the spans will be sent.
	// If this value is not set, the default value of 127.0.0.1:2000 will be used.
	DaemonEndpoint *v3.SocketAddress `protobuf:"bytes,1,opt,name=daemon_endpoint,json=daemonEndpoint,proto3" json:"daemon_endpoint,omitempty"`
	// The name of the X-Ray segment. By default this will be set to the cluster name.
	SegmentName string `protobuf:"bytes,2,opt,name=segment_name,json=segmentName,proto3" json:"segment_name,omitempty"`
	// The location of a local custom sampling rules JSON file.
	// For an example of the sampling rules see:
	// `X-Ray SDK documentation
	// <https://docs.aws.amazon.com/xray/latest/devguide/xray-sdk-go-configuration.html#xray-sdk-go-configuration-sampling>`_
	SamplingRuleManifest *v3.DataSource `protobuf:"bytes,3,opt,name=sampling_rule_manifest,json=samplingRuleManifest,proto3" json:"sampling_rule_manifest,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *XRayConfig) Reset()         { *m = XRayConfig{} }
func (m *XRayConfig) String() string { return proto.CompactTextString(m) }
func (*XRayConfig) ProtoMessage()    {}
func (*XRayConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_6d491b3510a2e630, []int{0}
}

func (m *XRayConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_XRayConfig.Unmarshal(m, b)
}
func (m *XRayConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_XRayConfig.Marshal(b, m, deterministic)
}
func (m *XRayConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_XRayConfig.Merge(m, src)
}
func (m *XRayConfig) XXX_Size() int {
	return xxx_messageInfo_XRayConfig.Size(m)
}
func (m *XRayConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_XRayConfig.DiscardUnknown(m)
}

var xxx_messageInfo_XRayConfig proto.InternalMessageInfo

func (m *XRayConfig) GetDaemonEndpoint() *v3.SocketAddress {
	if m != nil {
		return m.DaemonEndpoint
	}
	return nil
}

func (m *XRayConfig) GetSegmentName() string {
	if m != nil {
		return m.SegmentName
	}
	return ""
}

func (m *XRayConfig) GetSamplingRuleManifest() *v3.DataSource {
	if m != nil {
		return m.SamplingRuleManifest
	}
	return nil
}

func init() {
	proto.RegisterType((*XRayConfig)(nil), "envoy.config.trace.v3.XRayConfig")
}

func init() { proto.RegisterFile("envoy/config/trace/v3/xray.proto", fileDescriptor_6d491b3510a2e630) }

var fileDescriptor_6d491b3510a2e630 = []byte{
	// 327 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x91, 0xc1, 0x4a, 0xeb, 0x40,
	0x14, 0x86, 0x49, 0x2f, 0x5c, 0xe8, 0xf4, 0x72, 0x2f, 0x84, 0xab, 0x96, 0x2e, 0x34, 0x6d, 0x11,
	0xbb, 0x90, 0x19, 0x68, 0x76, 0xee, 0xac, 0xba, 0x53, 0x29, 0x29, 0x48, 0x77, 0xe1, 0x34, 0x39,
	0x8d, 0x83, 0xc9, 0x9c, 0x30, 0x33, 0x09, 0xcd, 0xda, 0x8d, 0xcf, 0xe0, 0xfb, 0xf8, 0x5e, 0xd2,
	0x24, 0x45, 0x0a, 0xd9, 0x85, 0xcc, 0xf7, 0xcd, 0x7f, 0xce, 0x3f, 0xcc, 0x43, 0x55, 0x52, 0x25,
	0x22, 0x52, 0x5b, 0x99, 0x08, 0xab, 0x21, 0x42, 0x51, 0xfa, 0x62, 0xa7, 0xa1, 0xe2, 0xb9, 0x26,
	0x4b, 0xee, 0x49, 0x4d, 0xf0, 0x86, 0xe0, 0x35, 0xc1, 0x4b, 0x7f, 0x34, 0x39, 0x12, 0x23, 0xd2,
	0xb5, 0x07, 0x71, 0xac, 0xd1, 0x98, 0x46, 0x1d, 0x5d, 0x74, 0x32, 0x1b, 0x30, 0xd8, 0x02, 0xe3,
	0x22, 0xce, 0x41, 0x80, 0x52, 0x64, 0xc1, 0x4a, 0x52, 0x46, 0x94, 0xa8, 0x8d, 0x24, 0x25, 0x55,
	0xd2, 0x22, 0x67, 0x25, 0xa4, 0x32, 0x06, 0x8b, 0xe2, 0xf0, 0xd1, 0x1c, 0x4c, 0xde, 0x7b, 0x8c,
	0xad, 0x03, 0xa8, 0xee, 0xea, 0xdb, 0xdd, 0x47, 0xf6, 0x2f, 0x06, 0xcc, 0x48, 0x85, 0xa8, 0xe2,
	0x9c, 0xa4, 0xb2, 0x43, 0xc7, 0x73, 0x66, 0x83, 0xf9, 0x94, 0x1f, 0x2d, 0xb0, 0x9f, 0x82, 0x97,
	0x3e, 0x5f, 0x51, 0xf4, 0x86, 0xf6, 0xb6, 0x99, 0x37, 0xf8, 0xdb, 0xb8, 0x0f, 0xad, 0xea, 0x8e,
	0xd9, 0x1f, 0x83, 0x49, 0x86, 0xca, 0x86, 0x0a, 0x32, 0x1c, 0xf6, 0x3c, 0x67, 0xd6, 0x0f, 0x06,
	0xed, 0xbf, 0x67, 0xc8, 0xd0, 0x7d, 0x61, 0xa7, 0x06, 0xb2, 0x3c, 0x95, 0x2a, 0x09, 0x75, 0x91,
	0x62, 0x98, 0x81, 0x92, 0x5b, 0x34, 0x76, 0xf8, 0xab, 0xce, 0xf5, 0xba, 0x73, 0xef, 0xc1, 0xc2,
	0x8a, 0x0a, 0x1d, 0x61, 0xf0, 0xff, 0xe0, 0x07, 0x45, 0x8a, 0x4f, 0xad, 0x7d, 0x73, 0xfd, 0xf9,
	0xf5, 0x71, 0x7e, 0xc5, 0x2e, 0xbb, 0x6a, 0x9f, 0x43, 0x9a, 0xbf, 0x02, 0xff, 0x59, 0x7b, 0x31,
	0x67, 0x53, 0x49, 0x4d, 0x52, 0xae, 0x69, 0x57, 0xf1, 0xce, 0xd7, 0x5a, 0xf4, 0xd7, 0x1a, 0xaa,
	0xe5, 0xbe, 0xb7, 0xa5, 0xb3, 0xf9, 0x5d, 0x17, 0xe8, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0x26,
	0x5c, 0x68, 0x0d, 0xfc, 0x01, 0x00, 0x00,
}
