package packbits

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func hexdata(s string) []byte {
	var r []byte
	for len(s) != 0 {
		var n string
		switch strings.IndexByte(s, ' ') {
		case -1:
			n = s
			s = ""
		case 2:
			n = s[:2]
			s = s[3:]
		default:
			panic("invalid hexdata")
		}
		x, err := strconv.ParseUint(n, 16, 8)
		if err != nil {
			panic("invalid hexdata: " + err.Error())
		}
		r = append(r, byte(x))
	}
	return r
}

var hexdigit = [16]byte{
	'0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}

func pbytes(d []byte) string {
	if len(d) == 0 {
		return "[]"
	}
	var r = make([]byte, len(d)*3+1)
	for i, b := range d {
		r[i*3] = ' '
		r[i*3+1] = hexdigit[b>>4]
		r[i*3+2] = hexdigit[b&15]
	}
	r[0] = '['
	r[len(r)-1] = ']'
	return string(r)
}

func TestUnpack(t *testing.T) {
	type tcase struct {
		packed   string
		unpacked string
	}
	cases := []tcase{
		{
			// Example from tech note.
			packed: "FE AA 02 80 00 2A FD AA 03 80 00 2A 22 F7 AA",

			unpacked: "AA AA AA 80 00 2A AA AA AA AA 80 00 2A 22 AA AA AA AA AA AA AA AA AA AA",
		},
		{
			packed:   "80",
			unpacked: "",
		},
		{}, // Empty is ok
	}
	for i, c := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			packed := hexdata(c.packed)
			expect := hexdata(c.unpacked)
			out, err := Unpack(packed)
			if err != nil {
				t.Error("Error:", err)
			} else if !bytes.Equal(out, expect) {
				t.Error("Incorrect ouptut")
				t.Logf("Got:    %s", pbytes(out))
				t.Logf("Expect: %s", pbytes(expect))
			}
		})
	}
}

func TestPack(t *testing.T) {
	type tcase struct {
		name string
		data []byte
	}
	cases := []tcase{
		{
			// Example from tech note 1023.
			"tn",
			hexdata("AA AA AA 80 00 2A AA AA AA AA 80 00 2A 22 AA AA AA AA AA AA AA AA AA AA"),
		},
		{
			"empty",
			nil,
		},
	}
	// Test runs surrounded by non-runs.
	for _, n := range []int{2, 3, 128, 129, 600} {
		// Create a string [1 2], [3]*n, [4 5]*2.
		x := make([]byte, n+4)
		x[0] = 1
		x[1] = 2
		for j := 2; j < n+2; j++ {
			x[j] = 3
		}
		x[n+2] = 4
		x[n+3] = 5
		name := fmt.Sprintf("run%d", n)
		cases = append(
			cases,
			tcase{name, x[2 : n+2]},
			tcase{name + "_suffix", x[2:]},
			tcase{"prefix_" + name, x[:n+2]},
			tcase{"prefix_" + name + "_suffix", x},
		)
	}
	for _, c := range cases {
		in := c.data
		t.Run(c.name, func(t *testing.T) {
			out := Pack(in)
			unpacked, err := Unpack(out)
			if err != nil {
				t.Error("Error:", err)
				t.Logf("Input:    %s", pbytes(in))
				t.Logf("Packed:   %s", pbytes(out))
			} else if !bytes.Equal(in, unpacked) {
				t.Error("Packed data does not match")
				t.Logf("Input:    %s", pbytes(in))
				t.Logf("Packed:   %s", pbytes(out))
				t.Logf("Unpacked: %s", pbytes(unpacked))
			}
		})
	}
}
