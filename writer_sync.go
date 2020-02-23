package nptrace

import (
	"io"
	"sync"
)

type WriterSync interface {
	io.Writer
	sync.Locker
}

type writerSync struct {
	wr io.Writer
	sync.Mutex
}

func NewWriterSync(wr io.Writer) WriterSync {
	return &writerSync{
		wr:    wr,
		Mutex: sync.Mutex{},
	}
}

func (w writerSync) Write(p []byte) (int, error) {
	w.Lock()
	n, err := w.wr.Write(p)
	w.Unlock()
	return n, err
}
