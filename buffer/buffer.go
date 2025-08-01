// Package buffer implements an efficient PieceTable generic buffer with
// infinite undo/redo capabilities.
package buffer

import (
	"errors"
	"os"
	"slices"
	"unsafe"
)

var ErrorBottomOfUndoList = errors.New("reached bottom of undo list")
var ErrorTopOfUndoList = errors.New("reached top of undo list")
var ErrorOutOfBounds = errors.New("out of bounds")

// Using slices for representing pieces was kinda weird so I didn't.

var bufferSize = os.Getpagesize()

type Buffer[Content any] struct {
	// The first buffer never changes and does not respect the buffer size.
	buffers []backingBuffer[Content]
	pieces  []piece
	//edits   []edit
	size int // Cache
	//undoTop int
}

type piece struct {
	buffer int
	start  int
	length int
}

//type edit struct {
//	piece    piece
//	idx      int
//	deletion bool
//}

func New[Content any]() *Buffer[Content] {
	buffer := new(Buffer[Content])
	buffer.buffers = make([]backingBuffer[Content], 1)
	buffer.buffers[0] = newBackingBuffer[Content](buffer.bufferSize(), 0)
	buffer.pieces = append(buffer.pieces, piece{
		buffer: 0,
		start:  0,
		length: 0,
	})
	return buffer
}

func FromSlice[Content any](content []Content) *Buffer[Content] {
	buffer := new(Buffer[Content])
	buffer.buffers = make([]backingBuffer[Content], 2)

	// Here the memory we alloc is exactly the needed.
	buffer.buffers[0] = newBackingBuffer[Content](len(content), 0)
	buffer.buffers[1] = newBackingBuffer[Content](
		buffer.bufferSize(), len(content))

	for _, c := range content {
		buffer.buffers[0].append(c)
		buffer.size++
	}
	buffer.pieces = append(buffer.pieces, piece{
		buffer: 0,
		start:  0,
		length: buffer.buffers[0].size(),
	})

	return buffer
}

func Content[Content any](b *Buffer[Content]) []Content {
	content := make([]Content, 0, b.Size())

	for _, piece := range b.pieces {
		c := b.pieceContent(piece)
		content = append(content, c...)
	}

	return content
}

func (b *Buffer[Content]) Insert(idx int, r Content) error {
	pidx, disp, err := b.findPieceForInsertion(idx)
	if err != nil {
		return err
	}

	b.normalizeUndo()
	b.size++

	piec := b.pieces[pidx]
	buffer := len(b.buffers) - 1

	newPiece := piece{
		buffer: buffer,
		start:  b.buffers[buffer].size(),
		length: 1,
	}

	b.appendToBack(r)

	// If "appending" on the piece and the piece is pointing to the end of the
	// buffers, we literally append onto it.
	if disp == piec.length {
		buf, idx := b.indexByPiece(piec, piec.length-1)
		// Minus 2 because we need to exclude the item we just inserted.
		sm1 := b.buffers[buf].size() - 2
		if buf == buffer && idx == sm1 {
			b.pieces[pidx].length++
			return nil
		}

		// Else we insert the new piece.
		b.pieces = slices.Insert(b.pieces, pidx+1, newPiece)
		return nil
	}

	// If inserting in the beggining of a piece.
	if disp == 0 {
		b.pieces = slices.Insert(b.pieces, pidx, newPiece)
		return nil
	}

	// If inserting in the middle of a piece

	orig := b.pieces[pidx]

	// We make the existing piece the right one.
	b.pieces[pidx] = piece{
		buffer: orig.buffer,
		start:  orig.start + disp,
		length: orig.length - disp,
	}

	// Insert the new piece.
	b.pieces = slices.Insert(b.pieces, pidx, newPiece)

	// Insert the left piece.
	b.pieces = slices.Insert(b.pieces, pidx, piece{
		buffer: orig.buffer,
		start:  orig.start,
		length: disp,
	})

	return nil
}

