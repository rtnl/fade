package proto

import (
	"github.com/samber/mo"
)

type ResCode int32

const (
	ResCodeNone ResCode = iota
	ResCodeOk
	ResCodeFail
)

type Res struct {
	Key  string         `json:"#"`
	Code ResCode        `json:"c"`
	Data map[string]any `json:"d"`
}

func NewRes(code ResCode, key string) (res *Res) {
	res = new(Res)

	res.Key = key
	res.Code = code
	res.Data = make(map[string]any)

	return
}

func (res *Res) WithKey(value string) *Res {
	res.Key = value

	return res
}

func (res *Res) GetKey() string {
	return res.Key
}

func (res *Res) WithCode(value ResCode) *Res {
	res.Code = value

	return res
}

func (res *Res) GetCode() ResCode {
	return res.Code
}

func (res *Res) WithData(value map[string]any) *Res {
	res.Data = value

	return res
}

func (res *Res) GetData() map[string]any {
	return res.Data
}

func (res *Res) WithDataEntry(key string, value any) *Res {
	res.Data[key] = value

	return res
}

func (res *Res) GetDataEntry(key string) mo.Option[any] {
	value, ok := res.Data[key]

	if !ok {
		return mo.None[any]()
	} else {
		return mo.Some[any](value)
	}
}
