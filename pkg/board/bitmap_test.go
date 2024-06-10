package board

import (
	"testing"
)

func TestBitmaps(t *testing.T) {
	mask := getMaskForPosition(2, 3)
	t.Logf("Mask: %d\n", mask)

	var m uint64
	setBitByMask(&m, mask)
	if !isBitSetByPos(m, NewPosition(2, 3)) {
		t.Errorf("Bit not set properly")
	}

}
