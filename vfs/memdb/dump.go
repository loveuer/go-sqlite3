package memdb

import (
	"errors"
	"io"
)

// MemDB is a no-export memDB struct dump helper.
// Note: This struct does not provide concurrency safety,
// and you must manage concurrent access yourself.
type MemDB struct {
	bs     []byte
	offset int
}

func (m *MemDB) Close() error {
	return nil
}

func (m *MemDB) Seek(offset int64, whence int) (int64, error) {
	var newOffset int

	switch whence {
	case io.SeekStart:
		newOffset = int(offset)
	case io.SeekCurrent:
		newOffset = m.offset + int(offset)
	case io.SeekEnd:
		newOffset = len(m.bs) + int(offset)
	default:
		return 0, errors.New("invalid whence")
	}

	if newOffset < 0 || newOffset > len(m.bs) {
		return 0, io.EOF // Out of bounds
	}

	m.offset = newOffset
	return int64(m.offset), nil
}

func (m *MemDB) Read(p []byte) (n int, err error) {
	if m.offset >= len(m.bs) {
		return 0, io.EOF // No more data to read
	}

	n = copy(p, m.bs[m.offset:]) // Copy data to p
	m.offset += n                // Update the offset

	return n, nil // Return number of bytes read and no error
}

func (m *memDB) Dump() io.ReadSeekCloser {

	m.lockMtx.Lock()
	defer m.lockMtx.Unlock()

	size := m.size
	d := &MemDB{bs: make([]byte, 0, size)}

	for _, bs := range m.data {
		d.bs = append(d.bs, bs[:]...)
	}

	return d
}
