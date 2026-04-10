package param

import (
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

type Params struct {
	req    mcp.CallToolRequest
	err    error
	params map[string]any
}

func NewParams(req mcp.CallToolRequest) *Params {
	return &Params{
		req:    req,
		params: make(map[string]any),
	}
}

func (ps *Params) Get(key string) any {
	return ps.params[key]
}

func (ps *Params) set(key string, value any) {
	ps.params[key] = value
}

func (ps *Params) String(key, defaultVal string) *StringParam {
	return NewStringParam(key, ps.req.GetString(key, defaultVal), ps)
}

func (ps *Params) Int(key string, defaultValue int) *IntParam {
	return NewIntParam(key, ps.req.GetInt(key, defaultValue), ps)

}

func (ps *Params) Bool(key string, defaultValue bool) *BoolParam {
	return NewBoolParam(key, ps.req.GetBool(key, defaultValue), ps)
}

func (ps *Params) Object(key string, defaultVal map[string]any) *ObjectParam {
	if defaultVal == nil {
		panic("defaultVal cannot be nil")
	}

	if valRaw, ok := ps.req.GetArguments()[key]; ok {
		if val, ok := valRaw.(map[string]any); ok {
			return NewObjectParam(key, val, ps)
		}
		ps.err = fmt.Errorf("'%s' parameter must be an object", key)
	}
	return NewObjectParam(key, defaultVal, ps)
}

func (ps *Params) Validate(f func(p *Params) error) *Params {
	if ps.err != nil {
		return ps
	}
	if err := f(ps); err != nil {
		ps.err = err
	}
	return ps
}

// MutuallyExclusive 验证两个参数互斥（必须且只能存在一个）
func (ps *Params) MutuallyExclusive(key1, key2 string) *Params {
	if ps.err != nil {
		return ps
	}

	v1 := ps.params[key1]
	v2 := ps.params[key2]

	isV1Empty := IsEmpty(v1)
	isV2Empty := IsEmpty(v2)

	if isV1Empty && isV2Empty {
		ps.err = fmt.Errorf("either '%s' or '%s' must be provided", key1, key2)
	} else if !isV1Empty && !isV2Empty {
		ps.err = fmt.Errorf("'%s' and '%s' are mutually exclusive", key1, key2)
	}

	return ps
}

func (ps *Params) Err() error {
	return ps.err
}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}
	switch val := v.(type) {
	case string:
		return val == ""
	case int:
		return val == 0
	case bool:
		return !val
	case map[string]any:
		return len(val) == 0
	case map[string]string:
		return len(val) == 0
	}
	return false
}
