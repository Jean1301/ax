package ax_test

import (
	"reflect"

	. "github.com/jmalloc/ax/src/ax"
	"github.com/jmalloc/ax/src/ax/internal/messagetest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MessageType", func() {
	message := TypeOf(&messagetest.Message{})

	Describe("TypeOf", func() {
		It("returns a message type with the correct name", func() {
			Expect(message.Name).To(Equal("ax.internal.messagetest.Message"))
		})

		It("returns a message type with the correct struct type", func() {
			Expect(message.StructType).To(Equal(reflect.TypeOf(messagetest.Message{})))
		})
	})

	Describe("TypeByName", func() {
		It("returns a message type with the correct name", func() {
			mt, ok := TypeByName("ax.internal.messagetest.Message")
			Expect(ok).To(BeTrue())
			Expect(mt.Name).To(Equal("ax.internal.messagetest.Message"))
		})

		It("returns a message type with the correct struct type", func() {
			mt, ok := TypeByName("ax.internal.messagetest.Message")
			Expect(ok).To(BeTrue())
			Expect(mt.StructType).To(Equal(reflect.TypeOf(messagetest.Message{})))
		})

		It("returns false if the message name is not registered", func() {
			_, ok := TypeByName("ax.internal.messagetest.Unknown")
			Expect(ok).To(BeFalse())
		})
	})

	Context("when the message is generic", func() {
		Describe("IsCommand", func() {
			It("returns false", func() {
				Expect(message.IsCommand()).To(BeFalse())
			})
		})

		Describe("IsEvent", func() {
			It("returns false", func() {
				Expect(message.IsEvent()).To(BeFalse())
			})
		})
	})

	Context("when the message is a command", func() {
		command := TypeOf(&messagetest.Command{})

		Describe("IsCommand", func() {
			It("returns true", func() {
				Expect(command.IsCommand()).To(BeTrue())
			})
		})

		Describe("IsEvent", func() {
			It("returns false", func() {
				Expect(command.IsEvent()).To(BeFalse())
			})
		})
	})

	Context("when the message is an event", func() {
		event := TypeOf(&messagetest.Event{})

		Describe("IsCommand", func() {
			It("returns false", func() {
				Expect(event.IsCommand()).To(BeFalse())
			})
		})

		Describe("IsEvent", func() {
			It("returns true", func() {
				Expect(event.IsEvent()).To(BeTrue())
			})
		})
	})

	Describe("ToSet", func() {
		It("returns a set containing only this message type", func() {
			Expect(
				message.ToSet(),
			).To(Equal(
				TypesOf(&messagetest.Message{}),
			))
		})
	})

	Describe("New", func() {
		It("returns a pointer to a new instance of the message struct", func() {
			Expect(message.New()).To(Equal(&messagetest.Message{}))
		})
	})

	Describe("PackageName", func() {
		It("returns the protocol buffers package name", func() {
			Expect(message.PackageName()).To(Equal("ax.internal.messagetest"))
		})

		It("returns an empty string if the message is not in a package", func() {
			mt := TypeOf(&messagetest.NoPackage{})
			Expect(mt.PackageName()).To(Equal(""))
		})
	})
})
