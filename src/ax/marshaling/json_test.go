package marshaling_test

import (
	"encoding/json"

	"github.com/golang/protobuf/proto"
	. "github.com/jmalloc/ax/src/ax/marshaling"
	"github.com/jmalloc/ax/src/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MarshalJSON", func() {

	message := &messagetest.NonAxMessage{
		Value: "<value>",
	}

	It("marshals the message using JSON", func() {

		_, data, err := MarshalJSON(message)
		Expect(err).ShouldNot(HaveOccurred())

		var m messagetest.NonAxMessage
		err = json.Unmarshal(data, &m)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(m).Should(Equal(*message))
	})

	It("includes the protocol information in the content-type", func() {
		ct, _, err := MarshalJSON(message)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(ct).To(Equal(
			"application/json; proto=ax.internal.messagetest.NonAxMessage",
		))
	})

	It("returns an error if the protocol name is not available", func() {
		var m proto.Message // no concrete value

		_, _, err := MarshalJSON(m)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("UnmarshalJSON", func() {

	message := &messagetest.NonAxMessage{
		Value: "<value>",
	}
	_, data, err := MarshalJSON(message)
	if err != nil {
		panic(err)
	}

	It("unmarshals the message using the protocol specified in the content-type", func() {

		m, err := UnmarshalJSON(
			"application/json; proto=ax.internal.messagetest.NonAxMessage",
			data,
		)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(proto.Equal(m, message)).To(BeTrue())
	})

	It("returns an error if the content-type is invalid", func() {
		_, err := UnmarshalJSON("", data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content-type is not specific to JSON encoding", func() {
		_, err := UnmarshalJSON("application/x-unknown", data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the content-type does not specify protocol as a content type parameter", func() {
		_, err := UnmarshalJSON("application/json", data)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if message type is unregistered", func() {
		_, err := UnmarshalJSON(
			"application/json; ax.internal.messagetest.NonExistingType",
			data,
		)
		Expect(err).Should(HaveOccurred())
	})
})
