package bitmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmark_Bitmap16_Set(b *testing.B) {
	var bm Bitmap16
	for i := 0; i < b.N; i++ {
		bm.Set(1)
	}
}

func Test_Bitmap16_Set(t *testing.T) {
	t.Run("must set the specified bit", func(t *testing.T) {
		t.Run("0", func(t *testing.T) {
			var b Bitmap16
			b.Set(0)
			assert.Equal(t, Bitmap16{1}, b)
		})
		t.Run("2", func(t *testing.T) {
			var b Bitmap16
			b.Set(2)
			assert.Equal(t, Bitmap16{4}, b)
		})
		t.Run("0,2", func(t *testing.T) {
			var b Bitmap16
			b.Set(0)
			b.Set(2)
			assert.Equal(t, Bitmap16{5}, b)
		})

		t.Run("15", func(t *testing.T) {
			var b Bitmap16
			b.Set(15)
			assert.Equal(t, Bitmap16{32768}, b)
		})

		t.Run("16", func(t *testing.T) {
			var b Bitmap16
			b.Set(16)
			assert.Equal(t, Bitmap16{0, 1}, b)
		})
	})
	t.Run("must do nothing if the specified bit is already == 1", func(t *testing.T) {
		var b Bitmap16
		b.Set(0)
		b.Set(0)
		assert.Equal(t, Bitmap16{1}, b)
	})
}

func Test_Bitmap16_Remove(t *testing.T) {
	var b Bitmap16

	b.Remove(0)
	assert.Nil(t, b)

	b.Set(0)
	b.Set(1)
	b.Set(100)
	b.Remove(100)
	assert.Equal(t, Bitmap16{3, 0, 0, 0, 0, 0, 0}, b)
	b.Remove(1)
	assert.Equal(t, Bitmap16{1, 0, 0, 0, 0, 0, 0}, b)
	b.Remove(0)
	assert.Equal(t, Bitmap16{0, 0, 0, 0, 0, 0, 0}, b)
}

func Benchmark_Bitmap16_Xor(b *testing.B) {
	var bm Bitmap16
	for i := 0; i < b.N; i++ {
		bm.Xor(1)
	}
}

func Test_Bitmap16_Xor(t *testing.T) {
	t.Run("must invert the specified bit", func(t *testing.T) {
		var b Bitmap16
		b.Xor(0)
		assert.Equal(t, Bitmap16{1}, b)
		b.Xor(1)
		assert.Equal(t, Bitmap16{3}, b)
		b.Xor(0)
		assert.Equal(t, Bitmap16{2}, b)
		b.Xor(1)
		assert.Equal(t, Bitmap16{0}, b)
	})
}

func Test_Bitmap16_IsEmpty(t *testing.T) {
	t.Run("must return true if it's empty", func(t *testing.T) {
		var b Bitmap16
		assert.True(t, b.IsEmpty())
		b.Xor(1)
		b.Xor(1)
		assert.True(t, b.IsEmpty())
	})
	t.Run("must return false if it is not empty", func(t *testing.T) {
		var b Bitmap16
		b.Set(1)
		assert.False(t, b.IsEmpty())
	})
}

func Test_Bitmap16_Has(t *testing.T) {
	var b Bitmap16
	assert.False(t, b.Has(0))
	b.Xor(0)
	assert.True(t, b.Has(0))
	b.Xor(0)
	assert.False(t, b.Has(0))
}

func Benchmark_Bitmap16_CountDiff(b *testing.B) {
	var b1, b2 Bitmap16
	b1.Set(1)
	b1.Set(2)
	b2.Set(31)
	for i := 0; i < b.N; i++ {
		b1.CountDiff(b2)
	}
}

func Test_Bitmap16_CountDiff(t *testing.T) {
	t.Run("must return 0 if bitmaps are equal", func(t *testing.T) {
		var b1, b2 Bitmap16
		b1.Set(1)
		b2.Set(1)
		assert.Equal(t, 0, b1.CountDiff(b2))
	})

	t.Run("must return correct count of different bits", func(t *testing.T) {
		var b1, b2 Bitmap16
		b1.Set(1)
		b1.Set(2)
		b1.Set(64)

		b2.Set(31)
		b2.Set(64)
		b2.Set(650)

		assert.Equal(t, 4, b1.CountDiff(b2))
	})
}

func Test_Bitmap16_Or(t *testing.T) {
	var b1, b2 Bitmap16
	b1.Set(0)
	b1.Set(1)
	b1.Set(100)

	b2.Set(0)
	b2.Set(1)
	b2.Set(2)
	b2.Set(101)
	b2.Set(128)

	b1.Or(b2)

	assert.Equal(t, Bitmap16{7, 0, 0, 0, 0, 0, 48, 0, 1}, b1)
	assert.True(t, b1.Has(0))
	assert.True(t, b1.Has(1))
	assert.True(t, b1.Has(2))
	assert.True(t, b1.Has(100))
	assert.True(t, b1.Has(101))
	assert.True(t, b1.Has(128))
	assert.False(t, b2.Has(100))
}

func Test_Bitmap16_And(t *testing.T) {
	var b1, b2 Bitmap16
	b1.Set(0)
	b1.Set(1)
	b1.Set(100)

	b2.Set(0)
	b2.Set(1)
	b2.Set(2)
	b2.Set(101)
	b2.Set(128)

	b1.And(b2)

	assert.Equal(t, Bitmap16{3, 0, 0, 0, 0, 0, 0}, b1)
	assert.True(t, b1.Has(0))
	assert.True(t, b1.Has(1))
	assert.False(t, b1.Has(2))
	assert.False(t, b1.Has(100))
	assert.False(t, b1.Has(101))
	assert.False(t, b1.Has(128))
	assert.True(t, b2.Has(2))
}

func Test_Bitmap16_Shrink(t *testing.T) {
	var b Bitmap16
	b.Set(1)
	b.Set(100)
	b.Remove(100)
	b.Shrink()
	assert.Equal(t, Bitmap16{2}, b)
}

func Test_Bitmap16_Clone(t *testing.T) {
	var b1 Bitmap16
	b1.Set(0)
	b2 := b1.Clone()

	assert.Equal(t, b1, b2)

	b2.Set(2)
	assert.Equal(t, Bitmap16{5}, b2)
	assert.Equal(t, Bitmap16{1}, b1)
}

func Test_Bitmap16_Range(t *testing.T) {
	var b1 Bitmap16
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

func Benchmark_Bitmap16_String(b *testing.B) {
	bm := Bitmap16{0, 5, 1000}
	for i := 0; i < b.N; i++ {
		_ = bm.String()
	}
}

func Test_Bitmap16_String(t *testing.T) {
	b := Bitmap16{}
	assert.Equal(t, "", b.String())

	b = Bitmap16{0, 5, 100}
	assert.Equal(t, "0|5|100", b.String())
}

func Test_FromString16(t *testing.T) {
	t.Run("must return error if unable to parse the string", func(t *testing.T) {
		_, err := FromString16("qwe")
		assert.Error(t, err)
	})
	t.Run("must parse the string correctly", func(t *testing.T) {
		v, err := FromString16("")
		assert.Nil(t, err)
		assert.Equal(t, Bitmap16{}, v)

		v, err = FromString16("0")
		assert.Nil(t, err)
		assert.Equal(t, Bitmap16{0}, v)

		v, err = FromString16("0|5")
		assert.Nil(t, err)
		assert.Equal(t, Bitmap16{0, 5}, v)
	})
}
