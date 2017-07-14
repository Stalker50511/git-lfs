package pack

import (
	"bytes"
	"encoding/binary"
)

const (
	V2 IndexVersion = 2
)

func v2Search(idx *Index, name []byte, at int64) (*IndexEntry, int, error) {
	var sha [20]byte
	if _, err := idx.readAt(sha[:], v2ShaOffset(at)); err != nil {
		return nil, 0, err
	}

	cmp := bytes.Compare(name, sha[:])
	if cmp != 0 {
		return nil, cmp, nil
	}

	var offs [4]byte
	if _, err := idx.readAt(offs[:], v2SmallOffsetOffset(at, int64(idx.Count()))); err != nil {
		return nil, 0, err
	}

	loc := uint64(binary.BigEndian.Uint32(offs[:]))
	if loc&0x80000000 > 0 {
		var offs [8]byte
		if _, err := idx.readAt(offs[:], int64(loc&0x7fffffff)); err != nil {
			return nil, 0, err
		}

		loc = binary.BigEndian.Uint64(offs[:])
	}
	return &IndexEntry{PackOffset: loc}, 0, nil
}

func v2ShaOffset(at int64) int64 {
	// Skip the packfile index header and the L1 fanout table.
	return OffsetV2Start +
		// Skip until the desired name in the sorted names table.
		(ObjectNameWidth * at)
}

func v2SmallOffsetOffset(at, total int64) int64 {
	// Skip the packfile index header and the L1 fanout table.
	return OffsetV2Start +
		// Skip the name table.
		(ObjectNameWidth * total) +
		// Skip the CRC table.
		(ObjectCRCWidth * total) +
		// Skip until the desired index in the small offsets table.
		(ObjectSmallOffsetWidth * at)
}
