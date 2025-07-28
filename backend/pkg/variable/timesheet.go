package variable

type TimeSheetListStatusEnum string

const (
	TimeSheetListStatusOpen      TimeSheetListStatusEnum = "open"
	TimeSheetListStatusClosed    TimeSheetListStatusEnum = "closed"
	TimeSheetListStatusReview    TimeSheetListStatusEnum = "review"
	TimeSheetListStatusFinalized TimeSheetListStatusEnum = "finalized"
	TimeSheetListStatusArchived  TimeSheetListStatusEnum = "archived"
)

type TimeSheetStatusEnum string

const (
	TimeSheetStatusDraft     TimeSheetStatusEnum = "draft"
	TimeSheetStatusSubmitted TimeSheetStatusEnum = "submitted"
	TimeSheetStatusApproved  TimeSheetStatusEnum = "approved"
	TimeSheetStatusRejected  TimeSheetStatusEnum = "rejected"
	TimeSheetStatusLocked    TimeSheetStatusEnum = "locked"
)

type LeaveTypeEnum string

const (
	LeaveTypeAnnual       LeaveTypeEnum = "annual_leave"
	LeaveTypePersonal     LeaveTypeEnum = "personal_leave"
	LeaveTypeBusinessTrip LeaveTypeEnum = "business_trip"
	LeaveTypeRemoteWork   LeaveTypeEnum = "remote_work"
	LeaveTypeUnpaid       LeaveTypeEnum = "unpaid_leave"
	LeaveTypeOther        LeaveTypeEnum = "other_leave"
)
