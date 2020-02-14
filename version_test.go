package dist

import "testing"

func TestVersion(t *testing.T){
	v := NewVersion(0,0, 1)
	t.Log(v.Major(), v.Minor(), v.Revision())
	t.Log(v)
	if v.Major() != 0 || v.Minor() != 0 || v.Revision() != 1 ||
		v.String() != "0.0.1" {
		t.Error("version error")
	}
}