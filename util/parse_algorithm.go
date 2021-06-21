package util

import (
	"strconv"

	"github.com/pkg/errors"
)

func GetInt(object interface{}) (resVal int64, err error) {

	resVal = 0
	err = nil

	switch object.(type) {
	case float64:
		resVal = int64(object.(float64))
	case string:
		if resVal, err = strconv.ParseInt(object.(string), 10, 64); err != nil {
			err = errors.WithStack(err)
		}
	default:
		err = errors.New("Expected integer, but use other values")
	}

	return resVal, err
}

func GetFloat(object interface{}) (resVal float64, err error) {

	resVal = 0.0
	err = nil

	switch object.(type) {
	case float64:
		resVal = object.(float64)
	case string:
		if resVal, err = strconv.ParseFloat(object.(string), 64); err != nil {
			err = errors.WithStack(err)
		}
	default:
		err = errors.New("Expected floating number, but use other values")
	}

	return resVal, err
}

func GetString(object interface{}) (resVal string, err error) {

	resVal = ""
	err = nil

	switch object.(type) {
	case string:
		resVal = object.(string)
	default:
		err = errors.New("Expected string, but use other values")
	}

	return resVal, err
}
