package logrotate

import (
	"io"
	"log"
	"time"

	"github.com/zkcrescent/chaos/utils"
	ec "github.com/zkcrescent/chaos/utils/exit-chan"
)

const (
	DefaultEmitTimeout = time.Millisecond * 10
)

const (
	ErrorLogOverflow   = utils.Error("log overflow, discard.")
	ErrorUnexpectError = utils.Error("unexpect error.")
)

type WriterInterceptor func(*Writer, []byte) error

type Writer struct {
	dst          io.WriteCloser
	interceptors []WriterInterceptor
}

func (w *Writer) intercept(b []byte) error {
	if w.interceptors != nil {
		for _, i := range w.interceptors {
			if i == nil {
				continue
			}
			if err := i(w, b); err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Writer) Close() error {
	return w.dst.Close()
}

func (w *Writer) Write(b []byte) (int, error) {
	if err := w.intercept(b); err != nil {
		return 0, err
	}

	return w.dst.Write(b)
}

type BufferWriter struct {
	Writer

	lc          chan []byte
	stopc       *ec.ExitChan
	emitTimeout time.Duration
}

func NewBufferWriter(w io.WriteCloser, buffer int, ins ...WriterInterceptor) *BufferWriter {
	return &BufferWriter{
		Writer: Writer{
			dst:          w,
			interceptors: ins,
		},
		lc:          make(chan []byte, buffer),
		stopc:       ec.New(),
		emitTimeout: DefaultEmitTimeout,
	}
}

func (w *BufferWriter) SetEmitTimeout(d time.Duration) *BufferWriter {
	w.emitTimeout = d
	return w
}

func (w *BufferWriter) Close() error {
	if w.stopc.Exited() {
		return nil
	}

	w.stopc.Close()
	return nil
}

func (w *BufferWriter) Write(b []byte) (int, error) {
	if err := w.intercept(b); err != nil {
		return 0, err
	}

	select {
	case <-w.stopc.Chan():
		return 0, io.EOF
	case <-time.After(w.emitTimeout):
		return 0, ErrorLogOverflow
	case w.lc <- b:
		return len(b), nil
	}

	return 0, ErrorUnexpectError
}

func (w *BufferWriter) Loop() {
	for {
		select {
		case <-w.stopc.Chan():
			return
		case l := <-w.lc:
			if n, err := w.dst.Write(l); err != nil {
				log.Println("BufferWriter.Write failed:", err)
			} else if ll := len(l); n != ll {
				log.Printf("BufferWriter.Write partial: %v of %v", n, ll)
			}
		}
	}
}
