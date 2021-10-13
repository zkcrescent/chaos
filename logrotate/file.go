package logrotate

import (
	"io"
	"os"

	"github.com/zkcrescent/chaos/utils"
)

type WriterInit func(*File, io.WriteCloser) (io.WriteCloser, error)

type File struct {
	fileName  *fileName
	rotate    *Rotate
	newWriter WriterInit

	writer io.WriteCloser
}

func FileWriter(path string, opts ...Option) (*File, error) {
	fn := toFileName(path)
	if fn == nil {
		return nil, utils.Error("path is empty")
	}

	f := &File{
		fileName: fn,
	}

	for _, o := range opts {
		if o == nil {
			continue
		}
		if err := o(f); err != nil {
			return nil, err
		}
	}

	dst, err := f.Attach()
	if err != nil {
		return nil, err
	}
	if f.rotate != nil {
		dsts, err := dst.Stat()
		if err != nil {
			return nil, err
		}
		f.rotate.init(dsts)
	}
	if f.newWriter == nil {
		f.newWriter = DefaultNewWriter
	}

	f.writer, err = f.newWriter(f, dst)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *File) Attach() (*os.File, error) {
	return os.OpenFile(f.fileName.String(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
}

func (f *File) Flush() (*os.File, error) {
	return os.OpenFile(f.fileName.String(), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
}

func (f *File) Write(b []byte) (int, error) {
	nb := make([]byte, len(b))
	copy(nb, b)
	return f.writer.Write(nb)
}
