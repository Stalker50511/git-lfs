package pack

import "fmt"

type IndexVersion uint32

const (
	VersionUnknown IndexVersion = iota
)

func (v IndexVersion) Width() int64 {
	switch v {
	case V1:
		return V1Width
	case V2:
		return V2Width
	}

	panic(fmt.Sprintf("git/odb/pack: width unknown for pack version %d", v))
}

func (v IndexVersion) Search(idx *Index, name []byte, at int64) (*IndexEntry, int, error) {
	switch v {
	case V1:
		return v1Search(idx, name, at)
	case V2:
		return v2Search(idx, name, at)
	}
	return nil, 0, &UnsupportedVersionErr{Got: uint32(v)}
}
