// Code generated by protoc-gen-go. DO NOT EDIT.
// source: src/internal/messagetest/nopackage.proto

package messagetest

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type NoPackage struct {
}

func (m *NoPackage) Reset()                    { *m = NoPackage{} }
func (m *NoPackage) String() string            { return proto.CompactTextString(m) }
func (*NoPackage) ProtoMessage()               {}
func (*NoPackage) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func init() {
	proto.RegisterType((*NoPackage)(nil), "NoPackage")
}

func init() { proto.RegisterFile("src/internal/messagetest/nopackage.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 89 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xd2, 0x2e, 0x2e, 0x4a, 0xd6,
	0x4f, 0xac, 0xd0, 0xcf, 0xcc, 0x2b, 0x49, 0x2d, 0xca, 0x4b, 0xcc, 0xd1, 0xcf, 0x4d, 0x2d, 0x2e,
	0x4e, 0x4c, 0x4f, 0x2d, 0x49, 0x2d, 0x2e, 0xd1, 0xcf, 0xcb, 0x2f, 0x48, 0x4c, 0xce, 0x4e, 0x4c,
	0x4f, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x57, 0xe2, 0xe6, 0xe2, 0xf4, 0xcb, 0x0f, 0x80, 0x08,
	0x39, 0xf1, 0x46, 0x71, 0x23, 0xa9, 0x4d, 0x62, 0x03, 0x2b, 0x31, 0x06, 0x04, 0x00, 0x00, 0xff,
	0xff, 0xe5, 0xa2, 0x93, 0x7b, 0x51, 0x00, 0x00, 0x00,
}
