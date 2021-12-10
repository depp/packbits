package packbits

const (
	maxRun = 128
)

// packOne packs a literal followed by a repeating run, and appends the packed
// data to dst. The number of input bytes packed is returned. If flush is false,
// not all input bytes may be packed, in an effort to merge nearby literals.
func packOne(dst, lit []byte, value byte, repeat int, flush bool) (out []byte, count int) {
	out = dst
	// Maximum literal length is 128, so pack as many 128-length literals as we
	// can.
	for len(lit) >= maxRun {
		out = append(out, maxRun-1)
		out = append(out, lit[:maxRun]...)
		lit = lit[maxRun:]
		count += maxRun
	}
	// Don't emit a run of repeats if it can be merged with the previous
	// literal. This allows two literals to be merged if they are separated by
	// runs of length 2, which saves one byte.
	if !flush && repeat <= 2 && len(lit) > 0 && len(lit) <= maxRun-2 {
		return
	}
	// Check if the run has an odd byte out, and add it to the previous
	// literal if possible. Otherwise, leave it for the next literal.
	var oddbyte bool
	if repeat&127 == 1 {
		repeat--
		if len(lit) > 0 {
			out = append(out, byte(len(lit)))
			out = append(out, lit...)
			out = append(out, value)
			count += len(lit) + 1
			lit = nil
		} else {
			oddbyte = true
		}
	}
	if len(lit) > 0 {
		out = append(out, byte(len(lit)-1))
		out = append(out, lit...)
		count += len(lit)
	}
	if repeat > 0 {
		count += repeat
		for repeat > maxRun {
			out = append(out, 257-maxRun, value)
			repeat -= maxRun
		}
		out = append(out, byte(257-repeat), value)
	}
	if flush && oddbyte {
		out = append(out, 0, value)
		count++
	}
	return
}

// Pack compresses the data using the PackBits compression scheme.
func Pack(data []byte) []byte {
	const maxRun = 128
	var out []byte
	// Loop invariants:
	// - lstart <= rstart <= i
	// - data[start:] has not been packed
	// - most recent runlen bytes are equal to prev
	var prev byte
	var start, runlen int
	for i, b := range data {
		if runlen == 0 {
			// We are not in a run, check if a run starts.
			if b == prev && start < i {
				runlen = 2
			}
		} else if b != prev {
			// We are in the byte after a run.
			var count int
			out, count = packOne(out, data[start:i-runlen], prev, runlen, false)
			start += count
			runlen = 0
		} else {
			// We are in a run which continues.
			runlen++
		}
		prev = b
	}
	out, _ = packOne(out, data[start:len(data)-runlen], prev, runlen, true)
	return out
}
