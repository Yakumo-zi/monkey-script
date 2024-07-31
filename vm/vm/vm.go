package vm

import (
	"fmt"
	"interpreter/object"
	"vm/code"
	"vm/compiler"
)

const StackSize = 2048

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

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

func (v *VM) LastPoppedStackElem() object.Object {
	return v.stack[v.sp]
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
		case code.OpDiv, code.OpAdd, code.OpMul, code.OpSub, code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := v.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			_, err := v.pop()
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := v.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := v.push(False)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (v *VM) executeBinaryOperation(op code.Opcode) error {
	right, err := v.pop()
	if err != nil {
		return err
	}
	left, err := v.pop()
	if err != nil {
		return err
	}
	leftType := left.Type()
	rightType := right.Type()
	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return v.executeIntegerOperation(op, left, right)
	}
	if leftType == object.BOOLEAN_OBJ && rightType == object.BOOLEAN_OBJ {
		return v.executeBooleanOperation(op, left, right)
	}
	return nil

}

func (v *VM) executeBooleanOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value
	if op == code.OpEqual {
		v.push(&object.Boolean{Value: leftValue == rightValue})

	} else {
		v.push(&object.Boolean{Value: leftValue != rightValue})
	}
	return nil
}

func (v *VM) executeIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value
	var err error
	switch op {
	case code.OpAdd:
		err = v.push(&object.Integer{Value: leftValue + rightValue})
	case code.OpSub:
		err = v.push(&object.Integer{Value: leftValue - rightValue})
	case code.OpMul:
		err = v.push(&object.Integer{Value: leftValue * rightValue})
	case code.OpDiv:
		if rightValue == 0 {
			return fmt.Errorf("can't div zero")
		}
		err = v.push(&object.Integer{Value: leftValue / rightValue})
	case code.OpEqual:
		if leftValue == rightValue {
			err = v.push(True)
		} else {
			err = v.push(False)
		}
	case code.OpGreaterThan:
		if leftValue > rightValue {
			err = v.push(True)
		} else {
			err = v.push(False)
		}
	case code.OpNotEqual:
		if leftValue == rightValue {
			err = v.push(False)
		} else {
			err = v.push(True)
		}
	}
	return err
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
