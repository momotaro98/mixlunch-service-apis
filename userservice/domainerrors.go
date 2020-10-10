package userservice

import (
	"fmt"

	"github.com/momotaro98/mixlunch-service-api/domainerror"
)

const (
	DuplicateUserRegisterErrorCode domainerror.ErrorCode = iota + 301
	DuplicateUserBlockRegisterErrorCode
	InconsistencyUserBlockErrorCode
)

type DuplicateUserRegisterError struct {
	userId string
}

var _ domainerror.DomainError = (*DuplicateUserRegisterError)(nil)

func NewDuplicateUserRegisterError(userId string) *DuplicateUserRegisterError {
	return &DuplicateUserRegisterError{
		userId: userId,
	}
}

func (e *DuplicateUserRegisterError) Error() string {
	return fmt.Sprintf("The user is already in DB. User ID: %s",
		e.userId)
}

func (e *DuplicateUserRegisterError) Code() domainerror.ErrorCode {
	return DuplicateUserRegisterErrorCode
}

type DuplicateUserBlockRegisterError struct {
	blocker string
	blockee string
}

var _ domainerror.DomainError = (*DuplicateUserBlockRegisterError)(nil)

func NewDuplicateUserBlockRegisterError(blocker, blockee string) *DuplicateUserBlockRegisterError {
	return &DuplicateUserBlockRegisterError{
		blocker: blocker,
		blockee: blockee,
	}
}

func (e *DuplicateUserBlockRegisterError) Error() string {
	return fmt.Sprintf("The user blocker pair is already in DB. Blocker User ID: %s, Blockee User ID: %s",
		e.blocker, e.blockee)
}

func (e *DuplicateUserBlockRegisterError) Code() domainerror.ErrorCode {
	return DuplicateUserBlockRegisterErrorCode
}

type InconsistencyUserBlockError struct {
	blocker string
	blockee string
}

var _ domainerror.DomainError = (*InconsistencyUserBlockError)(nil)

func NewInconsistencyUserBlockError(blocker, blockee string) *InconsistencyUserBlockError {
	return &InconsistencyUserBlockError{
		blocker: blocker,
		blockee: blockee,
	}
}

func (e *InconsistencyUserBlockError) Error() string {
	return fmt.Sprintf("The user blocker request has inconsistency. Check Blocker User ID: %s, Blockee User ID: %s",
		e.blocker, e.blockee)
}

func (e *InconsistencyUserBlockError) Code() domainerror.ErrorCode {
	return InconsistencyUserBlockErrorCode
}

// TODO: Add out of master ID scope (location ID, tag ID)

// TODO: Add Parse error (birthday, email?, languages?)
