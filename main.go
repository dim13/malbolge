package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode"
)

type Memory [59049]int // 3^10

func ReadProg(in io.Reader) Memory {
	var m Memory
	buf, err := ioutil.ReadAll(in)
	if err != nil {
		log.Fatal(err)
	}
	var n int
	for _, v := range buf {
		if !unicode.IsSpace(rune(v)) {
			m[n] = int(v)
			n++
		}
	}
	for n < len(m) {
		m[n] = Crazy(m[n-2], m[n-1])
		n++
	}
	return m
}

func Eval(mem Memory) {
	var A, C, D int
	for {
		switch (mem[C] + C) % 94 {
		case 4: // jmp [d]
			C = mem[D]
		case 5: // out a
			fmt.Printf("%c", byte(A))
		case 23: // in a
			fmt.Scanf("%c", &A)
		case 39: // rotr [d]; mov a, [d]
			mem[D] = RotR(mem[D])
			A = mem[D]
		case 40: // mov d, [d]
			D = mem[D]
		case 62: // crz [d], a; mov a, [d]
			mem[D] = Crazy(mem[D], A)
			A = mem[D]
		case 68: // nop
		case 81: // end
			return
		}
		mem[C] = Encrypt(mem[C])
		C = (C + 1) % len(mem)
		D = (D + 1) % len(mem)
	}
}

func RotR(x int) int {
	return x/3 + (x%3)*19683
}

//		crz	Input 2
//			0	1	2
// Input 1	0	1	0	0
//		1	1	0	2
//		2	2	2	1

func Crazy(x, y int) (ret int) {
	o := [][]int{
		{4, 3, 3, 1, 0, 0, 1, 0, 0},
		{4, 3, 5, 1, 0, 2, 1, 0, 2},
		{5, 5, 4, 2, 2, 1, 2, 2, 1},
		{4, 3, 3, 1, 0, 0, 7, 6, 6},
		{4, 3, 5, 1, 0, 2, 7, 6, 8},
		{5, 5, 4, 2, 2, 1, 8, 8, 7},
		{7, 6, 6, 7, 6, 6, 4, 3, 3},
		{7, 6, 8, 7, 6, 8, 4, 3, 5},
		{8, 8, 7, 8, 8, 7, 5, 5, 4},
	}
	for _, v := range []int{1, 9, 81, 729, 6561} {
		ret += o[x/v%9][y/v%9] * v
	}
	return
}

func Encrypt(x int) int {
	e := []int{
		57, 109, 60, 46, 84, 86, 97, 99, 96, 117,
		89, 42, 77, 75, 39, 88, 126, 120, 68, 108,
		125, 82, 69, 111, 107, 78, 58, 35, 63, 71,
		34, 105, 64, 53, 122, 93, 38, 103, 113, 116,
		121, 102, 114, 36, 40, 119, 101, 52, 123, 87,
		80, 41, 72, 45, 90, 110, 44, 91, 37, 92,
		51, 100, 76, 43, 81, 59, 62, 85, 33, 112,
		74, 83, 55, 50, 70, 104, 79, 65, 49, 67,
		66, 54, 118, 94, 61, 73, 95, 48, 47, 56,
		124, 106, 115, 98,
	}
	return e[x%94]
}

func main() {
	file := flag.String("file", "", "malbolge program file")
	flag.Parse()

	fd, err := os.Open(*file)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	Eval(ReadProg(fd))
}