func (b *Buffer[Content]) Delete(idx int) error {
	pidx, disp, err := b.findPieceWithIdx(idx)
	if err != nil {
		return err
	}

	b.normalizeUndo()
	b.size--

	piec := b.pieces[pidx]

	switch disp {
	// If removing from the top of the piece, we can simply decrease.
	case piec.length - 1:
		b.pieces[pidx].length--
	// If removing from the beggining of the piece, we can simply increase the
	// start.
	case 0:
		b.pieces[pidx].start++
		b.pieces[pidx].length--

		// If the piece begins at the end of the buffer.
		if b.pieces[pidx].start == b.buffers[b.pieces[pidx].buffer].size() {
			b.pieces[pidx].buffer++
			b.pieces[pidx].start = 0
		}
	default:
		// If we need to split the piece, we insert to the right.
		newb, newbdisp := b.indexByPiece(b.pieces[pidx], disp+1)
		newPiece := piece{
			buffer: newb,
			start:  newbdisp,
			length: b.pieces[pidx].length - (disp + 1),
		}
		b.pieces[pidx].length = disp
		b.pieces = slices.Insert(b.pieces, pidx+1, newPiece)
	}

	// If the length of the piece now is 0, we can remove it.
	if b.pieces[pidx].length == 0 {
		b.pieces = slices.Delete(b.pieces, pidx, pidx+1)
		// When removing, we may have sequential pieces. We can merge then.
		if pidx > 0 {
			b.tryMergePieces(pidx-1, pidx)
		}
	}

	return nil
}

func (b *Buffer[Content]) tryMergePieces(p1i, p2i int) {
	p1 := b.pieces[p1i]
	p2 := b.pieces[p2i]
	p1endbuf, p1enddisp := b.indexByPiece(p1, p1.length-1)

	// If the p1end is at the end of a buffer, we have to check wether the p2
	// begin is at the begin of the next one.
	if p1enddisp == b.buffers[p1endbuf].size()-1 {
		if p2.start == 0 && p2.buffer == p1endbuf+1 {
			b.mergePieces(p1i, p2i)
		}
	} else if p1endbuf == p2.buffer && p1enddisp == p2.start-1 {
		b.mergePieces(p1i, p2i)
	}
}

func (b *Buffer[Content]) mergePieces(p1i, p2i int) {
	removed := b.pieces[p2i]
	b.pieces = slices.Delete(b.pieces, p2i, p2i+1)
	b.pieces[p1i].length += removed.length
}

func (b *Buffer[Content]) Get(idx int) (Content, error) {
	var zero Content
	piec, disp, err := b.findPieceWithIdx(idx)
	if err != nil {
		return zero, err
	}

	buf, d := b.indexByPiece(b.pieces[piec], disp)
	return b.buffers[buf].content[d], nil
}

func (b *Buffer[Content]) Size() int {
	return b.size
}

func (b *Buffer[Content]) Undo() (int, error) {
	panic("todo")
}

func (b *Buffer[Content]) Redo() (int, error) {
	panic("todo")
}

func (b *Buffer[Content]) normalizeUndo() {
	//panic("todo")
	//b.edits = b.edits[:b.undoTop]
}

func (b *Buffer[Content]) findPieceWithIdx(idx int) (i int, d int, err error) {
	disp := 0
	for i, piece := range b.pieces {
		ndisp := piece.length + disp
		if ndisp > idx {
			return i, idx - disp, nil
		}
		disp = ndisp
	}

	return 0, 0, ErrorOutOfBounds
}

func (b *Buffer[Content]) findPieceForInsertion(idx int) (i int, d int, err error) {
	disp := 0
	for i, piece := range b.pieces {
		ndisp := piece.length + disp
		if ndisp >= idx {
			return i, idx - disp, nil
		}
		disp = ndisp
	}

	return 0, 0, ErrorOutOfBounds
}

func (b *Buffer[Content]) pieceContent(p piece) []Content {
	arr := make([]Content, 0, p.length)
	buf := p.buffer
	bdisp := p.start
	for range p.length {
		if bdisp >= b.buffers[buf].size() {
			bdisp = 0
			buf++
		}
		arr = append(arr, b.buffers[buf].content[bdisp])
		bdisp++
	}
	return arr
}

func (b *Buffer[Content]) indexByPiece(p piece, d int) (buffer int, bdisp int) {
	// If in the first (piece) buffer.
	if p.start+d < b.buffers[p.buffer].size() {
		return p.buffer, d + p.start
	}

	disp := b.buffers[p.buffer].size() - p.start
	buf := p.buffer + 1
	for {
		newdisp := disp + b.buffers[buf].size()
		if newdisp > d {
			return buf, newdisp - d
		}
		buf++
		disp = newdisp
	}
}

func (b *Buffer[Content]) bufferSize() int {
	var zero Content
	return bufferSize / int(unsafe.Sizeof(zero))
}
