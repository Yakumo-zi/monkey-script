package compiler

import (
	"fmt"
	"interpreter/ast"
	"interpreter/object"
	"vm/code"
)

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type EmittedInstruction struct {
	OpCode   code.Opcode
	Position int
}

func NewCompiler() *Compiler {
	return &Compiler{
		instructions:        code.Instructions{},
		constants:           []object.Object{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{OpCode: op, Position: pos}
	c.previousInstruction = previous
	c.lastInstruction = last
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "-":
			c.emit(code.OpMinus)
		case "!":
			c.emit(code.OpBang)

		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return nil
			}
			err = c.Compile(node.Left)
			if err != nil {
				return nil
			}
			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return nil
		}
		err = c.Compile(node.Right)
		if err != nil {
			return nil
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		case ">":
			c.emit(code.OpGreaterThan)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		pos := c.emit(code.OpJumpNotTruthy, 9999)
		err = c.Compile(node.Then)
		if err != nil {
			return err
		}
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}
		afterThenPos := len([]byte(c.instructions))
		c.changeOperand(pos, afterThenPos)

		if node.Else != nil {
			jumpPos := c.emit(code.OpJump, 9999)
			afterThenPos := len([]byte(c.instructions))
			c.changeOperand(pos, afterThenPos)
			if err = c.Compile(node.Else); err != nil {
				return err
			}
			if c.lastInstructionIsPop() {
				c.removeLastPop()
			}
			afterElsePos := len([]byte(c.instructions))
			c.changeOperand(jumpPos, afterElsePos)
		}
		return nil
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
		return nil
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	}
	return nil
}
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.instructions[opPos])
	newInstrunction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstrunction)

}

func (c *Compiler) replaceInstruction(pos int, newInstrunction []byte) {
	for i := 0; i < len(newInstrunction); i++ {
		c.instructions[pos+i] = newInstrunction[i]
	}
}

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.OpCode == code.OpPop
}
func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction

}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
