package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"starland-account/configs"
	"starland-account/internal/pkg/httpclientutil"

	"go.uber.org/zap"
)

type SendAlertMsgRequest struct {
	Content Content `json:"content"`
	MsgType string  `json:"msg_type"`
}

type Content struct {
	Text string `json:"text"`
}

type SendAlertMsgResponse struct {
	Code int64                  `json:"code"`
	Data map[string]interface{} `json:"data"`
	Msg  string                 `json:"msg"`
}

func SendAlertMsg(msg string) error {
	var err error
	ctx := context.Background()
	msg = fmt.Sprintf("[%s]: %s", configs.GetConfig().Env, msg)
	reqData := &SendAlertMsgRequest{MsgType: "text", Content: Content{Text: msg}}
	zap.S().Infof("reqData:%+v", reqData)
	reqBuf := new(bytes.Buffer)
	err = json.NewEncoder(reqBuf).Encode(reqData)
	if err != nil {
		return fmt.Errorf("req encode: %w", err)
	}
	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, "POST", configs.GetConfig().FeiShuAlertURL, reqBuf)
	if err != nil {
		return fmt.Errorf("post tts svc err: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	var resp *http.Response
	cl := httpclientutil.GetHttpClient()
	resp, err = cl.Do(req)
	if err != nil {
		return fmt.Errorf("cl.do err: %w", err)
	}

	defer func() {
		if e := resp.Body.Close(); e != nil {
			fmt.Println(e)
		}
	}()

	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}

	var res *SendAlertMsgResponse
	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return fmt.Errorf("json decode err: %w", err)
	}
	if res.Code != 0 {
		return fmt.Errorf("res err: %s", res.Msg)
	}

	return nil
}
