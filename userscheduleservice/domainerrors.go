package userscheduleservice

import (
	"fmt"
	"time"

	"github.com/momotaro98/mixlunch-service-api/domainerror"
)

const (
	InvalidDateTimeFormatCode domainerror.ErrorCode = iota + 100
	FromIsAfterToErrorCode
	DifferentDayFromAndToErrorCode
	DuplicateInOneDayErrorCode
	TheScheduleNotFoundErrorCode
	TimeRangeIsLessThanSpecifiedErrorCode
)

// InvalidDateTimeFormat

type InvalidDateTimeFormatError struct {
	SpecifiedDateTimeStr string
}

func NewInvalidDateTimeFormatError(specifiedDateTimeStr string) *InvalidDateTimeFormatError {
	return &InvalidDateTimeFormatError{
		SpecifiedDateTimeStr: specifiedDateTimeStr,
	}
}

func (e *InvalidDateTimeFormatError) Error() string {
	return fmt.Sprintf("Specified datetime format is wrong. Use RFC3339 format. Your datetime: %s",
		e.SpecifiedDateTimeStr)
}

func (e *InvalidDateTimeFormatError) Code() domainerror.ErrorCode {
	return InvalidDateTimeFormatCode
}

// FromIsAfterToError

type FromIsAfterToError struct {
	FromDate time.Time
	ToDate   time.Time
}

func NewFromIsAfterToError(fromDate time.Time, toDate time.Time) *FromIsAfterToError {
	return &FromIsAfterToError{
		FromDate: fromDate,
		ToDate:   toDate,
	}
}

func (e *FromIsAfterToError) Error() string {
	return fmt.Sprintf(
		"fromDateTime is after toDateTime of the user schedule. fromDateTime: %s, toDateTime: %s",
		e.FromDate.String(), e.ToDate.String())
}

func (e *FromIsAfterToError) Code() domainerror.ErrorCode {
	return FromIsAfterToErrorCode
}

// DifferentDayFromAndToError

type DifferentDayFromAndToError struct {
	FromDate time.Time
	ToDate   time.Time
}

func NewDifferentDayFromAndToError(fromDate time.Time, toDate time.Time) *DifferentDayFromAndToError {
	return &DifferentDayFromAndToError{
		FromDate: fromDate,
		ToDate:   toDate,
	}
}

func (e *DifferentDayFromAndToError) Error() string {
	return fmt.Sprintf(
		"Days fromDateTime and endDateTime of the user schedule are different. fromDateTime: %s, toDateTime: %s",
		e.FromDate.String(), e.ToDate.String())
}

func (e *DifferentDayFromAndToError) Code() domainerror.ErrorCode {
	return DifferentDayFromAndToErrorCode
}

// DuplicateInOneDayError

type DuplicateInOneDayError struct {
	TargetDate time.Time
}

func NewDuplicateInOneDayError(targetDate time.Time) *DuplicateInOneDayError {
	return &DuplicateInOneDayError{
		TargetDate: targetDate,
	}
}

func (e *DuplicateInOneDayError) Error() string {
	return fmt.Sprintf("There is already a user schedule in the day. The user schedule date: %s",
		e.TargetDate.String())
}

func (e *DuplicateInOneDayError) Code() domainerror.ErrorCode {
	return DuplicateInOneDayErrorCode
}

// TheScheduleNotFoundError

type TheScheduleNotFoundError struct {
	TargetDate time.Time
}

func NewTheScheduleNotFoundError(targetDate time.Time) *TheScheduleNotFoundError {
	return &TheScheduleNotFoundError{
		TargetDate: targetDate,
	}
}

func (e *TheScheduleNotFoundError) Error() string {
	return fmt.Sprintf("The specified day doesn't have a user schedule. The specified date: %s",
		e.TargetDate.String())
}

func (e *TheScheduleNotFoundError) Code() domainerror.ErrorCode {
	return TheScheduleNotFoundErrorCode
}

// TimeRangeIsLessThanSpecifiedError

type TimeRangeIsLessThanSpecifiedError struct {
	FromDate time.Time
	ToDate   time.Time
}

func NewTimeRangeIsLessThanSpecifiedError(fromDate time.Time, toDate time.Time) *TimeRangeIsLessThanSpecifiedError {
	return &TimeRangeIsLessThanSpecifiedError{
		FromDate: fromDate,
		ToDate:   toDate,
	}
}

func (e *TimeRangeIsLessThanSpecifiedError) Error() string {
	return fmt.Sprintf(
		"Time range between fromDateTime and toDateTime must not be within %d minutes. fromDateTime: %s, toDateTime: %s",
		regulatedTimeDurationMinutes, e.FromDate.String(), e.ToDate.String())
}

func (e *TimeRangeIsLessThanSpecifiedError) Code() domainerror.ErrorCode {
	return TimeRangeIsLessThanSpecifiedErrorCode
}
