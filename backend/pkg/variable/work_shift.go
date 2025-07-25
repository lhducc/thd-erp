package variable

type WorkDayEnum string
type TimeOfDayEnum string

const (
	FullDay WorkDayEnum = "1"
	HaftDay WorkDayEnum = "0.5"
	NoWork  WorkDayEnum = "0"
)
