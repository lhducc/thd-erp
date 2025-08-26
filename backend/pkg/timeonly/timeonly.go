package timeonly

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TimeOnly represents a clock time without date. JSON format: "HH:MM:SS"
type TimeOnly struct {
	time.Time
}

const timeLayout = "15:04:05"

func (t *TimeOnly) UnmarshalJSON(b []byte) error {
	// accept null
	s := strings.Trim(string(b), `"`)
	if s == "null" || s == "" {
		t.Time = time.Time{}
		return nil
	}
	parsed, err := time.Parse(timeLayout, s)
	if err != nil {
		return fmt.Errorf("invalid time format: %w", err)
	}
	// Use parsed time on zero date (0000-01-01) but keep the Time value
	t.Time = parsed
	return nil
}

func (t TimeOnly) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return json.Marshal(t.Format(timeLayout))
}

// Scan implements sql.Scanner
func (t *TimeOnly) Scan(src interface{}) error {
	if src == nil {
		t.Time = time.Time{}
		return nil
	}
	switch v := src.(type) {
	case time.Time:
		// keep only the clock part
		tt := time.Date(1, 1, 1, v.Hour(), v.Minute(), v.Second(), v.Nanosecond(), time.UTC)
		t.Time = tt
		return nil
	case []byte:
		s := string(v)
		parsed, err := time.Parse(timeLayout, s)
		if err != nil {
			// try parsing full timestamp
			parsed2, err2 := time.Parse(time.RFC3339, s)
			if err2 != nil {
				return err
			}
			t.Time = parsed2
			return nil
		}
		t.Time = parsed
		return nil
	case string:
		parsed, err := time.Parse(timeLayout, v)
		if err != nil {
			parsed2, err2 := time.Parse(time.RFC3339, v)
			if err2 != nil {
				return err
			}
			t.Time = parsed2
			return nil
		}
		t.Time = parsed
		return nil
	default:
		return fmt.Errorf("cannot scan %T into TimeOnly", src)
	}
}

// Value implements driver.Valuer
func (t TimeOnly) Value() (driver.Value, error) {
	if t.Time.IsZero() {
		return nil, nil
	}
	return t.Format(timeLayout), nil
}

func (t TimeOnly) String() string {
	if t.Time.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}
