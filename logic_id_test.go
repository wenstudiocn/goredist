package dist

import (
"testing"
)

func TestLogicID(t *testing.T) {
	var logicId uint64 = 1010203
	s, err := NewLogicIDByID(logicId)
	if err != nil {
		t.Error("failed:", err)
		return
	}
	if s.CataId() != 1 {
		t.Errorf("cataId=%d", s.cataId)
	}
	t.Logf("%v.%v.%v.%v", s.funcId, s.cataId, s.subId, s.instId)

	s, err = NewLogicIDByID(1000010203)
	if err != nil {
		t.Log("error:", err)
	}
	if s != nil {
		t.Error("error2")
	}
}
