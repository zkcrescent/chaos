package config

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
	"path"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/juju/errors"
)

type File struct {
	Name string
	Size int64
	hash hash.Hash
}

func NewFile(name string) *File {
	return &File{Name: name, hash: md5.New()}
}

func (f *File) CheckSum() string {
	return fmt.Sprintf("%x", f.hash.Sum(nil))
}

func (f *File) Write(p []byte) (int, error) {
	n, err := f.hash.Write(p)
	if err == nil {
		f.Size += int64(n)
	}
	return n, nil
}

type FileStore interface {
	GetFile(filename string) (io.ReadCloser, *File, error)
	GetFileWithDir(filename, dir string) (io.ReadCloser, *File, error)
	SaveFile(filename string, content io.Reader) (*File, error)
	SaveWithDir(filename string, dir string, content io.Reader) (*File, error)
	Name() string
	Remove(filename string) error
}

type OSSStore struct {
	bucket *oss.Bucket
}

func NewOSSStore(bucket *oss.Bucket) *OSSStore {
	return &OSSStore{bucket}
}

func (fs *OSSStore) SignURL(filename string, secs int64) (string, error) {
	s, err := fs.bucket.SignURL(filename, oss.HTTPGet, secs)
	if err != nil {
		return "", errors.Annotatef(err, "sign url: %v", filename)
	}
	return s, nil
}

func (fs *OSSStore) SignURLWithDir(filename, dir string, secs int64) (string, error) {
	return fs.SignURL(path.Join(dir, filename), secs)
}

func (fs *OSSStore) GetFile(filename string) (io.ReadCloser, *File, error) {
	f := NewFile(filename)
	reader, err := fs.bucket.GetObject(filename)
	if err != nil {
		return nil, nil, errors.Annotatef(err, "get object: %v", filename)
	}
	return teeReadCloser(reader, f), f, nil
}

func (fs *OSSStore) GetFileWithDir(filename, dir string) (io.ReadCloser, *File, error) {
	return fs.GetFile(path.Join(dir, filename))
}

func (fs *OSSStore) SaveFile(filename string, content io.Reader) (*File, error) {
	f := NewFile(filename)
	if err := fs.bucket.PutObject(filename, io.TeeReader(content, f)); err != nil {
		return nil, errors.Annotatef(err, "put object: %v", filename)
	}
	return f, nil
}

func (s *OSSStore) SaveWithDir(filename string, dir string, content io.Reader) (*File, error) {
	return s.SaveFile(path.Join(dir, filename), content)
}

func (fs *OSSStore) Name() string {
	return "oss"
}

func (fs *OSSStore) Remove(filename string) error {
	return fs.bucket.DeleteObject(filename)
}

type readCloser struct {
	io.Reader
	io.Closer
}

func teeReadCloser(rc io.ReadCloser, w io.Writer) io.ReadCloser {
	tee := io.TeeReader(rc, w)
	return readCloser{tee, rc}
}
