package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

const (
	OpAdd uint8 = iota + 1
	OpSub
	OpRight
	OpLeft
	OpOut
	OpIn
	OpLoopStart
	OpLoopEnd
)

type Interpreter struct {
	Code     *opcode
	Register [32768]uint8
}

type opcode struct {
	next [2]*opcode // first is the next one, second exists if it's a loop
	Op   uint8
	Val  uint8
}

func (i *Interpreter) LoadCode(code []byte) {
	var loopstack []*opcode
	var lastop *opcode
	for cnt, rawop := range code {
		op := new(opcode)
		op.Val = 1
		switch rawop {
		case '>':
			op.Op = OpRight
		case '<':
			op.Op = OpLeft
		case '+':
			op.Op = OpAdd
		case '-':
			op.Op = OpSub
		case '.':
			op.Op = OpOut
		case ',':
			op.Op = OpIn
		case '[':
			op.Op = OpLoopStart
			loopstack = append(loopstack, op)
		case ']':
			op.Op = OpLoopEnd
			if len(loopstack) > 0 {
				op.next[1] = loopstack[len(loopstack)-1]
				loopstack = loopstack[:len(loopstack)-1]
				op.next[1].next[1] = op
			}
		}
		if cnt == 0 {
			i.Code = op
			lastop = op
		} else {
			if op.Op <= 4 && lastop.Op == op.Op {
				lastop.Val++
			} else {
				lastop.next[0] = op
				lastop = op
			}
		}
	}
}

func (i *Interpreter) CPU() {
	currentop := i.Code
	var registerindex uint16

	for {
		switch currentop.Op {
		case OpAdd:
			i.Register[registerindex] += currentop.Val
		case OpSub:
			i.Register[registerindex] -= currentop.Val
		case OpRight:
			registerindex += uint16(currentop.Val)
		case OpLeft:
			registerindex -= uint16(currentop.Val)
		case OpOut:
			fmt.Printf("%c", rune(i.Register[registerindex]))
		case OpIn:
			_, err := fmt.Scanf("%c", &i.Register[registerindex])
			if err != nil {
				panic(err)
			}
		case OpLoopStart:
			if i.Register[registerindex] == 0 {
				currentop = currentop.next[1]
				continue
			}
		case OpLoopEnd:
			if i.Register[registerindex] != 0 {
				currentop = currentop.next[1]
				continue
			}
		}
		if currentop.next[0] == nil {
			break
		}
		currentop = currentop.next[0]

	}
}

func main() {
	flag.Parse()
	code, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	var i Interpreter
	i.LoadCode([]byte(code))
	start := time.Now()
	i.CPU()
	runtime := time.Since(start).Round(time.Millisecond).String()
	fmt.Println("\n===============================\nRUNTIME:", runtime)
}
