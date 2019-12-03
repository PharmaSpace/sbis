package sbis

import (
	"fmt"
)

var (
	ErrAuthConfigNotFound = fmt.Errorf("sbis | authorization data is not specified")
	ErrEmptyINN           = fmt.Errorf("sbis | empty required inn params")
)
