package param

type BoolParam struct {
	k  string
	v  bool
	ps *Params
}

func NewBoolParam(key string, value bool, ps *Params) *BoolParam {
	return &BoolParam{
		k:  key,
		v:  value,
		ps: ps,
	}
}

// 获取参数值
func (b *BoolParam) Parse() bool { return b.v }
