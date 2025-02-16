package core

type Instruction byte

const (
	InstrPushInt Instruction = iota + 10
	InstrAdd
	InstrPushByte
	InstrPack
)

// Actualy a stack is FILO we need a Queue

// FIFO
type Queue struct {
	data []any
	sp   int
}

func NewQueue(size ...int) *Queue {
	if len(size) > 0 {
		return &Queue{
			data: make([]any, size[0]),
			sp:   0,
		}
	}
	return &Queue{
		data: []any{},
		sp:   0,
	}
}

// Push data to the Queue and returns the stack lenght
func (s *Queue) Push(data any) {
	s.data = append(s.data, data)
	s.sp++
}

func (s *Queue) Pop() any {
	data := s.data[0]
	s.data = s.data[1:]
	s.sp--
	return data
}

func (s *Queue) Last() any {

	return s.data[s.LastIndex()]
}

func (s *Queue) LastIndex() int {
	return len(s.data) - 1
}

type VM struct {
	data  []byte
	ip    int // instruction pointer
	queue *Queue
	sp    int // stack pointer
}

func NewVM(data []byte, queueSize ...int) *VM {
	if len(queueSize) == 0 {

		return &VM{
			data:  data,
			ip:    0,
			queue: NewQueue(),
			sp:    -1,
		}
	} else {
		return &VM{
			data:  data,
			ip:    0,
			queue: NewQueue(),
		}
	}
}

func (vm *VM) Run() error {
	for {
		instr := vm.data[vm.ip]
		if err := vm.Exec(Instruction(instr)); err != nil {
			return err
		}
		vm.ip++
		if vm.ip > len(vm.data)-1 {
			break
		}

	}

	return nil
}

func (vm *VM) Exec(instr Instruction) error {
	switch instr {
	case InstrPushInt:
		vm.queue.Push(int(vm.data[vm.ip-1]))
	case InstrAdd:
		a := vm.queue.Pop().(int)
		b := vm.queue.Pop().(int)
		c := a + b
		vm.queue.Push(c)
	case InstrPushByte:
		vm.queue.Push(vm.data[vm.ip-1])
	case InstrPack:
		n := vm.queue.Pop().(int)
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = vm.queue.Pop().(byte)
		}
		vm.queue.Push(b)
	}
	return nil
}
