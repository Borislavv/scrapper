package vointerface

import "time"

type Timestamper interface {
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}
