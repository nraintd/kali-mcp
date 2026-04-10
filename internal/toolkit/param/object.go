package param

import (
	"fmt"
)

type ObjectParam struct {
	k  string
	v  map[string]any
	ps *Params
}

func NewObjectParam(key string, value map[string]any, ps *Params) *ObjectParam {
	return &ObjectParam{
		k:  key,
		v:  value,
		ps: ps,
	}
}

func (o *ObjectParam) Required() *ObjectParam {
	if o.v == nil {
		o.ps.err = fmt.Errorf("'%s' parameter cannot be null", o.k)
	}
	return o
}

func (o *ObjectParam) Validate(f func(v map[string]any, k string) error) *ObjectParam {
	if o.ps.err != nil {
		return o
	}
	if err := f(o.v, o.k); err != nil {
		o.ps.err = err
	}
	return o
}

func (o *ObjectParam) ParseString() map[string]string {
	// 如果 o.v 是 nil，len(o.v) == 0，最后会返回 map[string]string{}
	res := make(map[string]string, len(o.v))
	for k, v := range o.v {
		res[k] = fmt.Sprint(v)
	}
	o.ps.set(o.k, res)
	return res
}
