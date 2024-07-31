package vm

import (
	"fmt"
	"interpreter/object"
	"vm/code"
	"vm/compiler"
)

const StackSize = 2048

type VM struct {
	constants    []object.Object
	instructions code.Instructions
	stack        []object.Object
	sp           int
}

func NewVM(bytecode *compiler.ByteCode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (v *VM) StackTop() object.Object {
	if v.sp == 0 {
		return nil
	}
	return v.stack[v.sp-1]
}

func (v *VM) Run() error {
	for ip := 0; ip < len(v.instructions); ip++ {
		op := code.Opcode(v.instructions[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(v.instructions[ip+1:])
			ip += 2
			err := v.push(v.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd:
			err := v.executeInfixExpression(code.OpAdd)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func (v *VM) executeInfixExpression(op code.Opcode) error {
	right, err := v.pop()
	if err != nil {
		return err
	}
	left, err := v.pop()
	if err != nil {
		return err
	}
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	switch op {
	case code.OpAdd:
		err = v.push(&object.Integer{Value: leftValue + rightValue})
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *VM) pop() (object.Object, error) {
	if v.sp <= 0 {
		return nil, fmt.Errorf("stack overflow")
	}
	ret := v.stack[v.sp-1]
	v.sp -= 1
	return ret, nil
}

func (v *VM) push(o object.Object) error {
	if v.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}
	v.stack[v.sp] = o
	v.sp += 1
	return nil
}
