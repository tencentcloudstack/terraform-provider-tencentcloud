package tencentcloud

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/likexian/gokit/assert"
)

func validateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}
	return
}

func validateNotEmpty(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) == 0 {
		errors = append(errors, fmt.Errorf("%s must not use empty string: %s", k, value))
	}
	return
}

func validateInstanceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	words := strings.Split(value, ".")
	if len(words) <= 1 {
		errors = append(errors, fmt.Errorf("the format of %s is invalid: %s, it should be like S1.SMALL1", k, value))
		return
	}
	return
}

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func validateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q must contain a valid CIDR, got error parsing: %s", k, err))
		return
	}
	if ipnet == nil || value != ipnet.String() {
		errors = append(errors, fmt.Errorf("%q must contain a valid network CIDR, expected %q, got %q", k, ipnet, value))
	}
	return
}

func validateIp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ip := net.ParseIP(value)
	if ip == nil {
		errors = append(errors, fmt.Errorf("%q must contain a valid IP: %s", k, value))
	}
	return
}

// NOTE not exactly strict, but ok for now
func validateIntegerInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"%q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

func validateIntegerMin(min int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		return
	}
}

func validateStringLengthInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := len(v.(string))
		if value < min {
			errors = append(errors, fmt.Errorf(
				"length of %q cannot be lower than %d: %d", k, min, value))
		}
		if value > max {
			errors = append(errors, fmt.Errorf(
				"length of %q cannot be higher than %d: %d", k, max, value))
		}
		return
	}
}

func validateKeyPairName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 25 || len(value) == 0 {
		errors = append(errors, fmt.Errorf("the length of %s must be 1-25: %s", k, value))
	}

	pattern := `^[a-zA-Z0-9_]+$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%s only support letters, numbers and \"_\": %s", k, value))
	}
	return
}

func validateAllowedStringValueIgnoreCase(ss []string) schema.SchemaValidateFunc {
	var upperStrs = make([]string, len(ss))
	for index, value := range ss {
		upperStrs[index] = strings.ToUpper(value)
	}
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !assert.IsContains(upperStrs, strings.ToUpper(value)) {
			errors = append(errors, fmt.Errorf("%q must contain a valid string value must in array %#v, got %q", k, ss, value))
		}
		return
	}
}

func validateAllowedStringValue(ss []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !assert.IsContains(ss, value) {
			errors = append(errors, fmt.Errorf("%q must contain a valid string value must in array %#v, got %q", k, ss, value))
		}
		return
	}
}

func validatePort(v interface{}, k string) (ws []string, errors []error) {
	value := 0
	switch t := v.(type) {
	case string:
		value, _ = strconv.Atoi(t)
	case int:
		value = t
	default:
		errors = append(errors, fmt.Errorf("%q data type error ", k))
		return
	}
	if value < 1 || value > 65535 {
		errors = append(errors, fmt.Errorf("%q must be a valid port between 1 and 65535", k))
	}
	return
}

func validateMysqlPassword(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 || len(value) < 8 {
		// errors = append(errors, fmt.Errorf("%s invalid password, len(password) must between 8 and 64,%s", k, value))
		errors = append(errors, fmt.Errorf("the length of %s must be 8-64: %s", k, value))
	}
	var match = make(map[string]bool)
	if strings.ContainsAny(value, "_+-&=!@#$%^*()") {
		match["alien"] = true
	}
	for i := 0; i < len(value); i++ {
		if len(match) >= 2 {
			break
		}
		if value[i] >= '0' && value[i] <= '9' {
			match["number"] = true
			continue
		}
		if (value[i] >= 'a' && value[i] <= 'z') || (value[i] >= 'A' && value[i] <= 'Z') {
			match["letter"] = true
			continue
		}
	}
	if len(match) < 2 {
		errors = append(errors, fmt.Errorf("the format of %s is invalid: %s, it must contain at least letters, Numbers, and characters(_+-&=!@#$%%^*())", k, value))
	}
	return
}

func validateAllowedIntValue(ints []int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if !assert.IsContains(ints, value) {
			errors = append(errors, fmt.Errorf("%q must contain a valid int value in array %#v, got %d", k, ints, value))
		}
		return
	}
}

// Only support lowercase letters, numbers and "-". It cannot be longer than 40 characters.
func validateCosBucketName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 40 || len(value) < 0 {
		errors = append(errors, fmt.Errorf("the length of %s must be 1-40: %s", k, value))
	}

	pattern := `^[a-z0-9-]+$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%s only support lowercase letters, numbers and \"-\": %s", k, value))
	}
	return
}

func validateCosBucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%s cannot be parsed as RFC3339 Timestamp Format: %s", k, value))
	}

	return
}

func validateAsConfigPassword(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 16 || len(value) < 8 {
		errors = append(errors, fmt.Errorf("the length of %s must be 8-16: %s", k, value))
	}
	var match = make(map[string]bool)
	if strings.ContainsAny(value, "()~!@#$%^&*-+={}[]:;',.?/") {
		match["alien"] = true
	}
	for i := 0; i < len(value); i++ {
		if len(match) >= 2 {
			break
		}
		if value[i] >= '0' && value[i] <= '9' {
			match["number"] = true
			continue
		}
		if (value[i] >= 'a' && value[i] <= 'z') || (value[i] >= 'A' && value[i] <= 'Z') {
			match["letter"] = true
			continue
		}
	}
	if len(match) < 2 {
		errors = append(errors, fmt.Errorf("%s is invalid, it must contains at least letters, numbers, and characters(_+-&=!@#$%%^*()): %s", k, value))
	}
	return
}

func validateAsScheduleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%s cannot be parsed as RFC3339 Timestamp Format: %s", k, value))
	}
	return
}

// check if string has given prefix, if no one prefix matches, errors will have error
func validateStringPrefix(prefix ...string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		for _, p := range prefix {
			if strings.HasPrefix(value, p) {
				return
			}
		}
		errors = append(errors, fmt.Errorf("%s doesn't have preifx %v", k, prefix))
		return
	}
}

// check if string has given suffix, if no one suffix matches, errors will have error
func validateStringSuffix(suffix ...string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		for _, s := range suffix {
			if strings.HasSuffix(value, s) {
				return
			}
		}
		errors = append(errors, fmt.Errorf("%s doesn't have suffix %v", k, suffix))
		return
	}
}

func validateCidrIp(v interface{}, k string) (ws []string, errs []error) {
	if _, err := validateIp(v, k); len(err) == 0 {
		return
	}

	if _, err := validateCIDRNetworkAddress(v, k); len(err) != 0 {
		errs = append(errs, fmt.Errorf("%s must be a valid IP address or a valid CIDR IP address: %s",
			k, v))
	}
	return
}

func validateStringNumber(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		errors = append(errors, fmt.Errorf("%s must be a number: %s", k, value))
	}
	return
}
