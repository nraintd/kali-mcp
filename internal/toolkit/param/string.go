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

// 去除字符串首尾空白或指定字符
//
// 如果不传参，即变长参数 cutset 为空，行为等价于 strings.TrimSpace(v)
//
// 如果传参，即变长参数 cutset 不为空，行为等价于 strings.Trim(v, strings.Join(cutset, ""))
func (s *StringParam) Trim(cutset ...string) *StringParam {
	return s.Process(func(v, k string) (string, error) {
		if len(cutset) == 0 {
			return strings.TrimSpace(v), nil
		}
		return strings.Trim(v, strings.Join(cutset, "")), nil
	})
}

// 空字符串校验
func (s *StringParam) NotEmpty() *StringParam {
	return s.NotEmptyWithError(
		fmt.Errorf("'%s' parameter cannot be empty", s.k))
}

// 空字符串校验，自定义错误
func (s *StringParam) NotEmptyWithError(err error) *StringParam {
	return s.Validate(func(v, k string) error {
		if v == "" {
			return err
		}
		return nil
	})
}

// 字符串必须在给定集合中
func (s *StringParam) In(allowed ...string) *StringParam {
	return s.Validate(func(v, k string) error {
		for _, a := range allowed {
			if v == a {
				return nil
			}
		}
		return fmt.Errorf("invalid value for '%s': %s. must be one of: %s",
			k, v, strings.Join(allowed, ", "))
	})
}

// 验证器
func (s *StringParam) Validate(f func(v, k string) error) *StringParam {
	return s.Process(func(v, k string) (string, error) {
		return v, f(v, k)
	})
}

// 自定义处理函数
func (s *StringParam) Process(f func(v, k string) (string, error)) *StringParam {
	if s.ps.err != nil {
		return s
	}
	v, err := f(s.v, s.k)
	if err != nil {
		s.ps.err = err
		return s
	}
	s.v = v
	return s
}

// 获取参数值
func (s *StringParam) Parse() string {
	s.ps.set(s.k, s.v)
	return s.v
}
