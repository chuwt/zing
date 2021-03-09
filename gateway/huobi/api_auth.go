package huobi

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/chuwt/zing/object"
	"net/url"
	"sort"
	"strings"
	"time"
)

type ApiAuth struct {
	Key    string `json:"API Key"`
	Secret string `json:"Secret Key"`
}

func (a *ApiAuth) NewSignParams() object.Params {
	return object.Params{
		"AccessKeyId":      a.Key,
		"SignatureMethod":  "HmacSHA256",
		"SignatureVersion": "2",
		"Timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05"),
	}
}

func (a *ApiAuth) NewWsSignParams() object.Params {
	return object.Params{
		"accessKey":        a.Key,
		"signatureMethod":  "HmacSHA256",
		"signatureVersion": "2.1",
		"timestamp":        time.Now().UTC().Format("2006-01-02T15:04:05"),
	}
}

func (a *ApiAuth) Sign(method, host, path string, params object.Params) string {
	paramList := make([]string, 0)
	if strings.Contains(host, "wss://") {
		host = strings.ReplaceAll(host, "wss://", "")
		host = strings.Split(host, "/")[0]
	} else {
		host = strings.ReplaceAll(host, "https://", "")
	}
	for key, value := range params {
		paramList = append(paramList, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
	}
	sort.Strings(paramList)
	h := hmac.New(sha256.New, []byte(a.Secret))
	h.Write([]byte(strings.Join([]string{method, host, path, strings.Join(paramList, "&")}, "\n")))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
