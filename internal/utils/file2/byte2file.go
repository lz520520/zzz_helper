package file2

import (
	"io"
	"io/fs"
	"reflect"
	"time"
)

// A file is a single file in the FS.
// It implements fs.FileInfo and fs.DirEntry.
type bytesFile struct {
	name string
	data []byte
	hash [16]byte // truncated SHA256 hash
}

func (f *bytesFile) Name() string               { return f.name }
func (f *bytesFile) Size() int64                { return int64(len(f.data)) }
func (f *bytesFile) ModTime() time.Time         { return time.Time{} }
func (f *bytesFile) IsDir() bool                { return false }
func (f *bytesFile) Sys() any                   { return nil }
func (f *bytesFile) Type() fs.FileMode          { return f.Mode().Type() }
func (f *bytesFile) Info() (fs.FileInfo, error) { return f, nil }

func (f *bytesFile) Mode() fs.FileMode {
	if f.IsDir() {
		return fs.ModeDir | 0555
	}
	return 0444
}

// An OpenFile is a regular file open for reading.
type OpenFile struct {
	file   *bytesFile // the file itself
	offset int64      // current read offset
}

func (f *OpenFile) Close() error               { return nil }
func (f *OpenFile) Stat() (fs.FileInfo, error) { return f.file, nil }

// Read reads up to len(b) bytes from the File and stores them in b.
// It returns the number of bytes read and any error encountered.
// At end of file, Read returns 0, io.EOF.
func (f *OpenFile) Read(b []byte) (int, error) {
	if f.offset >= int64(len(f.file.data)) {
		return 0, io.EOF
	}
	if f.offset < 0 {
		return 0, &fs.PathError{Op: "read", Path: f.file.name, Err: fs.ErrInvalid}
	}
	n := copy(b, f.file.data[f.offset:])
	f.offset += int64(n)
	return n, nil
}

var (
//_ io.Seeker   = (*OpenFile)(nil)
//_ fs.FileInfo = (*bytesFile)(nil)
)

func (f *OpenFile) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case 0:
		// offset += 0
	case 1:
		offset += f.offset
	case 2:
		offset += int64(len(f.file.data))
	}
	if offset < 0 || offset > int64(len(f.file.data)) {
		return 0, &fs.PathError{Op: "seek", Path: f.file.name, Err: fs.ErrInvalid}
	}
	f.offset = offset
	return offset, nil
}

func (f *OpenFile) Write(b []byte) (int, error) {
	if reflect.ValueOf(f).IsNil() {
		return 0, &fs.PathError{Op: "write", Path: f.file.name, Err: fs.ErrInvalid}
	}
	f.file.data = append(f.file.data, b...)
	return len(b), nil
}

func OpenBytesFile(b []byte) *OpenFile {
	return &OpenFile{
		file: &bytesFile{
			data: b,
		},
		offset: 0,
	}
}
