package bitmap

import (
	"math/bits"
	"strconv"
	"strings"
)

type Bitmap32 []uint32

// Set set n-th bit to 1
func (b *Bitmap32) Set(n uint32) {
	block, bit := n>>5, n%32
	b.grow(block)
	(*b)[block] |= (1 << bit)
}

// Remove set n-th bit to 0
func (b *Bitmap32) Remove(n uint32) {
	block, bit := n>>5, n%32
	if uint32(len(*b)) <= block {
		return
	}
	(*b)[block] &= ^(1 << bit)
}

// Xor invert n-th bit
func (b *Bitmap32) Xor(n uint32) {
	block, val := n>>5, n%32
	b.grow(block)
	(*b)[block] ^= (1 << val)
}

// IsEmpty check if the bitmap has any bit set to 1
func (b *Bitmap32) IsEmpty() bool {
	for i := range *b {
		if (*b)[i] > 0 {
			return false
		}
	}

	return true
}

// Has check if n-th bit is set to 1
func (b *Bitmap32) Has(n uint32) bool {
	block, val := n>>5, n%32
	if uint32(len(*b)) <= block {
		return false
	}

	return (*b)[block]&(1<<val) > 0
}

// CountDiff count different bits in two bitmaps
func (b *Bitmap32) CountDiff(b2 Bitmap32) int {
	diff := 0
	max := len(*b)
	if len(b2) > max {
		max = len(b2)
	}

	for i := 0; i < max; i++ {
		if len(b2) <= i {
			diff += bits.OnesCount32((*b)[i])
			continue
		}
		if len(*b) <= i {
			diff += bits.OnesCount32((b2)[i])
			continue
		}

		diff += bits.OnesCount32((*b)[i] ^ b2[i])
	}

	return diff
}

// Or in-place OR operation with another bitmap
func (b *Bitmap32) Or(b2 Bitmap32) {
	b.grow(uint32(len(b2) - 1))
	for i := 0; i < len(b2); i++ {
		if b2[i] == 0 {
			continue
		}
		(*b)[i] |= b2[i]
	}
}

// And in-place And operation with another bitmap
func (b *Bitmap32) And(b2 Bitmap32) {
	for i := 0; i < len(b2) && i < len(*b); i++ {
		(*b)[i] &= b2[i]
	}
}

// Shrink remove zero elements at the end of the map
func (b *Bitmap32) Shrink() {
	shrinkedIndex := len(*b)
	for i := len(*b) - 1; i >= 0; i-- {
		if (*b)[i] != 0 {
			shrinkedIndex = i + 1
			break
		}
	}

	if shrinkedIndex != len(*b) {
		newSlice := make(Bitmap32, shrinkedIndex)
		copy(newSlice, (*b)[:shrinkedIndex])
		*b = newSlice
	}
}

// Clone create a copy of the bitmap
func (b *Bitmap32) Clone() Bitmap32 {
	clone := make(Bitmap32, len(*b))
	copy(clone, *b)

	return clone
}

// Range call the passed callback with all bits set to 1.
// If the callback returns false, the method exits
func (b *Bitmap32) Range(f func(n uint32) bool) {
	for i, block := range *b {
		for block != 0 {
			tz := bits.TrailingZeros32(block)
			bitIndex := uint32(i*32 + tz)

			if !f(bitIndex) {
				return
			}

			block &= block - 1
		}
	}
}

func (b *Bitmap32) String() string {
	var sb strings.Builder

	for i := range *b {
		sb.WriteString(strconv.FormatUint(uint64((*b)[i]), 10))
		if i != len(*b)-1 {
			sb.WriteString("|")
		}
	}

	return sb.String()
}

func FromString32(str string) (Bitmap32, error) {
	if str == "" {
		return Bitmap32{}, nil
	}

	nums := strings.Split(str, "|")
	result := make(Bitmap32, 0, len(nums))
	for _, num := range nums {
		v, err := strconv.ParseUint(num, 10, 32)
		if err != nil {
			return nil, err
		}
		result = append(result, uint32(v))
	}
	return result, nil
}

func (b *Bitmap32) grow(length uint32) {
	if length+1 > uint32(len(*b)) {
		*b = append(*b, make(Bitmap32, length+1-uint32(len(*b)))...)
	}
}
