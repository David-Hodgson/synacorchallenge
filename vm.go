package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	halt = uint16(0)
	set = uint16(1)
	push = uint16(2)
	pop = uint16(3)
	eq = uint16(4)
	gt = uint16(5)
	jmp = uint16(6)
	jt = uint16(7)
	jf = uint16(8)
	add = uint16(9)
	mult = uint16(10)
	mod = uint16(11)
	and = uint16(12)
	or = uint16(13)
	not = uint16(14)
	rmem = uint16(15)
	wmem = uint16(16)
	call = uint16(17)
	ret = uint16(18)
	out = uint16(19)
	in = uint16(20)
	noop = uint16(21)
)


type Stack interface {

	Pop() uint16
	Push(uint16)
}

type stack struct {
	top int16
	values []uint16
}

func NewStack() stack{
	newStack := stack{-1,make([]uint16,100)}
	return newStack
}

func (s *stack) Pop() uint16 {
	if s.top > -1 {
		s.top--
		return s.values[s.top+1]
	}
	//TODO return err and let calling code handle it or have is empty method
	panic("Popping from an empty stack")
}

func (s *stack) Push(value uint16) {
	s.top += 1

	if int(s.top) > len(s.values)-1 {
		panic("Trying to add to full stack")
	}
	s.values[s.top] = value
}
/**

halt: 0
  stop execution and terminate the program
set: 1 a b
  set register <a> to the value of <b>
push: 2 a
  push <a> onto the stack
pop: 3 a
  remove the top element from the stack and write it into <a>; empty stack = error
eq: 4 a b c
  set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
gt: 5 a b c
  set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
jmp: 6 a
  jump to <a>
jt: 7 a b
  if <a> is nonzero, jump to <b>
jf: 8 a b
  if <a> is zero, jump to <b>
add: 9 a b c
  assign into <a> the sum of <b> and <c> (modulo 32768)
mult: 10 a b c
  store into <a> the product of <b> and <c> (modulo 32768)
mod: 11 a b c
  store into <a> the remainder of <b> divided by <c>
and: 12 a b c
  stores into <a> the bitwise and of <b> and <c>
or: 13 a b c
  stores into <a> the bitwise or of <b> and <c>
not: 14 a b
  stores 15-bit bitwise inverse of <b> in <a>
rmem: 15 a b
  read memory at address <b> and write it to <a>
wmem: 16 a b
  write the value from <b> into memory at address <a>
call: 17 a
  write the address of the next instruction to the stack and jump to <a>
ret: 18
  remove the top element from the stack and jump to it; empty stack = halt
out: 19 a
  write the character represented by ascii code <a> to the terminal
in: 20 a
  read a character from the terminal and write its ascii code to <a>; it can be assumed that once input starts, it will continue until a newline is encountered; this means that you can safely read whole lines from the keyboard and trust that they will be fully read
noop: 21
  no operation
 */

var registers = make(map[uint16]uint16)

//TODO registers should be a type and this method should take that type
func clearRegisters() {
	registers = make(map[uint16]uint16)

}

func getValue(value uint16) uint16 {

	if value > uint16(32767) {
		return registers[value]
	}

	return value
}

