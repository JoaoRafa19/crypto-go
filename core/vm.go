package core

type Instruction byte

const (
	InstrPush Instruction = iota + 10
	InstrAdd
)

type Stack struct {
	arr []byte
}

// Push data to the Stack and returns the stack lenght
func (s *Stack) Push(data byte) int {
	s.arr = append(s.arr, data)
	return len(s.arr) - 1
}

func (s *Stack) Pop() byte {
	data := s.arr[s.LastIndex()]
	s.arr = s.arr[:s.LastIndex()-1]
	return data
}

func (s *Stack) Last() byte {
	return s.arr[s.LastIndex()]
}

func (s *Stack) LastIndex() int {
	return len(s.arr) - 1
}

type VM struct {
	data  []byte
	ip    int // instruction pointer
	stack *Stack
	sp    int // stack pointer
}

func NewVM(data []byte) *VM {
	return &VM{
		data:  data,
		ip:    0,
		stack: &Stack{},
		sp:    -1,
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
	case InstrPush:
		vm.pushStack(vm.data[vm.ip-1])
	case InstrAdd:
		a := vm.stack.arr[0]
		b := vm.stack.arr[1]
		c := a + b
		vm.pushStack(c)
	}
	return nil
}

func (vm *VM) pushStack(b byte) {
	vm.sp++
	vm.stack.Push(b)
}
