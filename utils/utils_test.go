package utils

import (
	"errors"
	"hashsnap-api-4/utils"
	"testing"
)

func TestCheckErr_ErrorPassed(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Passing error to CheckErr should panic")
		}
	}()
	err := errors.New("Sample error")
	utils.CheckErr(err, "", "")
}

func TestCheckErr_NilPassed(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Passing nil to CheckErr shound not panic")
		}
	}()
	utils.CheckErr(nil, "", "")
}