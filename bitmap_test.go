package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Benchmark_Bitmap_Set(b *testing.B) {
	var bm Bitmap
	for i := 0; i < b.N; i++ {
		bm.Set(1)
	}
}

func Test_Bitmap_Set(t *testing.T) {
	t.Run("must set the specified bit", func(t *testing.T) {
		t.Run("0", func(t *testing.T) {
			var b Bitmap
			b.Set(0)
			assert.Equal(t, Bitmap{1}, b)
		})
		t.Run("2", func(t *testing.T) {
			var b Bitmap
			b.Set(2)
			assert.Equal(t, Bitmap{4}, b)
		})
		t.Run("0,2", func(t *testing.T) {
			var b Bitmap
			b.Set(0)
			b.Set(2)
			assert.Equal(t, Bitmap{5}, b)
		})

		t.Run("63", func(t *testing.T) {
			var b Bitmap
			b.Set(63)
			assert.Equal(t, Bitmap{9223372036854775808}, b)
		})

		t.Run("64", func(t *testing.T) {
			var b Bitmap
			b.Set(64)
			assert.Equal(t, Bitmap{0, 1}, b)
		})
	})
	t.Run("must do nothing if the specified bit is already == 1", func(t *testing.T) {
		var b Bitmap
		b.Set(0)
		b.Set(0)
		assert.Equal(t, Bitmap{1}, b)
	})
}

func Test_Bitmap_Remove(t *testing.T) {
	var b Bitmap

	b.Remove(0)
	assert.Nil(t, b)

	b.Set(0)
	b.Set(100)
	b.Remove(100)
	assert.Equal(t, Bitmap{1, 0}, b)
}

func Benchmark_Bitmap_Xor(b *testing.B) {
	var bm Bitmap
	for i := 0; i < b.N; i++ {
		bm.Xor(1)
	}
}

func Test_Bitmap_Xor(t *testing.T) {
	t.Run("must invert the specified bit", func(t *testing.T) {
		var b Bitmap
		b.Xor(0)
		assert.Equal(t, Bitmap{1}, b)
		b.Xor(0)
		assert.Equal(t, Bitmap{0}, b)
	})
}

func Test_Bitmap_IsEmpty(t *testing.T) {
	t.Run("must return true if it's empty", func(t *testing.T) {
		var b Bitmap
		assert.True(t, b.IsEmpty())
		b.Xor(1)
		b.Xor(1)
		assert.True(t, b.IsEmpty())
	})
	t.Run("must return false if it is not empty", func(t *testing.T) {
		var b Bitmap
		b.Set(1)
		assert.False(t, b.IsEmpty())
	})
}

func Test_Bitmap_Has(t *testing.T) {
	var b Bitmap
	assert.False(t, b.Has(0))
	b.Xor(0)
	assert.True(t, b.Has(0))
	b.Xor(0)
	assert.False(t, b.Has(0))
}

func Benchmark_Bitmap_CountDiff(b *testing.B) {
	var b1, b2 Bitmap
	b1.Set(1)
	b1.Set(2)
	b2.Set(31)
	for i := 0; i < b.N; i++ {
		b1.CountDiff(b2)
	}
}

func Test_Bitmap_CountDiff(t *testing.T) {
	t.Run("must return 0 if bitmaps are equal", func(t *testing.T) {
		var b1, b2 Bitmap
		b1.Set(1)
		b2.Set(1)
		assert.Equal(t, 0, b1.CountDiff(b2))
	})

	t.Run("must return correct count of different bits", func(t *testing.T) {
		var b1, b2 Bitmap
		b1.Set(1)
		b1.Set(2)
		b1.Set(64)

		b2.Set(31)
		b2.Set(64)
		b2.Set(650)

		assert.Equal(t, 4, b1.CountDiff(b2))
	})
}

func Test_Bitmap_Clone(t *testing.T) {
	var b1 Bitmap
	b1.Set(0)
	b2 := b1.Clone()

	assert.Equal(t, b1, b2)

	b2.Set(2)
	assert.Equal(t, Bitmap{5}, b2)
	assert.Equal(t, Bitmap{1}, b1)
}

func Test_Bitmap_Range(t *testing.T) {
	var b1 Bitmap
	b1.Set(0)
	b1.Set(1)
	b1.Set(2)
	b1.Set(1000)
	b1.Set(10000)

	var items []uint32
	b1.Range(func(n uint32) bool {
		items = append(items, n)
		if n == 1000 {
			return false
		}

		return true
	})

	require.Equal(t, []uint32{0, 1, 2, 1000}, items)
}
