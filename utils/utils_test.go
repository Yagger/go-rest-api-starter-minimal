package utils

import (
	"errors"
	"testing"
)

func TestCheckErr_ErrorPassed(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Passing error to CheckErr should panic")
		}
	}()
	err := errors.New("Sample error")
	CheckErr(err, "", "")
}

func TestCheckErr_NilPassed(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Error("Passing nil to CheckErr shound not panic")
		}
	}()
	CheckErr(nil, "", "")
}

func TestIsEmailValid(t *testing.T) {
	if IsEmailValid("aaa") {
		t.Error("Email \"aaa\" should not be valid because it is not real email address")
	}
	if IsEmailValid("aaa@bbb.") {
		t.Error("Email \"aaa@bbb.\" should not be valid because it is not real email address")
	}
	if !IsEmailValid("aaa@bbb.com") {
		t.Error("Email \"aaa@bbb.com\" should be valid because it is real email address")
	}
	if IsEmailValid("aa!@#a@bbb.com") {
		t.Error("Email \"aa!@#a@bbb.com\" should not be valid because it is not real email address")
	}
}

func TestIsPasswordValid(t *testing.T) {
	if IsPasswordValid("aaa") {
		t.Error("Password \"aaa\" should not be valid password because it is less than 8 characters long")
	}
	if IsPasswordValid("aaabbbccc") {
		t.Error("Password \"aaabbbccc\" should not be valid password because it has no uppercase letters")
	}
	if IsPasswordValid("AAABBBCCC") {
		t.Error("Password \"AAABBBCCC\" should not be valid password because it has no lowercase letters")
	}
	if IsPasswordValid("AAAbbbCCC") {
		t.Error("Password \"AAAbbbCCC\" should not be valid password because it has no digits")
	}
	if IsPasswordValid("0123456789") {
		t.Error("Password \"0123456789\" should not be valid password because it has no letters")
	}
	if !IsPasswordValid("AAAbbb000") {
		t.Error("Password \"AAAbbb000\" should be valid because it is at least 8 characters long, has upper and lowercase letters and digits")
	}
}
