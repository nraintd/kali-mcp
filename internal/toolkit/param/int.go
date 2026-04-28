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
	return i.GteWithError(0, err)
}

// 大于
func (i *IntParam) Gt(v int) *IntParam {
	return i.GtWithError(v, fmt.Errorf("'%s' parameter must be greater than %d", i.k, v))
}

// 大于，自定义错误
func (i *IntParam) GtWithError(compare int, err error) *IntParam {
	return i.Validate(func(v int, k string) error {
		if v <= compare {
			return err
		}
		return nil
	})
}

// 大于等于
func (i *IntParam) Gte(v int) *IntParam {
	return i.GteWithError(v, fmt.Errorf("'%s' parameter must be greater than or equal to %d", i.k, v))
}

// 大于等于，自定义错误
func (i *IntParam) GteWithError(compare int, err error) *IntParam {
	return i.Validate(func(v int, k string) error {
		if v < compare {
			return err
		}
		return nil
	})
}

// 小于
func (i *IntParam) Lt(v int) *IntParam {
	return i.LtWithError(v, fmt.Errorf("'%s' parameter must be less than %d", i.k, v))
}

// 小于，自定义错误
func (i *IntParam) LtWithError(compare int, err error) *IntParam {
	return i.Validate(func(v int, k string) error {
		if v >= compare {
			return err
		}
		return nil
	})
}

// 小于等于
func (i *IntParam) Lte(v int) *IntParam {
	return i.LteWithError(v, fmt.Errorf("'%s' parameter must be less than or equal to %d", i.k, v))
}

// 小于等于，自定义错误
func (i *IntParam) LteWithError(compare int, err error) *IntParam {
	return i.Validate(func(v int, k string) error {
		if v > compare {
			return err
		}
		return nil
	})
}

// 等于
func (i *IntParam) Eq(v int) *IntParam {
	return i.EqWithError(v, fmt.Errorf("'%s' parameter must be equal to %d", i.k, v))
}

// 等于，自定义错误
func (i *IntParam) EqWithError(compare int, err error) *IntParam {
	return i.Validate(func(v int, k string) error {
		if v != compare {
			return err
		}
		return nil
	})
}

// 不等于
func (i *IntParam) NotEq(v int) *IntParam {
	return i.NotEqWithError(v, fmt.Errorf("'%s' parameter must not be equal to %d", i.k, v))
}

// 不等于，自定义错误
func (i *IntParam) NotEqWithError(compare int, err error) *IntParam {
	return i.Validate(func(v int, k string) error {
		if v == compare {
			return err
		}
		return nil
	})
}

// 验证器
func (i *IntParam) Validate(f func(v int, k string) error) *IntParam {
	return i.Process(func(v int, k string) (int, error) {
		return v, f(v, k)
	})
}

// 自定义处理函数
func (i *IntParam) Process(f func(v int, k string) (int, error)) *IntParam {
	if i.ps.err != nil {
		return i
	}
	v, err := f(i.v, i.k)
	if err != nil {
		i.ps.err = err
		return i
	}
	i.v = v
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
