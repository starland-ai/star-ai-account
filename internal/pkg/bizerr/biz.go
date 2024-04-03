package bizerr

import (
	"errors"
	"fmt"
)

type BizError struct {
	msg   string
	code  ErrCode
	cause error
}

func NewBizError(msg string, code ErrCode) *BizError {
	return &BizError{
		msg:   msg,
		code:  code,
		cause: nil,
	}
}

func (e *BizError) Code() ErrCode {
	return e.code
}

func (e *BizError) Msg() string {
	return e.msg
}

func (e *BizError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("(%d) %s: %s", e.code, e.msg, e.cause.Error())
	}
	return fmt.Sprintf("(%d) %s", e.code, e.msg)
}

func (e *BizError) Unwrap() error {
	return e.cause
}

func (e *BizError) Cause() error {
	return e.cause
}

func (e *BizError) Wrap(err error) error {
	return &BizError{
		msg:   e.msg,
		code:  e.code,
		cause: err,
	}
}

func (e *BizError) Wrapf(err error, msg string, args ...interface{}) error {
	return &BizError{
		msg:   fmt.Sprintf("%s: %s", e.msg, fmt.Sprintf(msg, args...)),
		code:  e.code,
		cause: err,
	}
}

func (e *BizError) Errorf(format string, args ...interface{}) error {
	return e.Wrap(fmt.Errorf(format, args...))
}

func ErrorToBizError(err error) (bool, *BizError) {
	var bizErr *BizError
	if errors.As(err, &bizErr) {
		return true, bizErr
	} else {
		return false, nil
	}
}
