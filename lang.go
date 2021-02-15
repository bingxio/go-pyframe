//
// Copyright (c) 2021 bingxio. All rights reserved.
//

package main

import (
    "fmt"
    "errors"
)

// my custom stack structure
type Stack struct {
    Val []interface{}
}

// push the new value forward
func (s *Stack) Push(val interface{}) { s.Val = append(s.Val, val) }

// pop up the value at the top of stack
func (s *Stack) Pop(reverse bool) interface{} {
    if s.Len() == 0 {
        return nil
    }
    if !reverse {
        // first in, then out
        // FILO
        i := len(s.Val) - 1
        v := s.Val[i]
        s.Val = append(s.Val[:i], s.Val[i+1:]...)
        return v
    } else {
        // first in, first out
        // FIFO
        v := s.Val[0]
        s.Val = append(s.Val[:0], s.Val[1:]...)
        return v
    }
}

// return the length of stack
func (s Stack) Len() int { return len(s.Val) }

// return whether the stack value is empty
func (s Stack) Empty() bool { return len(s.Val) == 0 }

// copy stack from source
func (s *Stack) CopyOf(v Stack) { s.Val = append(s.Val, v.Val[0:v.Len()]...) }

// enumeration
type ObjType byte

const (
    Num  ObjType = iota // integer
    Str                 // string
    Func                // function
    None                // none
)

// Value in virtual machine
type Object interface {
    Type() ObjType      // type
    Value() interface{} // value
    Stringer() string   // stringer
}

// number
type NumObj struct{ v float64 }

func (i NumObj) Type() ObjType      { return Num }
func (i NumObj) Value() interface{} { return i.v }

func (i NumObj) Stringer() string {
    return fmt.Sprintf("<Num { Value = %f }>", i.v)
}

// string
type StrObj struct{ v string }

func (s StrObj) Type() ObjType      { return Str }
func (s StrObj) Value() interface{} { return s.v }

func (s StrObj) Stringer() string {
    return fmt.Sprintf("<Str { Value = '%s' }>", s.v)
}

// function
type FuncObj struct {
    name string
    args int
    code *Code
}

func (f FuncObj) Type() ObjType      { return Func }
func (f FuncObj) Value() interface{} { return f.name }

func (f FuncObj) Stringer() string {
    return fmt.Sprintf("<Func { Name = '%s' Args = %d }>", f.name, f.args)
}

// none
type NoneObj struct{}

func (_ NoneObj) Type() ObjType      { return None }
func (_ NoneObj) Value() interface{} { return "None" } // return a string for the prompt
func (_ NoneObj) Stringer() string   { return "<None>" }

// the four operations for two object
// them must be integer
func FourOperations(x, y Object, op byte) Object {
    if x.Type() == None || y.Type() == None {
        return nil
    }
    if x.Type() == Str || y.Type() == Str {
        return nil
    }
    a := x.Value().(float64)
    b := y.Value().(float64)
    //
    switch op {
    case OpAdd:
        return NumObj{a + b}
    case OpSub:
        return NumObj{a - b}
    case OpMul:
        return NumObj{a * b}
    case OpDiv:
        return NumObj{a / b}
    default:
        return nil
    }
}

// code object structure
type Code struct {
    bytes   Stack // bytecodes
    cons    Stack // constants
    vars    Stack // variables
    offsets Stack // type index, not after each bytecode
}

// make a code object
func NewCode(code Stack, cons Stack, vars Stack, offsets Stack) *Code {
    c := &Code{
        bytes:   code,
        cons:    cons,
        vars:    vars,
        offsets: offsets, // its reserved
    }
    return c
}

// symbol table
type Table map[string]Object

// find object by name
func (t Table) lookUp(name string) Object {
    val, ok := t[name]
    if !ok {
        return nil
    }
    return val
}

// return the length of symbol table
func (t Table) Len() int { return len(t) }

// frame structure
type Frame struct {
    code     *Code // code object
    stack    Stack // calculation stack
    localVar Table // local variables table
    builtin  Table // builtin objects
}

// make a frame
func NewFrame(code *Code) *Frame {
    f := &Frame{
        stack: Stack{
            Val: make([]interface{}, 0),
        },
        code:     code,
        localVar: Table{},
    }
    // debug information
    f.dissemble()
    return f
}

