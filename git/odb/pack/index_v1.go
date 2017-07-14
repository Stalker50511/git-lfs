package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
)

const (
	V1 IndexVersion = 1
)

var (
	ErrIndexOutOfBounds = errors.New("git/odb/pack: index is out of bounds")
)

func v1Search(idx *Index, name []byte, at int64) (*IndexEntry, int, error) {
	if at > int64(idx.Count()) {
		return nil, 0, ErrIndexOutOfBounds
	}

	var sha [20]byte

	if _, err := idx.readAt(sha[:], v1ShaOffset(at)); err != nil {
		return nil, 0, err
	}

	cmp := bytes.Compare(name, sha[:])
	if cmp != 0 {
		return nil, cmp, nil
	}

	var offs [4]byte
	if _, err := idx.readAt(offs[:], v1EntryOffset(at)); err != nil {
		return nil, 0, err
	}

	return &IndexEntry{
		PackOffset: uint64(binary.BigEndian.Uint32(offs[:])),
	}, 0, nil
}

func v1ShaOffset(at int64) int64 {
	// Skip forward until the desired entry.
	return v1EntryOffset(at) +
		// Skip past the 4-byte object offset in the desired entry to
		// the SHA1.
		ObjectSmallOffsetWidth
}

func v1EntryOffset(at int64) int64 {
	// Skip the L1 fanout table
	return OffsetV1Start +
		// Skip the object entries before the one located at "at"
		(ObjectEntryV1Width * at)
}
