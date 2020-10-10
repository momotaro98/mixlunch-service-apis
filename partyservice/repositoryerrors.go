package partyservice

type RepositoryErrorCode int

type RepositoryError interface {
	error
	Code() RepositoryErrorCode
}

const (
	DuplicatePrimaryKeyErrorCode RepositoryErrorCode = iota
	NoReferenceRowErrorCode
)

var RepoErrCodeMapToRDBMS = map[RepositoryErrorCode]uint16{
	// See https://dev.mysql.com/doc/refman/5.6/ja/error-messages-server.html
	DuplicatePrimaryKeyErrorCode: 1062,
	NoReferenceRowErrorCode:      1452,
}

type DuplicatePrimaryKeyError struct {
	originalError error
}

var _ RepositoryError = (*DuplicatePrimaryKeyError)(nil)

func NewDuplicatePrimaryKeyError(err error) *DuplicatePrimaryKeyError {
	return &DuplicatePrimaryKeyError{originalError: err}
}

func (e *DuplicatePrimaryKeyError) Error() string {
	return e.originalError.Error()
}

func (e *DuplicatePrimaryKeyError) Code() RepositoryErrorCode {
	return DuplicatePrimaryKeyErrorCode
}

type NoReferenceRowError struct {
	originalError error
}

var _ RepositoryError = (*NoReferenceRowError)(nil)

func NewNoReferenceRowError(err error) *NoReferenceRowError {
	return &NoReferenceRowError{originalError: err}
}

func (e *NoReferenceRowError) Error() string {
	return e.originalError.Error()
}

func (e *NoReferenceRowError) Code() RepositoryErrorCode {
	return NoReferenceRowErrorCode
}