// output detailed data of frame
func (f Frame) dissemble() {
    fmt.Println("Bytecodes: ")
    // variable vars
    sp := 0
    // offset of variables
    vp := 0
    // offset of offsets for frame, self incrementing
    op := 0
    // code object for current frame
    co := f.code

    for i := 0; i < f.code.bytes.Len(); i++ {
        // byte
        b := co.bytes.Val[i].(byte)

        switch b {
        case OpLoad:
            fmt.Printf("%20d: %s %21s (%d)\n",
                i,
                stringCode(b),
                co.vars.Val[vp],
                co.offsets.Val[op],
            )
            op++
            vp++
        case OpConst:
            obj := co.cons.Val[co.offsets.Val[op].(int)].(Object)

            fmt.Printf("%20d: %s %20v (%d)\n",
                i,
                stringCode(b),
                obj.Value(),
                co.offsets.Val[op],
            )
            op++
        case OpStore:
            fmt.Printf("%20d: %s %20s (%d)\n",
                i,
                stringCode(b),
                co.vars.Val[sp],
                co.offsets.Val[op],
            )
            op++
            sp++
        case OpCall:
            fmt.Printf("%20d: %s %24d\n",
                i,
                stringCode(b),
                co.offsets.Val[op],
            )
            op++
        // the following bytecode has no parameters
        case
            OpAdd,
            OpSub,
            OpMul,
            OpDiv,
            OpPut,
            OpPop,
            OpBuiltin,
            OpRet,
            OpAddO,
            OpSubO:
            fmt.Printf("%20d: %s \n", i, stringCode(b))
        }
    }
    fmt.Println("Constants: ")
    {
        for k, v := range co.cons.Val {
            fmt.Printf("%20d: %s\n", k, v.(Object).Stringer())
        }
    }
    fmt.Println("Variables: ")
    {
        for k, v := range co.vars.Val {
            fmt.Printf("%20d: '%s'\n", k, v)
        }
    }
    fmt.Println("<----------------------------->")
}

// output detailed data of symbol table in frame
func (f Frame) dissembleTable() {
    fmt.Println("Symbols: ")
    {
        i := 0
        for k, v := range f.localVar {
            fmt.Printf("%20d: '%s' -> %s\n", i, k, v.Stringer())
            i++
        }
    }
    fmt.Println("<----------------------------->")
}

// bytecode
// each is a byte size
const (
    OpLoad    byte = iota // load variable
    OpConst               // load constant
    OpStore               // assign to variable
    OpAdd                 // addition
    OpSub                 // subtraction
    OpMul                 // multiplication
    OpDiv                 // division
    OpPut                 // output top object
    OpPop                 // pop value
    OpCall                // call function
    OpBuiltin             // load built in object
    OpAddO                // top of stack plus one
    OpSubO                // top of stack subtract one
    OpRet                 // return
)

// stringer for bytecode
func stringCode(b byte) string {
    switch b {
    case OpLoad:
        return "LOAD"
    case OpConst:
        return "CONST"
    case OpStore:
        return "STORE"
    case OpAdd:
        return "ADD"
    case OpSub:
        return "SUB"
    case OpMul:
        return "MUL"
    case OpDiv:
        return "DIV"
    case OpPut:
        return "PUT"
    case OpPop:
        return "POP"
    case OpCall:
        return "CALL"
    case OpBuiltin:
        return "BUILTIN"
    case OpAddO:
        return "ADDO"
    case OpSubO:
        return "SUBO"
    case OpRet:
        return "RET"
    }
    return "<?>"
}

// virtual machine structure
type Vm struct {
    frames []*Frame // the current frame of the vm
    last   Object   // last frame returned object
}

// return the top frame of vm
func (v Vm) topFrame() *Frame {
    return v.frames[len(v.frames)-1]
}

// before function call
// pop the top of frame in frames
func (v *Vm) popFrame() {
    l := len(v.frames) - 1
    v.frames = append(v.frames[:l], v.frames[l+1:]...)
}

