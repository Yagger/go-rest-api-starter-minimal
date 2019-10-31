package utils

import (
	"fmt"
)

// CheckErr func
func CheckErr(err error, controller string, info string) {
	if err != nil {
		msg := fmt.Sprintf("ERROR | %s | %s | %s", controller, info, err.Error())
		panic(msg)
	}
}
