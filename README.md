# Bitmap

Yet another bitmap implementation written in go.

## Installation

```
$ go get -v github.com/f1monkey/bitmap
```

## Usage

```go
func main() {
    var b bitmap.Bitmap

    b.IsEmpty() // true

    b.Set(0)
    b.Has(0) // true
    b.Remove(0)
    b.Has(0) // false

    b.Xor(1000)
    b.Has(1000) // true
    b.Xor(1000)
    b.Has(1000) // false

    b2 := b.Clone() // copy of "b"

    b2.Range(func(n uint32) bool {
        fmt.PrintLn(n)
        return true
    })

    // to string, from string
    var b3 bitmap.Bitmap
    b3.Set(1)
    b3.Set(100)
    b3.String() // "2|68719476736"
    b4, err := bitmap.FromString("2|68719476736")

    // Bitmap32 is backed by []uint32 slice
    // Everything else is all the same
    var b32 bitmap.Bitmap32
}
```