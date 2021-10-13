package logrotate

import (
	"io"
	"os"
)

type Rotate struct {
	Total int
	Size  int
	state *rotateState
}

type rotateState struct {
	cur   int
	wrote int
}

func (r *Rotate) init(fi os.FileInfo) {
	r.state = &rotateState{
		wrote: int(fi.Size()),
	}
}

func (r *Rotate) RotateInterceptor(f *File) WriterInterceptor {
	return func(w *Writer, b []byte) error {
		r.state.wrote += len(b)
		if r.state.wrote >= r.Size {
			if err := r.Rotate(f, w); err != nil {
				return err
			}
			r.state.wrote = len(b)
		}
		return nil
	}
}

func (r *Rotate) Rotate(f *File, w *Writer) error {
	var (
		located bool
		ptr     string
		fl      = map[string]os.FileInfo{}
	)
	for !located {
		r.state.cur = r.Next()

		var (
			cur  os.FileInfo
			next os.FileInfo
			err  error
			ok   bool
		)

		cp := f.fileName.StringInNumber(r.state.cur)
		if cur, ok = fl[cp]; !ok {
			cur, err = os.Stat(cp)
			if os.IsNotExist(err) {
				ptr = cp
				located = true
				break
			} else if err != nil {
				return err
			}
			fl[cp] = cur
		}

		np := f.fileName.StringInNumber(r.Next())
		if next, ok = fl[np]; !ok {
			next, err = os.Stat(np)
			if os.IsNotExist(err) {
				ptr = np
				located = true
				break
			} else if err != nil {
				return err
			}
			fl[np] = next
		}

		if next.ModTime().Before(cur.ModTime()) {
			ptr = np
			located = true
			break
		}
	}

	if err := os.Rename(f.fileName.String(), ptr); err != nil {
		return err
	}

	nf, err := f.Flush()
	if err != nil {
		return err
	}
	w.dst.Close()
	w.dst = nf
	return nil
}

func (r *Rotate) Next() int {
	cur := r.state.cur + 1
	if cur > r.Total {
		cur %= r.Total
	}
	return cur
}

func (r *Rotate) NewWriter(f *File, w io.WriteCloser) (io.WriteCloser, error) {
	return &Writer{
		dst: w,
		interceptors: []WriterInterceptor{
			r.RotateInterceptor(f),
		},
	}, nil
}
