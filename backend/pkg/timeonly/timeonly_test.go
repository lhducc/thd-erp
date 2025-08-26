package timeonly

import (
	"encoding/json"
	"testing"
)

func TestMarshalUnmarshalJSON(t *testing.T) {
	var t1 TimeOnly
	if err := json.Unmarshal([]byte(`"03:30:00"`), &t1); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if t1.String() != "03:30:00" {
		t.Fatalf("expected 03:30:00 got %s", t1.String())
	}

	b, err := json.Marshal(t1)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		t.Fatalf("unmarshal string failed: %v", err)
	}
	if s != "03:30:00" {
		t.Fatalf("expected json string 03:30:00 got %s", s)
	}
}
