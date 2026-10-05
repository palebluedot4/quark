package pool

import (
	"compress/gzip"
	"io"
	"sync"
)

var gzipWriters = sync.Pool{
	New: func() any {
		return gzip.NewWriter(io.Discard)
	},
}

func GetGzipWriter(w io.Writer) *gzip.Writer {
	writer := gzipWriters.Get().(*gzip.Writer)
	writer.Reset(w)
	return writer
}

func PutGzipWriter(w *gzip.Writer) error {
	err := w.Close()
	w.Reset(io.Discard)
	gzipWriters.Put(w)
	return err
}
