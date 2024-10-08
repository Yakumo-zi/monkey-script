package vm

import (
	"fmt"
	"interpreter/object"

	"vm/code"
	"vm/compiler"
)

const StackSize = 2048
const MaxFrames = 1024
const GlobalSize = 65536

var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

type VM struct {
	constants  []object.Object
	stack      []object.Object
	globals    []object.Object
	frames     []*Frame
	frameIndex int
	sp         int
}

func NewVM(bytecode *compiler.ByteCode) *VM {

	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame
	return &VM{
		constants:  bytecode.Constants,
		stack:      make([]object.Object, StackSize),
		sp:         0,
		globals:    make([]object.Object, GlobalSize),
		frames:     frames,
		frameIndex: 1,
	}
}

func (v *VM) currentFrame() *Frame {
	return v.frames[v.frameIndex-1]
}

func (v *VM) pushFrame(f *Frame) {
	v.frames[v.frameIndex] = f
	v.frameIndex++
}

func (v *VM) popFrame() *Frame {
	v.frameIndex--
	return v.frames[v.frameIndex]
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
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for v.currentFrame().ip < len(v.currentFrame().Instructions())-1 {
		v.currentFrame().ip++
		ip = v.currentFrame().ip
		ins = v.currentFrame().Instructions()
		op = code.Opcode(ins[ip])
		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			v.currentFrame().ip += 2
			err := v.push(v.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpDiv, code.OpAdd, code.OpMul, code.OpSub, code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := v.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpBang:
			err := v.executeBangOpeartor()
			if err != nil {
				return err
			}

		case code.OpMinus:
			err := v.executeMinusOperator()
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
		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			v.currentFrame().ip = pos - 1
		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			v.currentFrame().ip += 2
			condition, err := v.pop()
			if err != nil {
				return err
			}
			if !isTruthy(condition) {
				v.currentFrame().ip = pos - 1
			}
		case code.OpNull:
			err := v.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			idx := code.ReadUint16(ins[ip+1:])
			v.currentFrame().ip += 2
			val, err := v.pop()
			if err != nil {
				return err
			}
			v.globals[idx] = val
		case code.OpGetGlobal:
			idx := code.ReadUint16(ins[ip+1:])
			v.currentFrame().ip += 2
			val := v.globals[idx]
			err := v.push(val)
			if err != nil {
				return err
			}
		case code.OpArray:
			num := code.ReadUint16(ins[ip+1:])
			v.currentFrame().ip += 2
			arr := &object.ArrayObject{Elements: make([]object.Object, num)}
			for num > 0 {
				elm, err := v.pop()
				if err != nil {
					return err
				}
				arr.Elements[num-1] = elm
				num -= 1
			}
			err := v.push(arr)
			if err != nil {
				return err
			}
		case code.OpHash:
			num := code.ReadUint16(ins[ip+1:])
			v.currentFrame().ip += 2
			err := v.buildHashPairs(int(num))
			if err != nil {
				return err
			}
		case code.OpIndex:
			idx, err := v.pop()
			if err != nil {
				return err
			}
			left, err := v.pop()
			if err != nil {
				return err
			}
			err = v.executeIndexExpression(left, idx)
			if err != nil {
				return err
			}
		case code.OpReturnValue:
			returnValue, err := v.pop()
			if err != nil {
				return err
			}
			frame := v.popFrame()
			v.sp = frame.basePointer - 1
			err = v.push(returnValue)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := v.popFrame()
			// 弹出CompiledFunction object
			v.sp = frame.basePointer - 1

			err := v.push(Null)
			if err != nil {
				return err
			}

		case code.OpCall:
			numArgs := ins[ip+1]
			v.currentFrame().ip += 1
			err := v.callFunction(int(numArgs))
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := ins[ip+1]
			v.currentFrame().ip += 1
			frame := v.currentFrame()
			v.stack[frame.basePointer+int(localIndex)], _ = v.pop()
		case code.OpGetLocal:
			localIndex := ins[ip+1]
			v.currentFrame().ip += 1
			frame := v.currentFrame()
			err := v.push(v.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (v *VM) callFunction(numArgs int) error {
	fn, ok := v.stack[v.sp-1-int(numArgs)].(*object.CompiledFunction)
	if !ok {
		return fmt.Errorf("calling non-funcion")
	}
	if numArgs != fn.NumParameters {
		return fmt.Errorf("wrong number of arguments, want %d, got %d", fn.NumParameters, numArgs)
	}
	frame := NewFrame(fn, v.sp-numArgs)
	v.pushFrame(frame)
	v.sp = frame.basePointer + fn.NumLocals
	return nil
}

func (v *VM) executeIndexExpression(left, idx object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && idx.Type() == object.INTEGER_OBJ:
		return v.executeArrayIndex(left, idx)
	case left.Type() == object.HASH_OBJ:
		return v.executeHashIndex(left, idx)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}

func (v *VM) executeArrayIndex(left, idx object.Object) error {
	arr := left.(*object.ArrayObject)
	i := idx.(*object.Integer).Value
	max := int64(len(arr.Elements) - 1)
	if i < 0 || i > max {
		return v.push(Null)
	}
	return v.push(arr.Elements[i])
}

func (v *VM) executeHashIndex(left, idx object.Object) error {
	hash := left.(*object.HashObject)
	key, ok := idx.(object.Hashable)
	if !ok {
		return fmt.Errorf("unusable as hash key: %s", idx.Type())
	}
	pair, ok := hash.Pairs[key.HashKey()]
	if !ok {
		return v.push(Null)
	}
	return v.push(pair.Value)
}

func (v *VM) buildHashPairs(num int) error {
	hash := &object.HashObject{Pairs: make(map[object.HashKey]object.HashPair)}
	for num > 0 {
		value, err := v.pop()
		if err != nil {
			return err
		}
		key, err := v.pop()
		if err != nil {
			return err
		}
		switch hashable := key.(type) {
		case object.Hashable:
			hash.Pairs[hashable.HashKey()] = object.HashPair{Key: key, Value: value}
		default:
			return fmt.Errorf("type %s can't support hashable", hashable.Type())
		}
		num -= 2
	}
	return v.push(hash)
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
func (v *VM) executeMinusOperator() error {
	operand, err := v.pop()
	if err != nil {
		return err
	}
	if result, ok := operand.(*object.Integer); ok {
		return v.push(&object.Integer{Value: -result.Value})
	}
	return fmt.Errorf("unsupported type:%s for negation operator", operand.Type())
}

func (v *VM) executeBangOpeartor() error {
	operand, err := v.pop()
	if err != nil {
		return err
	}
	switch operand {
	case True:
		return v.push(False)
	case False:
		return v.push(True)
	case Null:
		return v.push(True)
	default:
		return v.push(True)
	}

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
	if leftType == object.STRING_OBJ && rightType == object.STRING_OBJ {
		return v.executeStringOperation(op, left, right)
	}
	return fmt.Errorf("type misstach,%s %s %s", string(leftType), string(op), string(rightType))
}
func (v *VM) executeStringOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.StringObject).Value
	rightValue := right.(*object.StringObject).Value
	switch op {
	case code.OpAdd:
		return v.push(&object.StringObject{Value: leftValue + rightValue})

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
		return nil, fmt.Errorf("nothing in stack")
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
