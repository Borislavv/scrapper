package jobrunnerinterface

import "context"

type JobRunner interface {
	Run(ctx context.Context)
}
