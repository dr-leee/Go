package types

import (
	"time"
)

type Log struct {
	ID    int
	Login string
	Time  time.Time
}

/*
func (l Log) Error() string {
	return fmt.Sprintf("id log: %d", l.ID)
}
*/
