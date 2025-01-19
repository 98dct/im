package bitmap

type Bitmap struct {
	bits []byte
	size int
}

func NewBitmap(size int) *Bitmap {

	if size == 0 {
		size = 250
	}

	return &Bitmap{
		bits: make([]byte, size),
		size: 8 * size,
	}
}

func (b *Bitmap) Set(id string) {
	idx := hash(id) % b.size
	byteIdx := idx / 8
	bitIdx := idx % 8
	b.bits[byteIdx] |= 1 << bitIdx
}

func (b *Bitmap) IsSet(id string) bool {
	idx := hash(id) % b.size
	byteIdx := idx / 8
	bitIdx := idx % 8
	return (b.bits[byteIdx])&(1<<bitIdx) != 0
}

func (b *Bitmap) Export() []byte {
	return b.bits
}

func Load(byte []byte) *Bitmap {
	if len(byte) == 0 {
		return NewBitmap(0)
	}
	return &Bitmap{
		bits: byte,
		size: len(byte) * 8,
	}
}

func hash(id string) int {
	// 使用BKDR哈希算法
	seed := 131313 // 31 131 1313 13131 131313, etc
	hash := 0
	for _, c := range id {
		hash = hash*seed + int(c)
	}
	return hash & 0x7FFFFFFF
}
