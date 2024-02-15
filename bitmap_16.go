package bitmap

import (
	"math/bits"
	"strconv"
	"strings"
)

type Bitmap16 []uint16

// Set set n-th bit to 1
func (b *Bitmap16) Set(n uint32) {
	block, bit := n>>4, n%16
	b.grow(block)
	(*b)[block] |= (1 << bit)
}

// Remove set n-th bit to 0
func (b *Bitmap16) Remove(n uint32) {
	block, bit := n>>4, n%16
	if uint32(len(*b)) <= block {
		return
	}
	(*b)[block] &= ^(1 << bit)
}

// Xor invert n-th bit
func (b *Bitmap16) Xor(n uint32) {
	block, val := n>>4, n%16
	b.grow(block)
	(*b)[block] ^= (1 << val)
}

// IsEmpty check if the bitmap has any bit set to 1
func (b *Bitmap16) IsEmpty() bool {
	for i := range *b {
		if (*b)[i] > 0 {
			return false
		}
	}

	return true
}

// Has check if n-th bit is set to 1
func (b *Bitmap16) Has(n uint32) bool {
	block, val := n>>4, n%16
	if uint32(len(*b)) <= block {
		return false
	}

	return (*b)[block]&(1<<val) > 0
}

// CountDiff count different bits in two bitmaps
func (b *Bitmap16) CountDiff(b2 Bitmap16) int {
	diff := 0
	max := len(*b)
	if len(b2) > max {
		max = len(b2)
	}

	for i := 0; i < max; i++ {
		if len(b2) <= i {
			diff += bits.OnesCount16((*b)[i])
			continue
		}
		if len(*b) <= i {
			diff += bits.OnesCount16((b2)[i])
			continue
		}

		diff += bits.OnesCount16((*b)[i] ^ b2[i])
	}

	return diff
}

// Or in-place OR operation with another bitmap
func (b *Bitmap16) Or(b2 Bitmap16) {
	b.grow(uint32(len(b2) - 1))
	for i := 0; i < len(b2); i++ {
		if b2[i] == 0 {
			continue
		}
		(*b)[i] |= b2[i]
	}
}

// And in-place And operation with another bitmap
func (b *Bitmap16) And(b2 Bitmap16) {
	for i := 0; i < len(b2) && i < len(*b); i++ {
		(*b)[i] &= b2[i]
	}
}

// Shrink remove zero elements at the end of the map
func (b *Bitmap16) Shrink() {
	shrinkedIndex := len(*b)
	for i := len(*b) - 1; i >= 0; i-- {
		if (*b)[i] != 0 {
			shrinkedIndex = i + 1
			break
		}
	}

	if shrinkedIndex != len(*b) {
		newSlice := make(Bitmap16, shrinkedIndex)
		copy(newSlice, (*b)[:shrinkedIndex])
		*b = newSlice
	}
}

// Clone create a copy of the bitmap
func (b *Bitmap16) Clone() Bitmap16 {
	clone := make(Bitmap16, len(*b))
	copy(clone, *b)

	return clone
}

// Range call the passed callback with all bits set to 1.
// If the callback returns false, the method exits
func (b *Bitmap16) Range(f func(n uint32) bool) {
	for i, block := range *b {
		for block != 0 {
			tz := bits.TrailingZeros16(block)
			bitIndex := uint32(i*16 + tz)

			if !f(bitIndex) {
				return
			}

			block &= block - 1
		}
	}
}

func (b *Bitmap16) String() string {
	var sb strings.Builder

	for i := range *b {
		sb.WriteString(strconv.FormatUint(uint64((*b)[i]), 10))
		if i != len(*b)-1 {
			sb.WriteString("|")
		}
	}

	return sb.String()
}

func FromString16(str string) (Bitmap16, error) {
	if str == "" {
		return Bitmap16{}, nil
	}

	nums := strings.Split(str, "|")
	result := make(Bitmap16, 0, len(nums))
	for _, num := range nums {
		v, err := strconv.ParseUint(num, 10, 16)
		if err != nil {
			return nil, err
		}
		result = append(result, uint16(v))
	}
	return result, nil
}

func (b *Bitmap16) grow(length uint32) {
	if length+1 > uint32(len(*b)) {
		*b = append(*b, make(Bitmap16, length+1-uint32(len(*b)))...)
	}
}
