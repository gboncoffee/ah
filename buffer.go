package main

import "bufio"

type Buffer interface {
	Displacement(line int) int
	DisplacedReader(disp int) *bufio.Reader
	Lines() int
	Insert(disp int, content string)
}
