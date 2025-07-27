package main

type Buffer interface {
	Insert(idx int, r rune) error
	Delete(idx int) error
	Get(idx int) (rune, error)
	Size() int
}

type BufferReader struct {
	buffer Buffer
	disp   int
}

func (r *BufferReader) Read(out []rune) (int, error) {
	i := 0
	for i = range out {
		b, err := r.buffer.Get(r.disp)
		if err != nil {
			return i, err
		}
		out[i] = b
		r.disp++
	}
	return i, nil
}

func DisplacedReader(buffer Buffer, disp int) *BufferReader {
	return &BufferReader{
		buffer: buffer,
		disp:   disp,
	}
}
