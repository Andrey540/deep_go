package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	err  error
	errs []error
}

func (e *MultiError) Error() string {
	var multiErrs [][]error
	var err error
	err = e
	for {
		if multiErr, ok := err.(*MultiError); ok {
			multiErrs = append(multiErrs, multiErr.errs)
			err = multiErr.Unwrap()
		} else {
			break
		}
	}
	count := 0
	result := ""
	if err != nil {
		result = fmt.Sprintf("\t* %s", err)
		count++
	}
	for i := len(multiErrs) - 1; i >= 0; i-- {
		for _, err1 := range multiErrs[i] {
			result = fmt.Sprintf("%s\t* %s", result, err1)
			count++
		}
	}
	return fmt.Sprintf("%d errors occured:\n%s\n", count, result)
}

func (e *MultiError) Unwrap() error {
	return e.err
}

func (e *MultiError) Is(err error) bool {
	for _, err1 := range e.errs {
		if errors.Is(err1, err) {
			return true
		}
	}
	return false
}

func (e *MultiError) As(err any) bool {
	for _, err1 := range e.errs {
		if errors.As(err1, &err) {
			return true
		}
	}
	return false
}

func Append(err error, errs ...error) *MultiError {
	return &MultiError{
		err:  err,
		errs: errs,
	}
}

type myError struct{}

func (e *myError) Error() string {
	return "Error"
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)

	err = Append(err, errors.New("error 3"))
	err = Append(err, errors.New("error 4"), errors.New("error 5"))
	expectedMessage = "5 errors occured:\n\t* error 1\t* error 2\t* error 3\t* error 4\t* error 5\n"
	assert.EqualError(t, err, expectedMessage)
}

func TestMultiErrorIs(t *testing.T) {
	target := errors.New("error 1")
	err := Append(target, errors.New("error 2"))
	err = Append(err, errors.New("error 3"))

	assert.True(t, errors.Is(err, target))

	err = Append(target, errors.New("error 2"))
	assert.True(t, errors.Is(err, target))

	target = errors.New("error 2")
	err = nil
	err = Append(err, errors.New("error 1"), target)
	err = Append(err, errors.New("error 3"))

	assert.True(t, errors.Is(err, target))
}

func TestMultiErrorAs(t *testing.T) {
	target := &myError{}
	err := Append(errors.New("error 1"), errors.New("error 2"), target)

	assert.True(t, errors.As(err, &target))

	err = Append(target, errors.New("error 1"))
	err = Append(errors.New("error 2"), errors.New("error 3"))

	assert.True(t, errors.As(err, &target))
}
