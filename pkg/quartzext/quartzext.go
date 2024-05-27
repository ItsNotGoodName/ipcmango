package quartzext

import (
	"context"

	"github.com/reugn/go-quartz/quartz"
)

func NewServiceScheduler(s quartz.Scheduler) ServiceScheduler {
	return ServiceScheduler{
		Scheduler: s,
	}
}

type ServiceScheduler struct {
	quartz.Scheduler
}

func (s ServiceScheduler) String() string {
	return "quartzext.ServiceScheduler"
}

func (s ServiceScheduler) Serve(ctx context.Context) error {
	s.Start(ctx)
	s.Wait(context.Background())
	return nil
}

func NewJobDetail(job quartz.Job) *quartz.JobDetail {
	return quartz.NewJobDetail(job, quartz.NewJobKey(job.Description()))
}
