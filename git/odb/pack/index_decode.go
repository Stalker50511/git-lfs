package pack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	MagicWidth   = 4
	VersionWidth = 4
	V1Width      = 0
	V2Width      = MagicWidth + VersionWidth

	FanoutEntries    = 256
	FanoutEntryWidth = 4
	FanoutWidth      = FanoutEntries * FanoutEntryWidth

	OffsetV1Start = V1Width + FanoutWidth
	OffsetV2Start = V2Width + FanoutWidth

	ObjectNameWidth        = 20
	ObjectCRCWidth         = 4
	ObjectSmallOffsetWidth = 4
	ObjectLargeOffsetWidth = 8

	ObjectEntryV1Width = ObjectNameWidth + ObjectSmallOffsetWidth
	ObjectEntryV2Width = ObjectNameWidth + ObjectCRCWidth + ObjectSmallOffsetWidth
)

var (
	ErrShortFanout = errors.New("git/odb/pack: too short fanout table")

	indexHeader = []byte{0xff, 0x74, 0x4f, 0x63}
)

func DecodeIndex(r io.ReaderAt) (*Index, error) {
	version, err := decodeIndexHeader(r)
	if err != nil {
		return nil, err
	}

	fanout, err := decodeIndexFanout(r, version.Width())
	if err != nil {
		return nil, err
	}

	return &Index{
		version: version,
		fanout:  fanout,

		f: r,
	}, nil
}

func decodeIndexHeader(r io.ReaderAt) (IndexVersion, error) {
	hdr := make([]byte, 4)
	if _, err := r.ReadAt(hdr, 0); err != nil {
		return VersionUnknown, err
	}

	if bytes.Equal(hdr, indexHeader) {
		vb := make([]byte, 4)
		if _, err := r.ReadAt(vb, 4); err != nil {
			return VersionUnknown, err
		}

		version := IndexVersion(binary.BigEndian.Uint32(vb))
		switch version {
		case V1, V2:
			return version, nil
		}

		return version, &UnsupportedVersionErr{uint32(version)}
	}
	return V1, nil
}

func decodeIndexFanout(r io.ReaderAt, offset int64) ([]uint32, error) {
	b := make([]byte, 256*4)
	if _, err := r.ReadAt(b, offset); err != nil {
		if err == io.EOF {
			return nil, ErrShortFanout
		}
		return nil, err
	}

	fanout := make([]uint32, 256)
	for i, _ := range fanout {
		fanout[i] = binary.BigEndian.Uint32(b[(i * 4):])
	}

	return fanout, nil
}
