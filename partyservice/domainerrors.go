package partyservice

import (
	"fmt"

	"github.com/momotaro98/mixlunch-service-api/domainerror"
)

const (
	InvalidDateTimeFormatCode domainerror.ErrorCode = 200 + iota
	DuplicateReviewErrorCode
	InconsistencyReviewErrorCode
)

// InvalidDateTimeFormat

type InvalidDateTimeFormatError struct {
	SpecifiedDateTimeStr string
}

var _ domainerror.DomainError = (*InvalidDateTimeFormatError)(nil)

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

type DuplicateReviewError struct {
	PartyID  int
	Reviewer string
	Reviewee string
}

var _ domainerror.DomainError = (*DuplicateReviewError)(nil)

func NewDuplicateReviewError(partyID int, reviewer, reviewee string) *DuplicateReviewError {
	return &DuplicateReviewError{
		PartyID:  partyID,
		Reviewer: reviewer,
		Reviewee: reviewee,
	}
}

func (e *DuplicateReviewError) Error() string {
	return fmt.Sprintf("The review is already posted. party_id: %d, reviewer: %s, reviewee: %s",
		e.PartyID, e.Reviewer, e.Reviewee,
	)
}

func (e *DuplicateReviewError) Code() domainerror.ErrorCode {
	return DuplicateReviewErrorCode
}

type InconsistencyReviewError struct {
	PartyID  int
	Reviewer string
	Reviewee string
}

var _ domainerror.DomainError = (*InconsistencyReviewError)(nil)

func NewInconsistencyReviewError(partyID int, reviewer, reviewee string) *InconsistencyReviewError {
	return &InconsistencyReviewError{
		PartyID:  partyID,
		Reviewer: reviewer,
		Reviewee: reviewee,
	}
}

func (e *InconsistencyReviewError) Error() string {
	return fmt.Sprintf("The review post has incosistency. Check party_id: %d, reviewer: %s, reviewee: %s",
		e.PartyID, e.Reviewer, e.Reviewee,
	)
}

func (e *InconsistencyReviewError) Code() domainerror.ErrorCode {
	return InconsistencyReviewErrorCode
}
