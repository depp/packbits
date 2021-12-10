// Package packbits implements the PackBits lossless data compression scheme, as
// used on old Macintosh computers.
//
// See Apple Technical Note TN1023.
package packbits

import (
	"errors"
)

// ErrInvalidData indicates that the PackBits data is invalid and cannot be
// decompressed.
var ErrInvalidData = errors.New("invalid PackBits data")

// Unpack decompresses the data using the PackBits compression scheme.
func Unpack(data []byte) ([]byte, error) {
	var r []byte
	for len(data) > 0 {
		c := data[0]
		data = data[1:]
		if c <= 127 {
			n := int(c) + 1
			if len(data) < n {
				return nil, ErrInvalidData
			}
			r = append(r, data[:n]...)
			data = data[n:]
		} else if c != 128 {
			if len(data) == 0 {
				return nil, ErrInvalidData
			}
			n := 257 - int(c)
			d := data[0]
			data = data[1:]
			for i := 0; i < n; i++ {
				r = append(r, d)
			}
		}
		// Control byte == 128 -> no output, according to TN1024:
		//
		// "PackBits never generates the value -128 ($80) as a flag-counter
		// byte, but a few PackBits-like routines that are built into some
		// applications do. UnpackBits handles this situation by skipping any
		// flag-counter byte with this value and interpreting the next byte as
		// the next flag-counter byte. If you're writing your own
		// UnpackBits-like routine, make sure it handles this situation in the
		// same way."
	}
	return r, nil
}
