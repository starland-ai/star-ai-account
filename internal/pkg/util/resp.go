package util

import (
	"errors"
	"fmt"
	"starland-account/internal/pkg/bizerr"
)

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func MakeResponse(data interface{}) *Response {
	return &Response{
		Code: "0",
		Msg:  "",
		Data: data,
	}
}

func MakeErrResponse(err error) *Response {
	var e *bizerr.BizError
	if errors.As(err, &e) {

		return &Response{
			Code: fmt.Sprint(e.Code()),
			Msg:  err.Error(),
			Data: nil,
		}
	}
	return &Response{
		Code: "-1",
		Msg:  err.Error(),
		Data: nil,
	}
}

func MakeResponseWithMsg(msg string) *Response {
	return &Response{
		Code: "",
		Msg:  msg,
		Data: nil,
	}
}

func (res *Response) SetCode(code string) *Response {
	res.Code = code
	return res
}
