package logrotate

import (
	"io"
	"time"
)

type Option func(f *File) error

func DefaultNewWriter(f *File, w io.WriteCloser) (io.WriteCloser, error) {
	return w, nil
}

func Rotating(limitSize, backups int) Option {
	return func(f *File) error {
		f.rotate = &Rotate{
			Size:  limitSize,
			Total: backups,
		}
		if f.newWriter == nil {
			f.newWriter = f.rotate.NewWriter
		}
		return nil
	}
}

func Buffer(logSize int, emitTimeouts ...time.Duration) Option {
	return func(f *File) error {
		f.newWriter = func(f *File, w io.WriteCloser) (io.WriteCloser, error) {
			ins := []WriterInterceptor{}
			if f.rotate != nil {
				ins = append(ins, f.rotate.RotateInterceptor(f))
			}
			ww := NewBufferWriter(w, logSize, ins...)
			if emitTimeouts != nil && len(emitTimeouts) > 0 {
				ww.SetEmitTimeout(emitTimeouts[0])
			}
			go ww.Loop()
			return ww, nil
		}
		return nil
	}
}
