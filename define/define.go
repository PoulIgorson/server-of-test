// Package define implements functions separate of project.
package define

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Itoa convet int to string.
var Itoa = strconv.Itoa

// Atoi convet string to int.
var Atoi = strconv.Atoi

// itoa convet int to []byte.
var Itob = func(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// GetToday returns string today date.
func GetToday() string {
	ctime := time.Now()
	cday, cmonth := ctime.Day(), int(ctime.Month())
	day, month := Itoa(cday), Itoa(cmonth)
	if len(day) < 2 {
		day = "0" + day
	}
	if len(month) < 2 {
		month = "0" + month
	}
	return fmt.Sprintf("%v-%v-%v", day, month, ctime.Year())
}

// Dict convet interface{} to fiber.Map.
func Dict(dict interface{}) fiber.Map {
	return dict.(map[string]interface{})
}

// ErrorToStr convert []error to string.
func ErrorsToStr(errs []error) string {
	errors := ""
	for _, err := range errs {
		errors += err.Error() + ", "
	}
	return errors[:len(errors)-2]
}

// Contains returns
func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// VerifyEmail verifing email
func VerifyEmail(email string) bool {
	patterns := strings.Split(email, "@")
	if len(patterns) != 2 {
		return false
	}

	name, host := patterns[0], patterns[1]
	if len(name) == 0 || len(host) == 0 {
		return false
	}

	host_patterns := strings.Split(host, ".")
	if len(host_patterns) != 2 {
		return false
	}
	hostname, hostup := host_patterns[0], host_patterns[1]
	if len(hostname) == 0 || len(hostup) == 0 {
		return false
	}

	return true
}

// Hash returns hash data
func Hash(data []byte) string {
	hasher := sha1.New()
	hasher.Write(data)
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}
