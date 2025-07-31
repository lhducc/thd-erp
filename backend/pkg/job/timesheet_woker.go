package job

import "context"

func Worker(ctx context.Context, jobs <-chan TimesheetJob, results chan<- error) {
	for job := range jobs {
		select {
		case <-ctx.Done():
			results <- ctx.Err()
			return
		default:
			err := job.Service.CalculateForEmployee(job.Timesheet, job.TimesheetList)
			results <- err
		}
	}
}
