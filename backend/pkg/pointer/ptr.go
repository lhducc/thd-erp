package pointer

func Float64(v float64) *float64 {
	return &v
}

func Int(v int) *int {
	return &v
}

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}
