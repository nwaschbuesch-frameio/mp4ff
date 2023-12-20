package mp4

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/nwaschbuesch-frameio/mp4ff/bits"
)

func TestStppEncDec(t *testing.T) {
	testCases := []struct {
		namespace      string
		schemaLocation string
		mimeTypes      string
		hasBtrt        bool
	}{
		{
			namespace:      "NS",
			schemaLocation: "location",
			mimeTypes:      "image/png image/jpg",
			hasBtrt:        false,
		},
		{
			namespace:      "NS",
			schemaLocation: "location",
			mimeTypes:      "image/png image/jpg",
			hasBtrt:        true,
		},
		{
			namespace:      "NS",
			schemaLocation: "",
			mimeTypes:      "",
			hasBtrt:        false,
		},
		{
			namespace:      "NS",
			schemaLocation: "",
			mimeTypes:      "",
			hasBtrt:        true,
		},
	}
	for _, tc := range testCases {
		stpp := NewStppBox(tc.namespace, tc.schemaLocation, tc.mimeTypes)
		if tc.hasBtrt {
			btrt := &BtrtBox{}
			stpp.AddChild(btrt)
			if stpp.Btrt != btrt {
				t.Error("Btrt link is broken")
			}
		}
		boxDiffAfterEncodeAndDecode(t, stpp)
	}
}

func TestStppWithEmtptyLists(t *testing.T) {
	const hexData = ("00000040737470700000000000000001" +
		"687474703a2f2f7777772e77332e6f72" +
		"672f6e732f74746d6c00000000000014" +
		"62747274000003ce00003b5800000430")
	data, err := hex.DecodeString(hexData)
	if err != nil {
		t.Error(err)
	}
	sr := bits.NewFixedSliceReader(data)
	box, err := DecodeBoxSR(0, sr)
	if err != nil {
		t.Error(err)
	}
	stpp, ok := box.(*StppBox)
	if !ok {
		t.Error("not an stpp box")
	}
	if int(stpp.Size()) != len(data) {
		t.Errorf("stpp size %d not same as %d", stpp.Size(), len(data))
	}
	buf := bytes.Buffer{}
	err = stpp.Encode(&buf)
	if err != nil {
		t.Error(err)
	}
	outData := buf.Bytes()
	if !bytes.Equal(data, outData) {
		t.Error("written stpp box differs from read")
	}

}
