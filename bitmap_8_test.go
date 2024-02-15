package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmark_Bitmap8_Set(b *testing.B) {
	var bm Bitmap8
	for i := 0; i < b.N; i++ {
		bm.Set(1)
	}
}

func Test_Bitmap8_Set(t *testing.T) {
	t.Run("must set the specified bit", func(t *testing.T) {
		t.Run("0", func(t *testing.T) {
			var b Bitmap8
			b.Set(0)
			assert.Equal(t, Bitmap8{1}, b)
		})
		t.Run("2", func(t *testing.T) {
			var b Bitmap8
			b.Set(2)
			assert.Equal(t, Bitmap8{4}, b)
		})
		t.Run("0,2", func(t *testing.T) {
			var b Bitmap8
			b.Set(0)
			b.Set(2)
			assert.Equal(t, Bitmap8{5}, b)
		})

		t.Run("7", func(t *testing.T) {
			var b Bitmap8
			b.Set(7)
			assert.Equal(t, Bitmap8{0x80}, b)
		})

		t.Run("8", func(t *testing.T) {
			var b Bitmap8
			b.Set(8)
			assert.Equal(t, Bitmap8{0, 1}, b)
		})
	})
	t.Run("must do nothing if the specified bit is already == 1", func(t *testing.T) {
		var b Bitmap8
		b.Set(0)
		b.Set(0)
		assert.Equal(t, Bitmap8{1}, b)
	})
}

func Test_Bitmap8_Remove(t *testing.T) {
	var b Bitmap8

	b.Remove(0)
	assert.Nil(t, b)

	b.Set(0)
	b.Set(1)
	b.Set(10)
	b.Remove(10)
	assert.Equal(t, Bitmap8{3, 0}, b)
	b.Remove(1)
	assert.Equal(t, Bitmap8{1, 0}, b)
	b.Remove(0)
	assert.Equal(t, Bitmap8{0, 0}, b)
}

func Benchmark_Bitmap8_Xor(b *testing.B) {
	var bm Bitmap8
	for i := 0; i < b.N; i++ {
		bm.Xor(1)
	}
}

func Test_Bitmap8_Xor(t *testing.T) {
	t.Run("must invert the specified bit", func(t *testing.T) {
		var b Bitmap8
		b.Xor(0)
		assert.Equal(t, Bitmap8{1}, b)
		b.Xor(1)
		assert.Equal(t, Bitmap8{3}, b)
		b.Xor(0)
		assert.Equal(t, Bitmap8{2}, b)
		b.Xor(1)
		assert.Equal(t, Bitmap8{0}, b)
	})
}

func Test_Bitmap8_IsEmpty(t *testing.T) {
	t.Run("must return true if it's empty", func(t *testing.T) {
		var b Bitmap8
		assert.True(t, b.IsEmpty())
		b.Xor(1)
		b.Xor(1)
		assert.True(t, b.IsEmpty())
	})
	t.Run("must return false if it is not empty", func(t *testing.T) {
		var b Bitmap8
		b.Set(1)
		assert.False(t, b.IsEmpty())
	})
}

func Test_Bitmap8_Has(t *testing.T) {
	var b Bitmap8
	assert.False(t, b.Has(0))
	b.Xor(0)
	assert.True(t, b.Has(0))
	b.Xor(0)
	assert.False(t, b.Has(0))
}

func Benchmark_Bitmap8_CountDiff(b *testing.B) {
	var b1, b2 Bitmap8
	b1.Set(1)
	b1.Set(2)
	b2.Set(31)
	for i := 0; i < b.N; i++ {
		b1.CountDiff(b2)
	}
}

func Test_Bitmap8_CountDiff(t *testing.T) {
	t.Run("must return 0 if bitmaps are equal", func(t *testing.T) {
		var b1, b2 Bitmap8
		b1.Set(1)
		b2.Set(1)
		assert.Equal(t, 0, b1.CountDiff(b2))
	})

	t.Run("must return correct count of different bits", func(t *testing.T) {
		var b1, b2 Bitmap8
		b1.Set(1)
		b1.Set(2)
		b1.Set(64)

		b2.Set(31)
		b2.Set(64)
		b2.Set(650)

		assert.Equal(t, 4, b1.CountDiff(b2))
	})
}

func Test_Bitmap8_Or(t *testing.T) {
	var b1, b2 Bitmap8
	b1.Set(0)
	b1.Set(1)
	b1.Set(10)

	b2.Set(0)
	b2.Set(1)
	b2.Set(2)
	b2.Set(11)
	b2.Set(16)

	b1.Or(b2)

	assert.Equal(t, Bitmap8{0x7, 0xc, 0x1}, b1)
	assert.True(t, b1.Has(0))
	assert.True(t, b1.Has(1))
	assert.True(t, b1.Has(2))
	assert.True(t, b1.Has(10))
	assert.True(t, b1.Has(11))
	assert.True(t, b1.Has(16))
	assert.False(t, b2.Has(10))
}

func Test_Bitmap8_And(t *testing.T) {
	var b1, b2 Bitmap8
	b1.Set(0)
	b1.Set(1)
	b1.Set(10)

	b2.Set(0)
	b2.Set(1)
	b2.Set(2)
	b2.Set(11)
	b2.Set(16)

	b1.And(b2)

	assert.Equal(t, Bitmap8{3, 0}, b1)
	assert.True(t, b1.Has(0))
	assert.True(t, b1.Has(1))
	assert.False(t, b1.Has(2))
	assert.False(t, b1.Has(10))
	assert.False(t, b1.Has(11))
	assert.False(t, b1.Has(16))
	assert.True(t, b2.Has(2))
}

func Test_Bitmap8_Shrink(t *testing.T) {
	var b Bitmap8
	b.Set(1)
	b.Set(100)
	b.Remove(100)
	b.Shrink()
	assert.Equal(t, Bitmap8{2}, b)
}

func Test_Bitmap8_Clone(t *testing.T) {
	var b1 Bitmap8
	b1.Set(0)
	b2 := b1.Clone()

	assert.Equal(t, b1, b2)

	b2.Set(2)
	assert.Equal(t, Bitmap8{5}, b2)
	assert.Equal(t, Bitmap8{1}, b1)
}

func Test_Bitmap8_Range(t *testing.T) {
	var b1 Bitmap8
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

	assert.Equal(t, []uint32{0, 1, 2, 1000}, items)
}

func Benchmark_Bitmap8_String(b *testing.B) {
	bm := Bitmap8{0, 5, 100}
	for i := 0; i < b.N; i++ {
		_ = bm.String()
	}
}

func Test_Bitmap8_String(t *testing.T) {
	b := Bitmap8{}
	assert.Equal(t, "", b.String())

	b = Bitmap8{0, 5, 100}
	assert.Equal(t, "0|5|100", b.String())
}

func Test_FromString8(t *testing.T) {
	t.Run("must return error if unable to parse the string", func(t *testing.T) {
		_, err := FromString8("qwe")
		assert.Error(t, err)
	})
	t.Run("must parse the string correctly", func(t *testing.T) {
		v, err := FromString8("")
		assert.Nil(t, err)
		assert.Equal(t, Bitmap8{}, v)

		v, err = FromString8("0")
		assert.Nil(t, err)
		assert.Equal(t, Bitmap8{0}, v)

		v, err = FromString8("0|5")
		assert.Nil(t, err)
		assert.Equal(t, Bitmap8{0, 5}, v)
	})
}
