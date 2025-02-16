package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	s := NewQueue()

	s.Push(1)
	s.Push(2)
	s.Push(3)
	s.Push(4)
	s.Push(5)
	s.Push(6)

	value := s.Pop()

	assert.Equal(t, value, 1)
	value = s.Pop()

	assert.Equal(t, value, 2)

}

func TestVM(t *testing.T) {
	// 2 + 2
	// 2
	// push to stack
	// 2
	// push stack
	// add
	// 4
	// push stack
	data := []byte{0x2, 0x0a, 0x2, 0x0a, 0x0b}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	assert.Equal(t, 4, vm.queue.Last())
}
func TestVMString(t *testing.T) {
	// 2 + 2
	// 2
	// push to stack
	// 2
	// push stack
	// add
	// 4
	// push stack
	data := []byte{0x03, 0x0a, 0x61, 0x0c, 0x61, 0x0c, 0x62, 0x0c, 0x0d}
	vm := NewVM(data)
	assert.Nil(t, vm.Run())
	// assert.Equal(t, 4, vm.queue.Last())

	value := vm.queue.Pop().([]byte)

	fmt.Println(string(value))

	fmt.Println(vm.queue.data...)
}
