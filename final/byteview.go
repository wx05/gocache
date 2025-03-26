package final

// ByteView /*
type ByteView struct {
	b []byte
}

// /一个只读层的实现
// Len Len，实现lru.Value的方法
func (v ByteView) Len() int {
	return len(v.b)
}

func (v *ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v *ByteView) String() string {
	return string(v.b)
}
