package funcs

import (
	"time"
)

//need filtered period
//include start and end
type ExcluedPeriod struct {
	Start time.Time
	End   time.Time
	Scope string
}

type ExcluedPeriods struct {
	periods []ExcluedPeriod
}

