package actor

import (
	"fmt"
)

type Address struct {
	area string
	id   string
}

func NewAddress(area, id string) *Address {
	return &Address{
		area,
		id,
	}
}

func (this *Address) IsEqual(address *Address) bool {
	return this.area == address.area && this.id == address.id
}

func (this *Address) String() string {
	if this == nil {
		return "address nil"
	}
	return fmt.Sprintf("%s.%s", this.area, this.id)
}
