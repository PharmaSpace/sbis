package sbis

import (
	"fmt"
)

var (
	ErrAuthConfigNotFound = fmt.Errorf("sbis | authorization data is not specified")
	ErrEmptyREQ           = fmt.Errorf("sbis | empty required params")
)
