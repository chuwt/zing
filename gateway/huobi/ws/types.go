package ws

import "encoding/json"

type ResData struct {
	Action  string          `json:"action"`
	Ch      string          `json:"ch"`
	Data    json.RawMessage `json:"data"`
	Code    int             `json:"code"`
	Message string          `json:"message"`
}

type Ping struct {
	Ts int64 `json:"ts"`
}

type Pong struct {
	Action string `json:"action"`
	Data   *Ping  `json:"data"`
}

type Req struct {
	Action string `json:"action"`
	Ch     string `json:"ch"`
}

type Auth struct {
	Req
	Params AuthParams `json:"params"`
}

type AuthParams struct {
	AuthType         string `json:"authType"`
	AccessKeyId      string `json:"accessKey"`
	SignatureMethod  string `json:"signatureMethod"`
	SignatureVersion string `json:"signatureVersion"`
	Timestamp        string `json:"timestamp"`
	Signature        string `json:"signature"`
}
