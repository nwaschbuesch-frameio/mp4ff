package mp4

import (
	"io"

	"github.com/nwaschbuesch-frameio/mp4ff/bits"
)

// MfroBox - Movie Fragment Random Access Offset Box (mfro)
// Contained in : MfraBox (mfra)
type MfroBox struct {
	Version    byte
	Flags      uint32
	ParentSize uint32
}

// DecodeMfro - box-specific decode
func DecodeMfro(hdr BoxHeader, startPos uint64, r io.Reader) (Box, error) {
	data, err := readBoxBody(r, hdr)
	if err != nil {
		return nil, err
	}
	sr := bits.NewFixedSliceReader(data)
	return DecodeMfroSR(hdr, startPos, sr)
}

// DecodeMfroSR - box-specific decode
func DecodeMfroSR(hdr BoxHeader, startPos uint64, sr bits.SliceReader) (Box, error) {
	versionAndFlags := sr.ReadUint32()

	b := &MfroBox{
		Version:    byte(versionAndFlags >> 24),
		Flags:      versionAndFlags & flagsMask,
		ParentSize: sr.ReadUint32(),
	}
	return b, sr.AccError()
}

// Type - return box type
func (b *MfroBox) Type() string {
	return "mfro"
}

// Size - return calculated size
func (b *MfroBox) Size() uint64 {
	return uint64(boxHeaderSize + 4 + 4)
}

// Encode - write box to w
func (b *MfroBox) Encode(w io.Writer) error {
	sw := bits.NewFixedSliceWriter(int(b.Size()))
	err := b.EncodeSW(sw)
	if err != nil {
		return err
	}
	_, err = w.Write(sw.Bytes())
	return err
}

// EncodeSW - box-specific encode to slicewriter
func (b *MfroBox) EncodeSW(sw bits.SliceWriter) error {
	err := EncodeHeaderSW(b, sw)
	if err != nil {
		return err
	}
	versionAndFlags := (uint32(b.Version) << 24) + b.Flags
	sw.WriteUint32(versionAndFlags)
	sw.WriteUint32(b.ParentSize)
	return sw.AccError()
}

// Info - write box-specific information
func (b *MfroBox) Info(w io.Writer, specificBoxLevels, indent, indentStep string) error {
	bd := newInfoDumper(w, indent, b, int(b.Version), b.Flags)
	bd.write(" - parentSize: %d", b.ParentSize)
	return bd.err
}
