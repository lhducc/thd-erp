package variable

type WorkDayEnum float64
type TimeOfDayEnum string

const (
	FullDay WorkDayEnum = 1
	HaftDay WorkDayEnum = 0.5
	NoWork  WorkDayEnum = 0
)
