package actor

import (
	"fmt"
)

type PID struct {
	area string
	id   string
}

func NewPID(area, id string) *PID {
	return &PID{
		area,
		id,
	}
}

func (this *PID) IsEqual(pid *PID) bool {
	return this.area == pid.area && this.id == pid.id
}

func (this *PID) String() string {
	return fmt.Sprintf("%s.%s", this.area, this.id)
}
