package bitmap

import (
	"math/bits"
	"strconv"
	"strings"
)

type Bitmap Bitmap64

type Bitmap64 []uint64

// Set set n-th bit to 1
func (b *Bitmap64) Set(n uint32) {
	block, bit := n>>6, n%64
	b.grow(block)
	(*b)[block] |= (1 << bit)
}

// Remove set n-th bit to 0
func (b *Bitmap64) Remove(n uint32) {
	block, bit := n>>6, n%64
	if uint32(len(*b)) <= block {
		return
	}
	(*b)[block] &= ^(1 << bit)
}

// Xor invert n-th bit
func (b *Bitmap64) Xor(n uint32) {
	block, val := n>>6, n%64
	b.grow(block)
	(*b)[block] ^= (1 << val)
}

// IsEmpty check if the bitmap has any bit set to 1
func (b *Bitmap64) IsEmpty() bool {
	for i := range *b {
		if (*b)[i] > 0 {
			return false
		}
	}

	return true
}

// Has check if n-th bit is set to 1
func (b *Bitmap64) Has(n uint32) bool {
	block, val := n>>6, n%64
	if uint32(len(*b)) <= block {
		return false
	}

	return (*b)[block]&(1<<val) > 0
}

// CountDiff count different bits in two bitmaps
func (b *Bitmap64) CountDiff(b2 Bitmap64) int {
	diff := 0
	max := len(*b)
	if len(b2) > max {
		max = len(b2)
	}

	for i := 0; i < max; i++ {
		if len(b2) <= i {
			diff += bits.OnesCount64((*b)[i])
			continue
		}
		if len(*b) <= i {
			diff += bits.OnesCount64((b2)[i])
			continue
		}

		diff += bits.OnesCount64((*b)[i] ^ b2[i])
	}

	return diff
}

// Or in-place OR operation with another bitmap
func (b *Bitmap64) Or(b2 Bitmap64) {
	b.grow(uint32(len(b2) - 1))
	for i := 0; i < len(b2); i++ {
		if b2[i] == 0 {
			continue
		}
		(*b)[i] |= b2[i]
	}
}

// And in-place And operation with another bitmap
func (b *Bitmap64) And(b2 Bitmap64) {
	for i := 0; i < len(b2) && i < len(*b); i++ {
		(*b)[i] &= b2[i]
	}
}

// Shrink remove zero elements at the end of the map
func (b *Bitmap64) Shrink() {
	shrinkedIndex := len(*b)
	for i := len(*b) - 1; i >= 0; i-- {
		if (*b)[i] != 0 {
			shrinkedIndex = i + 1
			break
		}
	}

	if shrinkedIndex != len(*b) {
		newSlice := make(Bitmap64, shrinkedIndex)
		copy(newSlice, (*b)[:shrinkedIndex])
		*b = newSlice
	}
}

// Clone create a copy of the bitmap
func (b *Bitmap64) Clone() Bitmap64 {
	clone := make(Bitmap64, len(*b))
	copy(clone, *b)

	return clone
}

// Range call the passed callback with all bits set to 1.
// If the callback returns false, the method exits
func (b *Bitmap64) Range(f func(n uint32) bool) {
	for i, block := range *b {
		for block != 0 {
			tz := bits.TrailingZeros64(block)
			bitIndex := uint32(i*64 + tz)

			if !f(bitIndex) {
				return
			}

			block &= block - 1
		}
	}
}

func (b *Bitmap64) String() string {
	var sb strings.Builder

	for i := range *b {
		sb.WriteString(strconv.FormatUint((*b)[i], 10))
		if i != len(*b)-1 {
			sb.WriteString("|")
		}
	}

	return sb.String()
}

func FromString(str string) (Bitmap64, error) {
	if str == "" {
		return Bitmap64{}, nil
	}

	nums := strings.Split(str, "|")
	result := make(Bitmap64, 0, len(nums))
	for _, num := range nums {
		v, err := strconv.ParseUint(num, 10, 64)
		if err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

func (b *Bitmap64) grow(length uint32) {
	if length+1 > uint32(len(*b)) {
		*b = append(*b, make(Bitmap64, length+1-uint32(len(*b)))...)
	}
}