// push new frame to frames of vm
func (v *Vm) pushFrame(f *Frame) { v.frames = append(v.frames, f) }

// return the reversed value at offset stack
func popOffset(c *Code) int {
    return c.offsets.Pop(true).(int)
}

// push new value to stack
func pushStack(f *Frame, v interface{}) { f.stack.Push(v) }

// return the value at the top of stack
func popStack(f *Frame) Object {
    return f.stack.Pop(false).(Object)
}

// return to the top of the stack just to see
func peekStack(f Frame) Object {
    return f.stack.Val[f.stack.Len()-1].(Object)
}

// evaluate
func (v *Vm) eval() error {
    f := v.topFrame()
    // offset of bytecode
    ip := 0
    // current code object of frame
    co := f.code

    for ip < co.bytes.Len() {
        // byte
        b := co.bytes.Val[ip].(byte)

        switch b {
        case OpLoad:
            nam := co.vars.Val[popOffset(co)].(string)
            val := f.localVar.lookUp(nam)
            // symbol not found
            if val == nil {
                return errors.New(
                    fmt.Sprintf("symbol '%s' not found", nam),
                )
            }
            pushStack(f, val)
        case OpConst:
            // constant load to stack
            pushStack(f, co.cons.Val[popOffset(co)])
        case OpStore:
            val := popStack(f)
            nam := co.vars.Val[popOffset(co)].(string)
            // add to symbol table
            f.localVar[nam] = val
        // + - * /
        case OpAdd, OpSub, OpMul, OpDiv:
            x := popStack(f)
            y := popStack(f)
            //
            val := FourOperations(x, y, b)
            if val == nil {
                return errors.New(
                    fmt.Sprintf(
                        "four arithmetic error \n\n"+
                            "OpA: %v\n"+
                            "OpB: %v",
                        x,
                        y,
                    ),
                )
            }
            pushStack(f, val)
        case OpAddO, OpSubO:
            val := popStack(f)
            //
            if val.Type() != Num {
                return errors.New("only integers can be evaluated")
            }
            obj := NumObj{val.(NumObj).v}

            if b == OpAddO {
                obj.v += 1
            }
            if b == OpSubO {
                obj.v -= 1
            }
            // restore
            pushStack(f, obj)
        case OpPut:
            fmt.Println("PUT: ", peekStack(*f).Value())
        case OpPop:
            // pop up the value at the top of the stack
            // nothing else!!
            popStack(f)
        case OpCall:
            val := popStack(f).(FuncObj)
            // new frame
            n := NewFrame(val.code)
            // copy the run value that must exist
            // for once again
            t := new(Stack)
            t.CopyOf(n.code.offsets)
            // function arguments
            if off := popOffset(co); off != 0 {
                if off != val.args {
                    return errors.New(
                        fmt.Sprintf(
                            "function parameter error \n\n"+
                                "Require: %d\n"+
                                "Have: %d",
                            val.args,
                            off,
                        ),
                    )
                }
                // data loading error
                if v.topFrame().stack.Len() < off {
                    return errors.New(
                        "data loading error",
                    )
                }
                // store variables to frame of function
                for i := 0; i < off; i++ {
                    // parameter name
                    x := val.code.vars.Val[i].(string)
                    // object
                    y := popStack(f)
                    // store to new frame
                    n.localVar[x] = y
                }
            }
            // push the new frame to the top
            v.pushFrame(n)
            // eval top frame
            if err := v.eval(); err != nil {
                panic(err)
            }
            // in order to do it again next time
            n.code.offsets = *t
            // destroy frame
            v.popFrame()
            // have return value
            if v.last != nil {
                // push current frame
                pushStack(f, v.last)
                // assign nil
                v.last = nil
            }
        case OpRet:
            // return value cache between frames
            v.last = popStack(f)
            // complete the execution
            v.dissemble()
            // output detailed data of symbol table in frame
            f.dissembleTable()
        }
        ip++
    }
    return nil
}

// output detailed data of vm
func (v Vm) dissemble() {
    fmt.Println("Stack: ")
    {
        for k, v := range v.topFrame().stack.Val {
            fmt.Printf("%20d: %v\n", k, v.(Object).Stringer())
        }
    }
}

func main() {}
