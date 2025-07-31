package buffer

import (
	"errors"
	"slices"
	"strings"
)

// Using slices for representing pieces was kinda weird so I didn't.

const bufferSize = 4096

type Buffer struct {
	// The first buffer never changes and does not respect the buffer size.
	buffers [][]rune
	pieces  []piece
	edits   []edit
	size    int // Cache
}

type piece struct {
	buffer int
	start  int
	length int
}

type edit struct {
	piece    piece
	deletion bool
}

func FromString(content string) *Buffer {
	buffer := new(Buffer)
	buffer.buffers = make([][]rune, 2)

	// We make Go alloc a sane amount of memory (may be up to 4x more than we
	// actually need due to how UTF-8 works, but hey, we're doing only one
	// allocation, and who cares about virtual memory anyways?).
	buffer.buffers[0] = make([]rune, 0, len(content))
	buffer.buffers[1] = make([]rune, 0, bufferSize)

	for _, c := range content {
		buffer.buffers[0] = append(buffer.buffers[0], c)
		buffer.size++
	}
	buffer.pieces = append(buffer.pieces, piece{
		buffer: 0,
		start:  0,
		length: len(buffer.buffers[0]),
	})

	return buffer
}

func (b *Buffer) String() string {
	var builder strings.Builder
	builder.Grow((len(b.buffers)*(bufferSize-1) + len(b.buffers[0])) * 4)

	for _, piece := range b.pieces {
		content := b.pieceContent(piece)
		for _, r := range content {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

func (b *Buffer) Insert(idx int, r rune) error {
	pidx, disp, err := b.findPieceForInsertion(idx)
	if err != nil {
		return err
	}

	b.size++

	piec := b.pieces[pidx]
	buffer := len(b.buffers) - 1

	var ed *edit
	newEdit := edit{deletion: false, piece: piece{
		buffer: buffer,
		start:  len(b.buffers[buffer]),
		length: 0,
	}}
	if len(b.edits) <= 0 {
		b.edits = append(b.edits, newEdit)
		ed = &b.edits[len(b.edits)-1]
	} else {
		ed = &b.edits[len(b.edits)-1]
		if ed.deletion || ed.piece != piec || disp != piec.length {
			b.edits = append(b.edits, newEdit)
			ed = &b.edits[len(b.edits)-1]
		}
	}

	b.buffers[buffer] = append(b.buffers[buffer], r)
	if len(b.buffers[buffer]) == bufferSize {
		b.buffers = append(b.buffers, make([]rune, 0, bufferSize))
	}

	// If "appending" on the piece.
	if disp == piec.length {
		// If it's the same piece of the last insertion, we simply append.
		if piec == ed.piece {
			ed.piece.length++
			b.pieces[pidx].length++
			return nil
		}

		// Else we insert.
		ed.piece.length++
		b.pieces = slices.Insert(b.pieces, pidx+1, ed.piece)
		return nil
	}

	ed.piece.length++
	newPiece := ed.piece

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

func (b *Buffer) findPieceWithIdx(idx int) (i int, d int, err error) {
	disp := 0
	for i, piece := range b.pieces {
		ndisp := piece.length + disp
		if ndisp > idx {
			return i, idx - disp, nil
		}
		disp = ndisp
	}

	return 0, 0, errors.New("out of bounds")
}

func (b *Buffer) findPieceForInsertion(idx int) (i int, d int, err error) {
	disp := 0
	for i, piece := range b.pieces {
		ndisp := piece.length + disp
		if ndisp >= idx {
			return i, idx - disp, nil
		}
		disp = ndisp
	}

	return 0, 0, errors.New("out of bounds")
}

func (b *Buffer) pieceContent(p piece) []rune {
	arr := make([]rune, 0, p.length)
	buf := p.buffer
	bdisp := p.start
	for range p.length {
		if bdisp >= len(b.buffers[buf]) {
			bdisp = 0
			buf++
		}
		arr = append(arr, b.buffers[buf][bdisp])
		bdisp++
	}
	return arr
}

func (b *Buffer) indexByPiece(p piece, d int) (buffer int, bdisp int) {
	// If in the first (piece) buffer.
	if p.start+d < len(b.buffers[p.buffer]) {
		return p.buffer, d + p.start
	}

	disp := len(b.buffers[p.buffer]) - p.start
	buf := p.buffer + 1
	for {
		newdisp := disp + len(b.buffers[buf])
		if newdisp > d {
			return buf, newdisp - d
		}
		buf++
		disp = newdisp
	}
}

func (b *Buffer) Get(idx int) (rune, error) {
	piec, disp, err := b.findPieceWithIdx(idx)
	if err != nil {
		return ' ', err
	}

	buf, d := b.indexByPiece(b.pieces[piec], disp)
	return b.buffers[buf][d], nil
}

func (b *Buffer) Delete(idx int) error {
	pidx, disp, err := b.findPieceWithIdx(idx)
	if err != nil {
		return err
	}

	piec := b.pieces[pidx]

	buffer, _ := b.indexByPiece(b.pieces[pidx], disp)

	var ed *edit
	newEdit := edit{deletion: true, piece: piece{
		buffer: buffer,
		start:  0,
		length: 0,
	}}
	if len(b.edits) <= 0 {
		b.edits = append(b.edits, newEdit)
		ed = &b.edits[len(b.edits)-1]
	} else {
		ed = &b.edits[len(b.edits)-1]
		// TODO.
		if !ed.deletion {
			b.edits = append(b.edits, newEdit)
			ed = &b.edits[len(b.edits)-1]
		}
	}

	switch disp {
	// If removing from the top of the piece, we can simply decrease.
	case piec.length - 1:
		ed.piece.start--
		ed.piece.length++
		b.pieces[pidx].length--
	// If removing from the beggining of the piece, we can simply increase the
	// start.
	case 0:
		b.pieces[pidx].start++
		b.pieces[pidx].length--

		// If the piece begins at the end of the buffer.
		if b.pieces[pidx].start == len(b.buffers[b.pieces[pidx].buffer]) {
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
	}

	return nil
}

func (b *Buffer) Size() int {
	return b.size
}
