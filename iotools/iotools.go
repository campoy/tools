package iotools

import "io"

type readerAt struct {
	r   io.Reader
	buf []byte
}

func NewReaderAt(r io.Reader, bufSize int64) io.ReaderAt {
	return &readerAt{r, nil}
}

func (r *readerAt) ReadAt(p []byte, off int64) (n int, err error) {

	if off == 0 {
		return r.r.Read(p)
	}
	return 0, nil
}
