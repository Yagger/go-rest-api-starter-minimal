package utils

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"invoicing-app/conf"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"time"
)

func CheckErr(err error, controller string, info string) {
	if err != nil {
		msg := fmt.Sprintf("ERRORR | %s | %s | %s", controller, info, err.Error())
		panic(msg)
	}
}

func RandString(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	letterRunes := []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetBody(r *http.Request, obj interface{}) ([]byte, error) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bodyBytes, obj)
	return bodyBytes, err
}

type Payload interface {
	Parse(r *http.Request) (validationErrors string, err error)
}

func ParsePayload(w http.ResponseWriter, r *http.Request, payload Payload) (isValid bool) {
	validationErrors, err := payload.Parse(r)
	if err != nil {
		validationErrors = "Problem with payload data. Make sure to pass a valid JSON."
	}
	if validationErrors != "" {
		w.WriteHeader(400)
		fmt.Fprintf(w, validationErrors)
		return false
	}
	return true
}

func IsEmailValid(value string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$", value)
	return match
}

func IsPasswordValid(value string) bool {
	var (
		match bool
	)
	if match, _ = regexp.MatchString("^.{8,}$", value); !match {
		return false
	}
	if match, _ = regexp.MatchString("[a-z]+", value); !match {
		return false
	}
	if match, _ = regexp.MatchString("[A-Z]+", value); !match {
		return false
	}
	if match, _ = regexp.MatchString("[0-9]+", value); !match {
		return false
	}
	return true
}

func HashPassword(password, salt string) string {
	key := []byte(conf.SessionSecret)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(password + salt))
	hash := mac.Sum(nil)
	hexString := hex.EncodeToString(hash)
	return hexString
}

func TimeLayoutChange(s string) string {
	var (
		err        error
		t          time.Time
		fromLayout = "2006-01-02T15:04:05Z"
		toLayout   = "2006-01-02 15:04:05"
		result     string
	)
	t, err = time.Parse(fromLayout, s)
	CheckErr(err, "TimeLayoutChange", fmt.Sprintf("Failed to parse time string \"%s\" into Time object", s))
	result = t.Format(toLayout)
	// Dealing with the weird Golang zero time parsing from Mysql. It parses 0000-00-00 00:00:00 as 0001-01-01 00:00:00, which is not what most clients expect.
	if result == "0001-01-01 00:00:00" {
		result = "0000-00-00 00:00:00"
	}
	return result
}
