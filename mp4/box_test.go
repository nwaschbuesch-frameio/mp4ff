package mp4

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/nwaschbuesch-frameio/mp4ff/bits"
)

// TestBadBoxAndRemoveBoxDecoder checks that we can avoid decoder error by removing a BoxDecode.
//
// The box is then interpreted as an UnknownBox and its data is not further processed with decoded.
func TestBadBoxAndRemoveBoxDecoder(t *testing.T) {
	badMetaBox := (`000000416d6574610000002168646c7300000000000000006d64746100000000` +
		`000000000000000000000000106b657973000000000000000000000008696c7374`)
	data, err := hex.DecodeString(badMetaBox)
	if err != nil {
		t.Error(err)
	}
	sr := bits.NewFixedSliceReader(data)
	_, err = DecodeBoxSR(0, sr)
	if err == nil {
		t.Errorf("reading bad meta box should have failed")
	}
	sr = bits.NewFixedSliceReader(data)
	RemoveBoxDecoder("meta")
	defer SetBoxDecoder("meta", DecodeMeta, DecodeMetaSR)
	box, err := DecodeBoxSR(0, sr)
	if err != nil {
		t.Error(err)
	}
	_, ok := box.(*MetaBox)
	if ok {
		t.Errorf("box should not be MetaBox")
	}
	unknown, ok := box.(*UnknownBox)
	if !ok {
		t.Errorf("box should be unknown")
	}
	if unknown.Type() != "meta" {
		t.Errorf("unknown type %q instead of meta", unknown.Type())
	}
	b := bytes.Buffer{}
	err = unknown.Encode(&b)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(data, b.Bytes()) {
		t.Errorf("written unknown differs")
	}
}
