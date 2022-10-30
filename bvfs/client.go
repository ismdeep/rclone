package bvfs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	endpoint    string
	accessToken string
}

func NewClient() *Client {
	return &Client{
		endpoint:    "http://127.0.0.1:9000",
		accessToken: "",
	}
}

func (receiver *Client) NewRequest(method string, uri string, body interface{}) (*http.Request, error) {
	var buf []byte
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buf = raw
	}

	v := bytes.NewBuffer(buf)

	req, err := http.NewRequest(method,
		fmt.Sprintf("%v%v", receiver.endpoint, uri),
		v)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", receiver.accessToken))
	return req, nil
}

// Do function
func (receiver *Client) Do(ctx context.Context, req *http.Request, v interface{}) error {
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("Warn:", err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	type RespBody struct {
		Msg  string      `json:"msg"`
		Data interface{} `json:"data"`
	}

	if resp.StatusCode != http.StatusOK {
		var respBody RespBody
		if err := json.Unmarshal(body, &respBody); err == nil {
			return errors.New(string(body))
		}

		return errors.New(respBody.Msg)
	}

	if v != nil {
		var respBody RespBody
		if err := json.Unmarshal(body, &respBody); err != nil {
			return err
		}

		raw, err := json.Marshal(respBody.Data)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(raw, v); err != nil {
			return err
		}
	}

	return nil
}

func (receiver *Client) SignIn() *Client {
	type SignInReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type SignInResp struct {
		Token string `json:"token"`
	}

	req, err := receiver.NewRequest(http.MethodPost, "/api/v1/sign-in", &SignInReq{
		Username: "admin",
		Password: "123456",
	})
	if err != nil {
		panic(err)
	}

	var signInData SignInResp
	if err := receiver.Do(context.Background(), req, &signInData); err != nil {
		panic(err)
	}

	fmt.Println("Token:", signInData.Token)

	receiver.accessToken = signInData.Token

	return receiver
}

func (receiver *Client) MkDir(ctx context.Context, folder string) error {
	type MkDirReq struct {
		Dir string `json:"dir"`
	}
	req, err := receiver.NewRequest(http.MethodPost, "/api/v1/mkdir", &MkDirReq{Dir: folder})
	if err != nil {
		return err
	}

	if err := receiver.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}
