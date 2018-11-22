package message_test

import (
	"testing"

	"github.com/veselj/encoding/message"

	. "github.com/onsi/gomega"
)

func Test_ShouldFailWrongOutParam(t *testing.T) {
	g := NewGomegaWithT(t)
	// nil out param
	err := message.Unmarshal([]byte("0"), nil)
	g.Expect(err).NotTo(BeNil())
	g.Expect(err.Error()).To(Equal("input not a pointer or nil"))
	// out param not a pointer
	var i int
	err = message.Unmarshal([]byte("03"), i)
	g.Expect(err).NotTo(BeNil())
	g.Expect(err.Error()).To(Equal("input not a pointer or nil"))
	// out param not a struct
	err = message.Unmarshal([]byte("03"), &i)
	g.Expect(err).NotTo(BeNil())
	g.Expect(err.Error()).To(Equal("input must point to a struct"))
}

func Test_SingleItemStruct(t *testing.T) {
	g := NewGomegaWithT(t)
	type intS struct {
		IntField int `type:"int" len:"3" padding:"0"`
	}
	var ints intS
	err := message.Unmarshal([]byte(`001`), &ints)
	g.Expect(err).To(BeNil())
	g.Expect(ints.IntField).To(Equal(1))
}
