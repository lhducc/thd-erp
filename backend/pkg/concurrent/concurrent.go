package concurrent

type ConcurrencyConfig struct {
	TimesheetWorkers int
	UploadWorkers    int
	MaxRetries       int
}

var AppConcurrencyConfig = ConcurrencyConfig{
	TimesheetWorkers: 10,
	UploadWorkers:    8,
	MaxRetries:       3,
}
