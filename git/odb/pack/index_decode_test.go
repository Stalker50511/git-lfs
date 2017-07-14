package pack

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeIndexV1(t *testing.T) {
	idx, err := DecodeIndex(bytes.NewReader(make([]byte, FanoutWidth)))

	assert.NoError(t, err)
	assert.Equal(t, V1, idx.version)
	assert.EqualValues(t, 0, idx.Count())
}

func TestDecodeIndexV1InvalidFanout(t *testing.T) {
	idx, err := DecodeIndex(bytes.NewReader(make([]byte, FanoutWidth-1)))

	assert.Equal(t, ErrShortFanout, err)
	assert.Nil(t, idx)
}

func TestDecodeIndexV2(t *testing.T) {
	buf := make([]byte, 0, V2Width+FanoutWidth)
	buf = append(buf, 0xff, 0x74, 0x4f, 0x63)
	buf = append(buf, 0x0, 0x0, 0x0, 0x2)
	for i := 0; i < FanoutEntries; i++ {
		x := make([]byte, 4)

		binary.BigEndian.PutUint32(x, uint32(3))

		buf = append(buf, x...)
	}

	idx, err := DecodeIndex(bytes.NewReader(buf))

	assert.NoError(t, err)
	assert.Equal(t, V2, idx.version)
	assert.EqualValues(t, 3, idx.Count())
}

func TestDecodeIndexV2InvalidFanout(t *testing.T) {
	buf := make([]byte, 0, V2Width+FanoutWidth-FanoutEntryWidth)
	buf = append(buf, 0xff, 0x74, 0x4f, 0x63)
	buf = append(buf, 0x0, 0x0, 0x0, 0x2)
	buf = append(buf, make([]byte, FanoutWidth-1)...)

	idx, err := DecodeIndex(bytes.NewReader(buf))

	assert.Equal(t, ErrShortFanout, err)
	assert.Nil(t, idx)
}

func TestDecodeIndexUnsupportedVersion(t *testing.T) {
	buf := make([]byte, 0, V2Width)
	buf = append(buf, 0xff, 0x74, 0x4f, 0x63)
	buf = append(buf, 0x0, 0x0, 0x0, 0x3)

	idx, err := DecodeIndex(bytes.NewReader(buf))

	assert.EqualError(t, err, "git/odb/pack: unsupported version: 3")
	assert.Nil(t, idx)
}

func TestDecodeIndexEmptyContents(t *testing.T) {
	idx, err := DecodeIndex(bytes.NewReader(make([]byte, 0)))

	assert.Equal(t, io.EOF, err)
	assert.Nil(t, idx)
}
