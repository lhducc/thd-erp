package struct_support

import "erp/backend/pkg/variable"

var (
	ValidRepeatTypes = map[variable.RepeatTypeEnum]struct{}{
		variable.Weekly:  {},
		variable.Monthly: {},
	}

	ValidWeekdays = map[variable.WeekdayEnum]struct{}{
		variable.Monday:    {},
		variable.Tuesday:   {},
		variable.Wednesday: {},
		variable.Thursday:  {},
		variable.Friday:    {},
		variable.Saturday:  {},
		variable.Sunday:    {},
	}

	ValidStatus = map[variable.StatusWorkSchedule]struct{}{
		variable.Active:   {},
		variable.InActive: {},
		variable.Expired:  {},
	}
)
