// File: internal/test_import.go
package internal

import (
	"fmt"
	"nekosync/internal/domain/user"
)

func Check() {
    fmt.Println(user.User{}) // should compile fine if import works
}