package proto

import (
	"maps"

	"github.com/google/uuid"
	"github.com/samber/mo"
)

type Req struct {
	Key    string         `json:"#"`
	Token  string         `json:"t"`
	Method string         `json:"m"`
	Data   map[string]any `json:"d"`
}

func NewReq() (req *Req) {
	req = new(Req)

	req.Key = uuid.New().String()
	req.Data = make(map[string]any)

	return
}

func (req *Req) Clone() (other *Req) {
	other = new(Req)

	other.Key = req.Key
	other.Token = req.Token
	other.Method = req.Method
	other.Data = maps.Clone(req.Data)

	return
}

func (req *Req) GetKey() string {
	return req.Key
}

func (req *Req) WithToken(value string) *Req {
	req.Token = value

	return req
}

func (req *Req) GetToken() string {
	return req.Token
}

func (req *Req) WithMethod(value string) *Req {
	req.Method = value

	return req
}

func (req *Req) GetMethod() string {
	return req.Method
}

func (req *Req) WithData(value map[string]any) *Req {
	req.Data = value

	return req
}

func (req *Req) GetData() map[string]any {
	return req.Data
}

func (req *Req) WithDataEntry(key string, value any) *Req {
	req.Data[key] = value

	return req
}

func (req *Req) GetDataEntry(key string) mo.Option[any] {
	value, ok := req.Data[key]

	if !ok {
		return mo.None[any]()
	} else {
		return mo.Some[any](value)
	}
}
