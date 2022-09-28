package gee_cache

// ByteView 用于存放缓存值
type ByteView struct {
	b []byte // readonly
}

// Len 实现了Value接口，返回缓存值大小
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice 返回b的拷贝，防止缓存值被外部修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// String ...
func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
