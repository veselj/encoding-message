package message_test

import (
	"testing"

	message "github.com/veselj/encoding-message"

	. "github.com/onsi/gomega"
)

func Test_SingleItemStruct_Int(t *testing.T) {
	g := NewGomegaWithT(t)
	type intS struct {
		IntField int `type:"int" len:"3" padding:"0"`
	}
	var ints intS
	err := message.Unmarshal([]byte(`001`), &ints)
	g.Expect(err).To(BeNil())
	g.Expect(ints.IntField).To(Equal(1))
}

func Test_SingleItemStruct_String(t *testing.T) {
	g := NewGomegaWithT(t)
	type intS struct {
		IntField int `type:"int" len:"3" padding:"0"`
	}
	var ints intS
	err := message.Unmarshal([]byte(`001`), &ints)
	g.Expect(err).To(BeNil())
	g.Expect(ints.IntField).To(Equal(1))
}
