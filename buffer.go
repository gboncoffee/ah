package main

type Buffer interface {
	Insert(idx int, r rune) error
	Delete(idx int) error
	Get(idx int) (rune, error)
	Size() int
	Undo() (int, error)
	Redo() (int, error)
}
