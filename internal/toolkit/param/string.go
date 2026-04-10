package param

import (
	"fmt"
	"strings"
)

type StringParam struct {
	k  string
	v  string
	ps *Params
}

func NewStringParam(key, value string, ps *Params) *StringParam {
	return &StringParam{
		k:  key,
		v:  value,
		ps: ps,
	}
}

// 去除字符串两端的空白字符
func (s *StringParam) TrimSpace() *StringParam {
	if s.ps.err != nil {
		return s
	}
	s.v = strings.TrimSpace(s.v)
	return s
}

// 空字符串校验
func (s *StringParam) NotEmpty() *StringParam {
	return s.NotEmptyWithError(fmt.Errorf("'%s' parameter cannot be empty", s.k))
}

// 空字符串校验，自定义错误
func (s *StringParam) NotEmptyWithError(err error) *StringParam {
	if s.ps.err != nil {
		return s
	}
	if s.v == "" {
		s.ps.err = err
	}
	return s
}

// 验证器
func (s *StringParam) Validate(f func(v, k string) error) *StringParam {
	return s.Custom(func(v, k string) (string, error) {
		return s.v, f(v, k)
	})
}

// 自定义处理函数
func (s *StringParam) Custom(f func(v, k string) (string, error)) *StringParam {
	if s.ps.err != nil {
		return s
	}
	vv, err := f(s.v, s.k)
	if err != nil {
		s.ps.err = err
		return s
	}
	s.v = vv
	return s
}

// 获取参数值
func (s *StringParam) Parse() string {
	s.ps.set(s.k, s.v)
	return s.v
}
