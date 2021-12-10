# PackBits Compression in Go

## Overview

This package provides an implementation of the PackBits compression scheme for Go. The PackBits compression scheme is a simple lossless run-length encoding scheme. It originally appeared as the compression scheme used by MacPaint on Macintosh computers in 1984.

### References

- [Apple Technical Note TN1023](https://web.archive.org/web/20080705155158/http://developer.apple.com/technotes/tn/tn1023.html)

- [Wikipedia: PackBits](https://en.wikipedia.org/wiki/PackBits)

## Example of Usage

```go
// Compress raw data.
data := []byte{
    0xAA, 0xAA, 0xAA, 0x80, 0x00, 0x2A, 0xAA, 0xAA,
    0xAA, 0xAA, 0x80, 0x00, 0x2A, 0x22, 0xAA, 0xAA,
    0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA, 0xAA,
}
packed := packbits.Pack(data)

// Decompress packed data.
unpacked, err := packbits.Unpack(packed)
```

## Efficiency

There is more than one way to encode a given block of raw data. I believe that this library will always choose one of the smallest possible encodings. In certain circumstances, this means merging adjancent blocks of data into a single block.

## License

PackBits is distributed under the terms of the terms of the MIT license. See LICENSE.txt for details.
