package utils

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

func IsEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
func IsPhoneValid(phone string) bool {
	phoneRegex := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return phoneRegex.MatchString(strings.TrimSpace(phone))
}

func IsUrlValid(URL string) bool {
	_, err := url.ParseRequestURI(URL)
	return err == nil
}

func IsValidUUID(input uuid.UUID) bool {
	_, err := uuid.Parse(input.String())
	return input != uuid.UUID{00000000 - 0000 - 0000 - 0000 - 000000000000} && err == nil
}

func ConvertIntToUUID(value int) (uuid.UUID, error) {
	intStr := strconv.Itoa(value)
	uuidFromInt, err := uuid.Parse(intStr)
	if err != nil {
		return uuid.Nil, err
	}
	return uuidFromInt, nil
}
