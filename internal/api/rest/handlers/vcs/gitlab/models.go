package gitlab

import "time"

type RequestParams struct {
	ProjectID int `validate:"required"`
	Since     time.Time
	Until     time.Time
	Status    string `validate:"omitempty,oneof=all open closed merged"`
}
