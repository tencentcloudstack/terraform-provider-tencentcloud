package common

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func ValidateNameRegex(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if _, err := regexp.Compile(value); err != nil {
		errors = append(errors, fmt.Errorf(
			"%q contains an invalid regular expression: %s",
			k, err))
	}
	return
}

func ValidateNotEmpty(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if value == "" {
		errors = append(errors, fmt.Errorf("%s must not use empty string: %s", k, value))
	}
	return
}

func ValidateInstanceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !strings.Contains(value, ".") {
		errors = append(errors, fmt.Errorf("the format of %s is invalid: %s, it should be like S1.SMALL1", k, value))
		return
	}
	return
}

// validateCIDRNetworkAddress ensures that the string value is a valid CIDR that
// represents a network address - it adds an error otherwise
func ValidateCIDRNetworkAddress(v interface{}, k string) (ws []string, errors []error) {
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

func ValidateIp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	ip := net.ParseIP(value)
	if ip == nil {
		errors = append(errors, fmt.Errorf("%q must contain a valid IP: %s", k, value))
	}
	return
}

func ValidateImageID(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if !strings.HasPrefix(value, "img-") {
		errors = append(errors, fmt.Errorf("the format of %q is invalid: %s, it should begin with `img-`", k, value))
	}
	return
}

// NOTE not exactly strict, but ok for now
func ValidateIntegerInRange(min, max int64) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := int64(v.(int))
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

func ValidateIntegerMin(min int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if value < min {
			errors = append(errors, fmt.Errorf(
				"%q cannot be lower than %d: %d", k, min, value))
		}
		return
	}
}

func ValidateStringLengthInRange(min, max int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		length := utf8.RuneCountInString(v.(string))
		if length < min {
			errors = append(errors, fmt.Errorf(
				"length of %q cannot be lower than %d: %d", k, min, length))
		}
		if length > max {
			errors = append(errors, fmt.Errorf(
				"length of %q cannot be higher than %d: %d", k, max, length))
		}
		return
	}
}

func ValidateKeyPairName(v interface{}, k string) (ws []string, errors []error) {
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

func ValidateAllowedStringValueIgnoreCase(ss []string) schema.SchemaValidateFunc {
	var upperStrs = make([]string, len(ss))
	for index, value := range ss {
		upperStrs[index] = strings.ToUpper(value)
	}
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !IsContains(upperStrs, strings.ToUpper(value)) {
			errors = append(errors, fmt.Errorf("%q must contain a valid string value must in array %#v, got %q", k, ss, value))
		}
		return
	}
}

func ValidateAllowedStringValue(ss []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(string)
		if !IsContains(ss, value) {
			errors = append(errors, fmt.Errorf("%q must contain a valid string value must in array %#v, got %q", k, ss, value))
		}
		return
	}
}

func ValidatePort(v interface{}, k string) (ws []string, errors []error) {
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

func ValidatePortRange(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	items := strings.Split(value, "-")
	if len(items) < 2 {
		errors = append(errors, fmt.Errorf("%q must be like minport-maxport format", k))
	}
	minPort, err := strconv.Atoi(items[0])
	if err != nil {
		errors = append(errors, err)
	}
	if minPort < 1 || minPort > 65535 {
		errors = append(errors, fmt.Errorf("%q min port must be a valid port between 1 and 65535", k))
	}
	maxPort, err := strconv.Atoi(items[1])
	if err != nil {
		errors = append(errors, err)
	}
	if maxPort < 1 || maxPort > 65535 {
		errors = append(errors, fmt.Errorf("%q max port must be a valid port between 1 and 65535", k))
	}
	if minPort > maxPort {
		errors = append(errors, fmt.Errorf("%q min port should not be greater than max port", k))
	}
	return
}

func ValidateMysqlPassword(v interface{}, k string) (ws []string, errors []error) {
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

func ValidateAllowedIntValue(ints []int) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (ws []string, errors []error) {
		value := v.(int)
		if !IsContains(ints, value) {
			errors = append(errors, fmt.Errorf("%q must contain a valid int value in array %#v, got %d", k, ints, value))
		}
		return
	}
}

// Only support lowercase letters, numbers and "-". It cannot be longer than 60 characters.
// specification: https://cloud.tencent.com/document/product/436/13312
func ValidateCosBucketName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 60 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("the length of %s must be 1-60: %s", k, value))
	}

	pattern := `^[a-z0-9]([a-z0-9-]*[a-z0-9])?-[0-9]{10}$`
	if match, _ := regexp.Match(pattern, []byte(value)); !match {
		errors = append(errors, fmt.Errorf("%s is not valid, please refer to the official documents: %s", k, value))
	}
	return
}

func ValidateCosBucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%s cannot be parsed as RFC3339 Timestamp Format: %s", k, value))
	}

	return
}

func ValidateAsConfigPassword(v interface{}, k string) (ws []string, errors []error) {
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

func ValidateAsScheduleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, value)
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%s cannot be parsed as RFC3339 Timestamp Format: %s", k, value))
	}
	return
}

// check if string has given prefix, if no one prefix matches, errors will have error
func ValidateStringPrefix(prefix ...string) schema.SchemaValidateFunc {
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
func ValidateStringSuffix(suffix ...string) schema.SchemaValidateFunc {
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

func ValidateCidrIp(v interface{}, k string) (ws []string, errs []error) {
	if _, err := ValidateIp(v, k); len(err) == 0 {
		return
	}

	if _, err := ValidateCIDRNetworkAddress(v, k); len(err) != 0 {
		errs = append(errs, fmt.Errorf("%s must be a valid IP address or a valid CIDR IP address: %s",
			k, v))
	}
	return
}

func ValidateStringNumber(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := strconv.ParseFloat(value, 64)
	if err != nil {
		errors = append(errors, fmt.Errorf("%s must be a number: %s", k, value))
	}
	return
}

func ValidateLowCase(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for _, c := range value {
		if unicode.IsUpper(c) {
			errors = append(errors, fmt.Errorf("%s must be a low case string: %s", k, value))
			return
		}
	}
	return
}

func ValidateTime(layout string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (wss []string, errs []error) {
		timeStr := v.(string)
		if _, err := time.Parse(layout, timeStr); err != nil {
			errs = append(errs, errors.Errorf("%s time format is invalid", k))
		}

		return
	}
}

func ValidateYaml(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if err := yaml.Unmarshal([]byte(value), make(map[interface{}]interface{})); err != nil {
		errors = append(errors, fmt.Errorf(
			"%s cannot be parsed as yaml Format, value: %s", k, value))
	}
	return
}

func ValidateTkeGpuDriverVersion(v interface{}, k string) (ws []string, errors []error) {
	value := v.(map[string]interface{})
	if len(value) > 0 {
		keySet := []string{"name", "version"}
		for _, paraKey := range keySet {
			if value[paraKey] == nil || strings.TrimSpace(value[paraKey].(string)) == "" {
				errors = append(errors, fmt.Errorf("%s in %s cannot be empty", paraKey, k))
			}
		}
	}
	return
}
