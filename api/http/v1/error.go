package v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrResp struct {
	Status int
	Code   int
	Msg    string
	Err    error
}

func (r *ErrResp) MarshalJSON() ([]byte, error) {
	msg := fmt.Sprintf("%s, details: %s", r.Msg, r.Err)
	data := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{
		Code: r.Code,
		Msg:  msg,
	}
	return json.Marshal(data)
}

func (r *ErrResp) WithMsgf(msg string, args ...interface{}) *ErrResp {
	r1 := *r
	r1.Msg = fmt.Sprintf(msg, args...)
	return &r1
}

func (r *ErrResp) WithErr(err error) *ErrResp {
	r1 := *r
	r1.Err = err
	return &r1
}

func (r *ErrResp) WithErrorf(format string, args ...interface{}) *ErrResp {
	return r.WithErr(fmt.Errorf(format, args...))
}

func (r *ErrResp) Abort(c *gin.Context) {
	if r.Err != nil {
		e := c.Error(r.Err)
		r.Err = e
	}
	c.AbortWithStatusJSON(r.Status, r)
}

var (
	ErrRespBadRequest = &ErrResp{
		Status: http.StatusBadRequest,
		Code:   http.StatusBadRequest,
		Msg:    http.StatusText(http.StatusBadRequest),
	}
	ErrRespInternalError = &ErrResp{
		Status: http.StatusInternalServerError,
		Code:   http.StatusInternalServerError,
		Msg:    http.StatusText(http.StatusInternalServerError),
	}
)