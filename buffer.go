package main

type Buffer interface {
	Insert(disp int, b byte) error
	Get(disp int) (byte, error)
}

type BufferReader struct {
	buffer Buffer
	disp   int
}

func (r *BufferReader) Read(out []byte) (int, error) {
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
