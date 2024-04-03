package bizerr

type ErrCode int32

func (c ErrCode) Int32() int32 {
	return int32(c)
}

const (
	AuthenticationFailed   ErrCode = 401
	BadRequest             ErrCode = 65531
	InternalError          ErrCode = 65534
	TimeOut                ErrCode = 65539
	VerificationCodeFailed ErrCode = 65540
	NotExist               ErrCode = 65541
	PostureNotExist        ErrCode = 65542
)

var (
	ErrInternalError          = NewBizError("internal bizerr", InternalError)
	ErrCheckQrResultError     = NewBizError("qr code check bizerr", InternalError)
	ErrVerificationCodeFailed = NewBizError("check verification code failed", VerificationCodeFailed)
	ErrModelNotExist          = NewBizError("model not exists", NotExist)
	ErrPostureNotExist        = NewBizError("posture not exists", NotExist)
	ErrClothingNotExist       = NewBizError("clothing not exists", NotExist)
	ErrAccountNotExist        = NewBizError("account not exists", NotExist)
	ErrActivityNotExist       = NewBizError("activity not exists", NotExist)
)
