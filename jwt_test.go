package dist

import (
	"testing"
)

func TestTokener(t *testing.T) {
	tokener := NewJwtTokener("1885df74d00dbbe19274c6d955feeb5b")
	ts, err := tokener.Token(map[string]interface{}{
		"username": "test",
		"password": "abc123",
	})
	if nil != err {
		t.Error(err)
		return
	}
	t.Log(ts)
	r, err := tokener.Parse(ts)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(r["username"].(string))
	t.Log(r["password"])
}
