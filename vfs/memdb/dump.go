package memdb

import (
	"errors"
	"io"
)

// MemDB is a no-export memDB struct dump helper.
// Note: This struct does not provide concurrency safety,
// and you must manage concurrent access yourself.
type MemDB struct {
	db *memDB
}

type MemDBDumper struct {
	bs     []byte
	offset int
}

func (m *MemDBDumper) Close() error {
	return nil
}

func (m *MemDBDumper) Seek(offset int64, whence int) (int64, error) {
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

func (m *MemDBDumper) Read(p []byte) (n int, err error) {
	if m.offset >= len(m.bs) {
		return 0, io.EOF // No more data to read
	}

	n = copy(p, m.bs[m.offset:]) // Copy data to p
	m.offset += n                // Update the offset

	return n, nil // Return number of bytes read and no error
}

func (m *MemDB) Dump() io.ReadSeekCloser {

	m.db.lockMtx.Lock()
	defer m.db.lockMtx.Unlock()

	size := m.db.size
	d := &MemDBDumper{bs: make([]byte, 0, size)}

	for _, bs := range m.db.data {
		d.bs = append(d.bs, bs[:]...)
	}

	return d
}
