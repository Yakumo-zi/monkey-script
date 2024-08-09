package compiler

import (
	"fmt"
	"interpreter/ast"
	"interpreter/object"
	"sort"
	"vm/code"
)

type EmittedInstruction struct {
	OpCode   code.Opcode
	Position int
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	symboltable *SymbolTable
	constants   []object.Object
	scopes      []CompilationScope
	scopeIndex  int
}

func NewCompiler() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	return &Compiler{
		constants:   []object.Object{},
		scopes:      []CompilationScope{mainScope},
		symboltable: NewSymbolTable(),
	}
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.symboltable = NewEnclosedSymbolTable(c.symboltable)
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
}
func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstrunction()
	c.symboltable = c.symboltable.Outer
	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	return instructions
}

func (c *Compiler) currentInstrunction() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}
func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{OpCode: op, Position: pos}
	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
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
	case *ast.LetStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		sym := c.symboltable.Define(node.Name.Value)
		if sym.Scope == GlobalScope {

			c.emit(code.OpSetGlobal, sym.Index)
		} else {
			c.emit(code.OpSetLocal, sym.Index)
		}
	case *ast.ArrayLiteral:
		for _, e := range node.Elements {
			err := c.Compile(e)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})
		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}
		c.emit(code.OpHash, len(node.Pairs)*2)
	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)
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
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}
		jumpPos := c.emit(code.OpJump, 9999)
		afterThenPos := len([]byte(c.currentInstrunction()))
		c.changeOperand(pos, afterThenPos)
		if node.Else != nil {
			if err = c.Compile(node.Else); err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		} else {
			c.emit(code.OpNull)
		}
		afterElsePos := len([]byte(c.currentInstrunction()))
		c.changeOperand(jumpPos, afterElsePos)
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.Identifier:
		sym, ok := c.symboltable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("variable %s not define", node.Value)
		}
		if sym.Scope == GlobalScope {
			c.emit(code.OpGetGlobal, sym.Index)
		} else {
			c.emit(code.OpGetLocal, sym.Index)
		}
	case *ast.StringLiteral:
		stringLiteral := &object.StringObject{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(stringLiteral))
	case *ast.FunctionLiteral:
		c.enterScope()
		for _, ident := range node.Parameters {
			err := c.Compile(ident)
			if err != nil {
				return err
			}
		}
		err := c.Compile(node.Body)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}
		numLocals := c.symboltable.numDefinitions
		instructions := c.leaveScope()
		compileFn := &object.CompiledFunction{Instructions: instructions, NumLocals: numLocals}
		c.emit(code.OpConstant, c.addConstant(compileFn))
	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return nil
		}
		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}
		c.emit(code.OpCall)
	}
	return nil
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))
	c.scopes[c.scopeIndex].lastInstruction.OpCode = code.OpReturnValue
}
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstrunction()[opPos])
	newInstrunction := code.Make(op, operand)
	c.replaceInstruction(opPos, newInstrunction)
}

func (c *Compiler) replaceInstruction(pos int, newInstrunction []byte) {
	ins := c.currentInstrunction()
	for i := 0; i < len(newInstrunction); i++ {
		ins[pos+i] = newInstrunction[i]
	}
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstrunction()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastInstruction.OpCode == op
}
func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction
	old := c.currentInstrunction()
	new := old[:last.Position]
	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)
	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstrunction())
	updateInstuction := append(c.currentInstrunction(), ins...)
	c.scopes[c.scopeIndex].instructions = updateInstuction
	return posNewInstruction

}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) ByteCode() *ByteCode {
	return &ByteCode{
		Instructions: c.currentInstrunction(),
		Constants:    c.constants,
	}
}

type ByteCode struct {
	Instructions code.Instructions
	Constants    []object.Object
}