func runProgram(program []uint16) {

	programCounter :=0
	clearRegisters()
	stack := NewStack();

	for ;; {

		if programCounter < 0 || programCounter > len(program)-1 {
			break;
		}

		instruction := program[programCounter]

		switch instruction {
		case halt:
			//halt: 0
			//stop execution and terminate the program
			fmt.Println("Halting at programCounter:", programCounter)
			programCounter = -1
			break;
		case set:
			//set: 1 a b
			//set register <a> to the value of <b>
			a := program[programCounter+1]
			b := getValue(program[programCounter+2])
			registers[a] = b
			programCounter += 3
		case push:
			//push: 2 a
			//	push <a> onto the stack
			a := program[programCounter+1]
			stack.Push(getValue(a))
			programCounter += 2
		case pop:
			//pop: 3 a
			//	remove the top element from the stack and write it into <a>; empty stack = error
			a := program[programCounter+1]
			value := stack.Pop()
			registers[a] = value
			programCounter += 2
		case eq:
			//eq: 4 a b c
			//	set <a> to 1 if <b> is equal to <c>; set it to 0 otherwise
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])

			if b == c {
				registers[register] = 1
			} else {
				registers[register] = 0
			}
			programCounter += 4

		case gt:
			//gt: 5 a b c
			//	set <a> to 1 if <b> is greater than <c>; set it to 0 otherwise
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])

			if b > c {
				registers[register] = 1
			} else {
				registers[register] = 0
			}
			programCounter += 4

		case jmp:
			// jmp 6 a
			//jump to <a>
			register := program[programCounter+1]
			jmpValue := getValue(register)
			programCounter = int(jmpValue)
		case jt:
			//jt: 7 a b
			//if <a> is nonzero, jump to <b>
			a := getValue(program[programCounter+1])
			b := getValue(program[programCounter+2])

 			if a != 0 {
				programCounter = int(b)
			} else {
				programCounter += 3
			}
		case jf:
			//jf: 8 a b
			//if <a> is zero, jump to <b>
			a := getValue(program[programCounter+1])
			b := getValue(program[programCounter+2])

			if a == 0 {
				programCounter = int(b)
			} else {
				programCounter += 3
			}
		case add:
			// add: 9 a b c
			// assign into <a> the sum of <b> and <c> (modulo 32768)
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])
			sum := (b+c) % 32768
			registers[register] = sum
			programCounter += 4
		case mult:
			//mult: 10 a b c
			//store into <a> the product of <b> and <c> (modulo 32768)
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])
			product := (b*c) % 32768
			registers[register] = product
			programCounter += 4
		case mod:
			//mod: 11 a b c
			//store into <a> the remainder of <b> divided by <c>
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])
			mod := b %c
			registers[register] = mod
			programCounter += 4
		case and:
			//and: 12 a b c
			//	stores into <a> the bitwise and of <b> and <c>
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])
			value := b & c
			registers[register] = value
			programCounter += 4
		case or:
			//or: 13 a b c
			//stores into <a> the bitwise or of <b> and <c>
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])
			value := b | c
			registers[register] = value
			programCounter += 4
		case not:
			//not: 14 a b
			//stores 15-bit bitwise inverse of <b> in <a>
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			value := b ^ 32767
			registers[register] = value
			programCounter += 3
		case rmem:
			//rmem: 15 a b
			//read memory at address <b> and write it to <a>
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			value := program[b]
			registers[register] = value
			programCounter += 3
		case wmem:
			//wmem: 16 a b
			//write the value from <b> into memory at address <a>
			a := getValue(program[programCounter+1])
			b := getValue(program[programCounter+2])
			program[a] = b
 			programCounter += 3
		case call:
			//call: 17 a
			//write the address of the next instruction to the stack and jump to <a>
			a := getValue(program[programCounter+1])
			stack.Push(uint16(programCounter+2))
			programCounter = int(a)
		case ret:
			//ret: 18
			//remove the top element from the stack and jump to it; empty stack = halt
			programCounter = int(stack.Pop())
		case out:
			// out: 19 a
			// write the character represented by ascii code <a> to the terminal
			register := program[programCounter+1]
			character := string(getValue(register))
			fmt.Print(character)
			programCounter += 2
		case noop:
			// noop: 20
			// no operations
			programCounter += 1
		default:
			fmt.Println("Unknown instruction:",instruction, " at line ", programCounter)

			fmt.Println(programCounter-5, " - ", program[programCounter-5])
			fmt.Println(programCounter-4, " - ", program[programCounter-4])
			fmt.Println(programCounter-3, " - ", program[programCounter-3])
			fmt.Println(programCounter-2, " - ", program[programCounter-2])
			fmt.Println(programCounter-1, " - ", program[programCounter-1])
			fmt.Println(programCounter, " - ", program[programCounter])
			fmt.Println(programCounter+1, " - ", program[programCounter+1])
			fmt.Println(programCounter+2, " - ", program[programCounter+2])

			fmt.Println("Registers: ", registers)
			panic("Quiting")
			programCounter++
		}

	}

	fmt.Println("Program finished with progam counter at ", programCounter)
}

//TODO this whole method needs tidying up
func ReadBinaryFile(filename string) []uint16 {

	file, err := os.Open(filename) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	var contents []uint16

	if file != nil {
//		fi, _ := file.Stat()

		buffer := make([]byte, 1024)
		for {


			count, err := file.Read(buffer)
			if err != nil {
				break
			} else {
				//log.Println("Read something.",buffer)
				//contents += buffer[0:count]..
				buffer = buffer[0:count]
				var pi = make([]uint16, count/2)

				buf := bytes.NewReader(buffer)
				//fmt.Println(buf)
				err := binary.Read(buf, binary.LittleEndian, &pi)
				if err != nil {
					fmt.Println("binary.Read failed:", err)
				}

				//fmt.Println(contents)
				contents = append(contents, pi...)
			}

		}
	} else {
		fmt.Println("File is nil")
	}

	file.Close()
	return contents
}

func main() {

	fmt.Println("=========================================================")
	fmt.Println("Synacor Challenge")
	fmt.Println()

	fileName := "D:\\temp\\synacor\\challenge.bin"
	fmt.Println("Reading program from: ", fileName)
	program := ReadBinaryFile(fileName)
	fmt.Println("Program has ", len(program), "instructions")
	fmt.Println("=========================================================")
	fmt.Println()

	//TODO fudge for program length
	//TODO add proper memory handling
	if len(program) < 32768 {
		ram := make([]uint16, 32768 -len(program))
		program = append(program, ram...)
	}

	runProgram(program)

}
