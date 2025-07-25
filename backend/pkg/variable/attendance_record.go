package variable

type StatusEnum string

const (
	Pending  StatusEnum = "pending"
	Approved StatusEnum = "approved"
	Rejected StatusEnum = "rejected"
)
