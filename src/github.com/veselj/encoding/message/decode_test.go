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

func Test_SingleItem_IntStruct(t *testing.T) {
	g := NewGomegaWithT(t)
	type intS struct {
		IntField int `type:"int" len:"3" padding:"0"`
	}
	var ints intS
	err := message.Unmarshal([]byte(`001`), &ints)
	g.Expect(err).To(BeNil())
	g.Expect(ints.IntField).To(Equal(1))
}

func Test_IntAndStringStruct(t *testing.T) {
	g := NewGomegaWithT(t)
	type Message struct {
		IntField int    `len:"2" padding:"0"`
		StrField string `len:"5" padding:" "`
	}
	var msg Message
	err := message.Unmarshal([]byte(`03hello`), &msg)
	g.Expect(err).To(BeNil())
	g.Expect(msg.IntField).To(Equal(3))
	g.Expect(msg.StrField).To(Equal("hello"))
}

func Test_SeparatedStruct(t *testing.T) {
	g := NewGomegaWithT(t)
	type Message struct {
		IntField int    `sep:"\x1c" padding:"0"`
		StrField string `sep:"\x1c" padding:" "`
	}
	var msg Message
	err := message.Unmarshal([]byte(`045greetings`), &msg)
	g.Expect(err).To(BeNil())
	g.Expect(msg.IntField).To(Equal(45))
	g.Expect(msg.StrField).To(Equal("greetings"))
}

// func Test_StructSlice(t *testing.T) {
// 	g := NewGomegaWithT(t)
// 	type Message struct {
// 		StrFields []string `sep:"\x1c" padding:" "`
// 	}
// 	var msg Message
// 	err := message.Unmarshal([]byte(`123hi`), &msg)
// 	g.Expect(err).To(BeNil())
// 	g.Expect(msg.StrFields[0]).To(Equal("123"))
// 	g.Expect(msg.StrFields[1]).To(Equal("hi"))
// }
