package variable

type (
	StatusWorkSchedule string
	RepeatTypeEnum     string
	WeekdayEnum        string
)

const (
	Sunday    WeekdayEnum = "sunday"
	Monday    WeekdayEnum = "monday"
	Tuesday   WeekdayEnum = "tuesday"
	Wednesday WeekdayEnum = "wednesday"
	Thursday  WeekdayEnum = "thursday"
	Friday    WeekdayEnum = "friday"
	Saturday  WeekdayEnum = "saturday"
)

const (
	Expired  StatusWorkSchedule = "expired"
	InActive StatusWorkSchedule = "inactive"
	Active   StatusWorkSchedule = "active"
)
const (
	Weekly  RepeatTypeEnum = "weekly"
	Monthly RepeatTypeEnum = "monthly"
)
