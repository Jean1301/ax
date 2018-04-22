package ax_test

//
// import (
// 	. "github.com/jmalloc/ax/src/ax"
// 	. "github.com/onsi/ginkgo"
// 	. "github.com/onsi/gomega"
// )
//
// var _ = Describe("AggregateBehavior", func() {
// 	var behavior *AggregateBehavior
// 	BeforeEach(func() {
// 		behavior = &AggregateBehavior{}
// 	})
//
// 	Describe("AggregateID", func() {
// 		It("returns the aggregate ID", func() {
// 			var id AggregateID
// 			id.GenerateUUID()
// 			behavior.SetAggregateID(id)
//
// 			Expect(behavior.AggregateID()).To(Equal(id))
// 		})
//
// 		It("panics if the aggregate ID has not been set", func() {
// 			Expect(func() {
// 				behavior.AggregateID()
// 			}).To(Panic())
// 		})
// 	})
//
// 	Describe("SetAggregateID", func() {
// 		It("panics if the aggregate ID has already been set", func() {
// 			var id AggregateID
// 			id.GenerateUUID()
// 			behavior.SetAggregateID(id)
//
// 			Expect(func() {
// 				behavior.SetAggregateID(id)
// 			}).To(Panic())
// 		})
// 	})
//
// 	Describe("Revision", func() {
// 		It("returns zero by default", func() {
// 			Expect(behavior.Revision()).To(BeNumerically("==", 0))
// 		})
//
// 		It("returns the revision", func() {
// 			behavior.SetRevision(123)
//
// 			Expect(behavior.Revision()).To(BeNumerically("==", 123))
// 		})
// 	})
// })
