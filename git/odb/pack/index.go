package pack

import (
	"io"
)

type Index struct {
	version IndexVersion
	fanout  []uint32

	f io.ReaderAt
}

func (i *Index) Count() int {
	return int(i.fanout[255])
}

func (i *Index) Entry(name []byte) (*IndexEntry, error) {
	left, right := i.bounds(name)

	for left < right {
		mid := (left + right) / 2

		entry, cmp, err := i.version.Search(i, name, mid)
		if err != nil {
			return nil, err
		}

		if cmp == 0 {
			return entry, nil
		} else if cmp < 0 {
			right = mid
		} else if cmp > 0 {
			left = mid
		}
	}
	return nil, nil
}

func (i *Index) readAt(p []byte, at int64) (n int, err error) {
	return i.f.ReadAt(p, at)
}

func (i *Index) bounds(name []byte) (left, right int64) {
	if name[0] == 0 {
		left = 0
	} else {
		left = int64(i.fanout[name[0]-1])
	}

	if name[0] == 255 {
		right = int64(i.Count())
	} else {
		right = int64(i.fanout[name[0]+1])
	}

	return
}
