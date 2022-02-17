package common

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/dongri/phonenumber"
)

func StrFirst(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j]
		}
		i++
	}

	return s
}

// IsEmail check email the real email string
func IsEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	return re.MatchString(email)
}

// IsNumeric ...
func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)

	return err == nil
}

/**
 * Add quote to string separate by (,)
 * Input : a,b,c Output : 'a','b','c'
 */
func AddQuote(input string) string {
	inputs := strings.Split(input, ",")
	var outputs []string
	for _, val := range inputs {
		outputs = append(outputs, "'"+val+"'")
	}
	return strings.Join(outputs[:], ",")
}

func IsPhone(s string) (phonenumber.ISO3166, string, bool) {
	d := phonenumber.GetISO3166ByNumber(s, true)
	if d.CountryName == "" {
		return d, "", false
	}

	number := phonenumber.ParseWithLandLine(s, d.Alpha3)
	if number == "" {
		return d, "", false
	}

	return d, number, true
}

func CheckStringContains(value string, listValue []string) bool {
	arrStr := new(ArrStr)
	if exist, _ := arrStr.InArray(strings.Trim(value, " "), listValue); len(value) > 0 && !exist {
		return false
	}

	return true
}

func IsDate(date string) bool {
	var regex, _ = regexp.Compile(`([12]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01]))`)
	isMatch := regex.MatchString(date)

	return isMatch
}
