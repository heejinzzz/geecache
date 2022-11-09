package geecache

// A ByteView holds an immutable view of bytes
type ByteView struct {
	bytes []byte
}

// Len returns the view's length
func (bv ByteView) Len() int {
	return len(bv.bytes)
}

// Bytes returns a copy of the data as a byte slice
func (bv ByteView) Bytes() []byte {
	return cloneBytes(bv.bytes)
}

func cloneBytes(bytes []byte) []byte {
	res := make([]byte, len(bytes))
	copy(res, bytes)
	return res
}

// String returns the data as a string, making a copy if necessary
func (bv ByteView) String() string {
	return string(bv.bytes)
}
