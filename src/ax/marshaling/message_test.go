package marshaling_test

import (
	"github.com/golang/protobuf/proto"
	"github.com/jmalloc/ax/src/ax/internal/messagetest"
	. "github.com/jmalloc/ax/src/ax/marshaling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalMessage", func() {
	message := &messagetest.Message{
		Value: "<value>",
	}

	It("marshals the message using protocol buffers", func() {
		_, data, err := MarshalMessage(message)
		Expect(err).ShouldNot(HaveOccurred())

		var m messagetest.Message
		err = proto.Unmarshal(data, &m)
		Expect(err).ShouldNot(HaveOccurred())

		Expect(proto.Equal(&m, message)).To(BeTrue())
	})

	It("includes the protocol information in the content-type", func() {
		ct, _, err := MarshalMessage(message)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ct).To(Equal("application/vnd.google.protobuf; proto=ax.internal.messagetest.Message"))
	})
})

var _ = Describe("UnmarshalMessage", func() {
	message := &messagetest.Message{
		Value: "<value>",
	}
	data, err := proto.Marshal(message)
	if err != nil {
		panic(err)
	}

	It("unmarshals the message using the protocol specified in the content-type", func() {
		m, err := UnmarshalMessage(
			"application/vnd.google.protobuf; proto=ax.internal.messagetest.Message",
			data,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(proto.Equal(m, message)).To(BeTrue())
	})

	It("returns an error if the content-type is invalid", func() {
		_, err := UnmarshalMessage("", data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content-type is not supported", func() {
		_, err := UnmarshalMessage("application/x-unknown", data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if an error occurs unmarshaling the protocol buffers message", func() {
		_, err := UnmarshalMessage(
			"application/vnd.google.protobuf; proto=ax.internal.messagetest.Unknown", // note unknown message type
			data,
		)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the buffer contains a protocol buffer message that is not an ax.Message", func() {
		data, err := proto.Marshal(&messagetest.NonAxMessage{})
		Expect(err).ShouldNot(HaveOccurred())

		_, err = UnmarshalMessage(
			"application/vnd.google.protobuf; proto=ax.internal.messagetest.NonAxMessage",
			data,
		)
		Expect(err).Should(HaveOccurred())
	})
})