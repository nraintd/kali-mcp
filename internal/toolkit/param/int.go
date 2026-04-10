package param

import "fmt"

type IntParam struct {
	k  string
	v  int
	ps *Params
}

func NewIntParam(key string, value int, ps *Params) *IntParam {
	return &IntParam{
		k:  key,
		v:  value,
		ps: ps,
	}
}

// 非零校验
func (i *IntParam) NotZero() *IntParam {
	return i.NotZeroWithError(fmt.Errorf("'%s' parameter cannot be zero", i.k))
}

// 非零校验，自定义错误
func (i *IntParam) NotZeroWithError(err error) *IntParam {
	return i.NotEqWithError(0, err)
}

// 非负校验
func (i *IntParam) NotNegative() *IntParam {
	return i.NotNegativeWithError(fmt.Errorf("'%s' parameter cannot be negative", i.k))
}

// 非负校验，自定义错误
func (i *IntParam) NotNegativeWithError(err error) *IntParam {
	return i.GtWithError(0, err)
}

// 大于
func (i *IntParam) Gt(v int) *IntParam {
	return i.GtWithError(v, fmt.Errorf("'%s' parameter must be greater than %d", i.k, v))
}

// 大于，自定义错误
func (i *IntParam) GtWithError(v int, err error) *IntParam {
	if i.ps.err != nil {
		return i
	}
	if i.v <= v {
		i.ps.err = err
	}
	return i
}

// 小于
func (i *IntParam) Lt(v int) *IntParam {
	return i.LtWithError(v, fmt.Errorf("'%s' parameter must be less than %d", i.k, v))
}

// 小于，自定义错误
func (i *IntParam) LtWithError(v int, err error) *IntParam {
	if i.ps.err != nil {
		return i
	}
	if i.v >= v {
		i.ps.err = err
	}
	return i
}

// 等于
func (i *IntParam) Eq(v int) *IntParam {
	return i.EqWithError(v, fmt.Errorf("'%s' parameter must be equal to %d", i.k, v))
}

// 等于，自定义错误
func (i *IntParam) EqWithError(v int, err error) *IntParam {
	if i.ps.err != nil {
		return i
	}
	if i.v != v {
		i.ps.err = err
	}
	return i
}

// 不等于
func (i *IntParam) NotEq(v int) *IntParam {
	return i.NotEqWithError(v, fmt.Errorf("'%s' parameter must not be equal to %d", i.k, v))
}

// 不等于，自定义错误
func (i *IntParam) NotEqWithError(v int, err error) *IntParam {
	if i.ps.err != nil {
		return i
	}
	if i.v == v {
		i.ps.err = err
	}
	return i
}

// 验证器
func (i *IntParam) Validate(f func(v int, k string) error) *IntParam {
	return i.Custom(func(v int, k string) (int, error) {
		return i.v, f(v, k)
	})
}

// 自定义处理函数
func (i *IntParam) Custom(f func(v int, k string) (int, error)) *IntParam {
	if i.ps.err != nil {
		return i
	}
	vv, err := f(i.v, i.k)
	if err != nil {
		i.ps.err = err
		return i
	}
	i.v = vv
	return i
}

// 获取参数值
func (i *IntParam) Parse() int {
	i.ps.set(i.k, i.v)
	return i.v
}

// 获取参数值，转换为 string
func (i *IntParam) ParseToString() string {
	i.ps.set(i.k, i.v)
	return fmt.Sprint(i.v)
}
