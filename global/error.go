package global

import (
	"fmt"
)

type UnknowTypeError struct {
	Type interface{}
}

func (u *UnknowTypeError) Error() string {
	return fmt.Sprintf("Unknow type: %v", u.Type)
}
