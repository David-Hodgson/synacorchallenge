package main

import (
	"fmt"
)

const (
	add = uint16(9)
	out = uint16(19)
)

var registers = make(map[uint16]uint16)

//TODO registers should be a type and this method should take that type
func clearRegisters() {
	registers = make(map[uint16]uint16)

}

func getValue(value uint16) uint16 {

	if value > 32767 {
		return registers[value]
	}

	return value
}

func runProgram(program []uint16) {

	programCounter :=0
	clearRegisters()

	for ;; {

		if programCounter < 0 || programCounter > len(program)-1 {
			break;
		}

		instruction := program[programCounter]
		//fmt.Println("Command:", instruction)

		switch instruction {
		case add:
			// add: 9 a b c
			// assign into <a> the sum of <b> and <c> (modulo 32768)
			register := program[programCounter+1]
			b := getValue(program[programCounter+2])
			c := getValue(program[programCounter+3])
			sum := (b+c) % 32768
			registers[register] = sum
			programCounter += 4
		case out:
			// out: 19 a
			// write the character represented by ascii code <a> to the terminal
			register := program[programCounter+1]
			character := string(getValue(registers[register]))
			fmt.Print(character)
			programCounter += 2
		default:
			fmt.Println("Unknown instruction:",instruction)
			programCounter++
		}

	}

}

func main() {

	fmt.Println("Synacor Challenge")

	//My first program

	//should do

	/*

	Store into register 0 the sum of 4 and the value contained in register 1.
  - Output to the terminal the character with the ascii code contained in register 0.


	 */
	program := []uint16 {9,32768,32769,4,19,32768}

	fmt.Println("Input program: ", program)

	runProgram(program)
}
