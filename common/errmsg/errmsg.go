package errmsg

import (
	"edu-project/consts"
	"errors"
)

type ErrMsg struct {
	ErrCode         consts.ErrorType
	ErrMsg          string
	ErrInternalInfo string
}

func NewProtocolErrorInfo(errInfo string) *ErrMsg {
	e := &ErrMsg{
		ErrCode:         consts.ErrorType_ErrorProtocol,
		ErrMsg:          consts.InternalErrMsg,
		ErrInternalInfo: errInfo,
	}
	return e
}

func NewNormalErrorInfo(errmsg string, errInfo string) *ErrMsg {
	e := &ErrMsg{
		ErrCode:         consts.ErrorType_ErrorNormal,
		ErrMsg:          errmsg,
		ErrInternalInfo: errInfo,
	}
	return e
}

func NewErrorDB(err error) *ErrMsg {
	if err == nil {
		return nil
	}
	var e *ErrMsg
	ok := errors.As(err, &e)
	if ok {
		return e
	}
	return &ErrMsg{
		ErrCode:         consts.ErrorType_ErrorDB,
		ErrMsg:          consts.InternalErrMsg,
		ErrInternalInfo: err.Error(),
	}
}
