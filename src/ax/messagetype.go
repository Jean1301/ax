package ax

import (
	"reflect"
	"strings"

	"github.com/golang/protobuf/proto"
)

// MessageType identifies a unique message by its Protocol Buffers name.
type MessageType struct {
	// Name is the fully-qualified Protocol Buffers name of the message type.
	Name string

	// StructType is the type of the struct used to represent the message.
	// Note that only a pointer-to-struct type will satisfy the Message
	// interface.
	StructType reflect.Type
}

// TypeOf returns the message type of m.
func TypeOf(m Message) MessageType {
	return MessageType{
		proto.MessageName(m),
		reflect.TypeOf(m).Elem(),
	}
}

// TypeByName return the message type by the fully qualified Protocol Buffer
// name.
func TypeByName(n string) (MessageType, bool) {
	rt := proto.MessageType(n)

	if rt == nil {
		return MessageType{}, false
	}

	return MessageType{
		n,
		rt.Elem(),
	}, true
}

// ToSet returns a MessageTypeSet containing mt as its only member.
func (mt MessageType) ToSet() MessageTypeSet {
	return MessageTypeSet{
		map[MessageType]struct{}{mt: struct{}{}},
	}
}

// New returns a new pointer to a zero-value message.
func (mt MessageType) New() Message {
	return reflect.New(mt.StructType).Interface().(Message)
}

// Package returns the Protocol Buffers package name for this message type.
func (mt MessageType) PackageName() string {
	i := strings.LastIndexByte(mt.Name, '.')
	if i == -1 {
		return ""
	}

	return mt.Name[:i]
}

// IsCommand returns true if the message type is a command.
func (mt MessageType) IsCommand() bool {
	return reflect.PtrTo(mt.StructType).Implements(commandType)
}

// IsEvent returns true if the message type is an event.
func (mt MessageType) IsEvent() bool {
	return reflect.PtrTo(mt.StructType).Implements(eventType)
}

// MessageTypeSet is a set of unique message types.
type MessageTypeSet struct {
	members map[MessageType]struct{}
}

// TypesOf returns the message types of the elements in m.
func TypesOf(m ...Message) MessageTypeSet {
	members := make(map[MessageType]struct{}, len(m))

	for _, v := range m {
		members[TypeOf(v)] = struct{}{}
	}

	return MessageTypeSet{members}
}

// Members returns the message types in the set.
func (s MessageTypeSet) Members() []MessageType {
	types := make([]MessageType, 0, len(s.members))

	for mt := range s.members {
		types = append(types, mt)
	}

	return types
}

// Len returns the number of types in the set.
func (s MessageTypeSet) Len() int {
	return len(s.members)
}

// Has returns true if mt is a member of the set.
func (s MessageTypeSet) Has(mt MessageType) bool {
	_, ok := s.members[mt]
	return ok
}

// Add returns a new set containing the members of this set and mt.
func (s MessageTypeSet) Add(mt MessageType) MessageTypeSet {
	if s.Has(mt) {
		return s
	}

	members := make(map[MessageType]struct{}, len(s.members)+1)

	members[mt] = struct{}{}

	for mt := range s.members {
		members[mt] = struct{}{}
	}

	return MessageTypeSet{members}
}

// Union returns the set union of s and o.
func (s MessageTypeSet) Union(o MessageTypeSet) MessageTypeSet {
	if o.Len() == 0 {
		return s
	} else if s.Len() == 0 {
		return o
	}

	members := make(map[MessageType]struct{}, len(s.members)+len(o.members))

	for mt := range s.members {
		members[mt] = struct{}{}
	}

	for mt := range o.members {
		members[mt] = struct{}{}
	}

	return MessageTypeSet{members}
}
