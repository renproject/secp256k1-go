package secp256k1

type FpOld struct {
	n [5]uint64
}

func (r *FpOld) secp256k1_fe_normalize() {
	t0, t1, t2, t3, t4 := r.n[0], r.n[1], r.n[2], r.n[3], r.n[4]

	/* Reduce t4 at the start so there will be at most a single carry from the first pass */
	var m uint64
	x := t4 >> 48
	t4 &= 0x0FFFFFFFFFFFF

	/* The first pass ensures the magnitude is 1, ... */
	t0 += x * 0x1000003D1
	t1 += (t0 >> 52)
	t0 &= 0xFFFFFFFFFFFFF
	t2 += (t1 >> 52)
	t1 &= 0xFFFFFFFFFFFFF
	m = t1
	t3 += (t2 >> 52)
	t2 &= 0xFFFFFFFFFFFFF
	m &= t2
	t4 += (t3 >> 52)
	t3 &= 0xFFFFFFFFFFFFF
	m &= t3

	/* ... except for a possible carry at bit 48 of t4 (i.e. bit 256 of the field element) */
	verifyCheck(t4>>49 == 0)

	/* At most a single final reduction is needed; check if the value is >= the field characteristic */
	b := (t4>>48 != 0) || ((t4 == 0x0FFFFFFFFFFFF) && (m == 0xFFFFFFFFFFFFF) && (t0 >= 0xFFFFEFFFFFC2F))
	if b {
		x = 1
	} else {
		x = 0
	}

	/* Apply the final reduction (for constant-time behaviour, we do it always) */
	t0 += x * 0x1000003D1
	t1 += (t0 >> 52)
	t0 &= 0xFFFFFFFFFFFFF
	t2 += (t1 >> 52)
	t1 &= 0xFFFFFFFFFFFFF
	t3 += (t2 >> 52)
	t2 &= 0xFFFFFFFFFFFFF
	t4 += (t3 >> 52)
	t3 &= 0xFFFFFFFFFFFFF

	/* If t4 didn't carry to bit 48 already, then it should have after any final reduction */
	verifyCheck(t4>>48 == x)

	/* Mask off the possible multiple of 2^256 from the final reduction */
	t4 &= 0x0FFFFFFFFFFFF

	r.n[0] = t0
	r.n[1] = t1
	r.n[2] = t2
	r.n[3] = t3
	r.n[4] = t4
}

func (r *FpOld) secp256k1_fe_normalize_weak() {
	t0, t1, t2, t3, t4 := r.n[0], r.n[1], r.n[2], r.n[3], r.n[4]

	/* Reduce t4 at the start so there will be at most a single carry from the first pass */
	x := t4 >> 48
	t4 &= 0x0FFFFFFFFFFFF

	/* The first pass ensures the magnitude is 1, ... */
	t0 += x * 0x1000003D1
	t1 += (t0 >> 52)
	t0 &= 0xFFFFFFFFFFFFF
	t2 += (t1 >> 52)
	t1 &= 0xFFFFFFFFFFFFF
	t3 += (t2 >> 52)
	t2 &= 0xFFFFFFFFFFFFF
	t4 += (t3 >> 52)
	t3 &= 0xFFFFFFFFFFFFF

	/* ... except for a possible carry at bit 48 of t4 (i.e. bit 256 of the field element) */
	verifyCheck(t4>>49 == 0)

	r.n[0] = t0
	r.n[1] = t1
	r.n[2] = t2
	r.n[3] = t3
	r.n[4] = t4
}

func (r *FpOld) secp256k1_fe_normalize_var() {
	t0, t1, t2, t3, t4 := r.n[0], r.n[1], r.n[2], r.n[3], r.n[4]

	/* Reduce t4 at the start so there will be at most a single carry from the first pass */
	var m uint64
	x := t4 >> 48
	t4 &= 0x0FFFFFFFFFFFF

	/* The first pass ensures the magnitude is 1, ... */
	t0 += x * 0x1000003D1
	t1 += (t0 >> 52)
	t0 &= 0xFFFFFFFFFFFFF
	t2 += (t1 >> 52)
	t1 &= 0xFFFFFFFFFFFFF
	m = t1
	t3 += (t2 >> 52)
	t2 &= 0xFFFFFFFFFFFFF
	m &= t2
	t4 += (t3 >> 52)
	t3 &= 0xFFFFFFFFFFFFF
	m &= t3

	/* ... except for a possible carry at bit 48 of t4 (i.e. bit 256 of the field element) */
	verifyCheck(t4>>49 == 0)

	/* At most a single final reduction is needed; check if the value is >= the field characteristic */
	b := (t4>>48 != 0) || ((t4 == 0x0FFFFFFFFFFFF) && (m == 0xFFFFFFFFFFFFF) && (t0 >= 0xFFFFEFFFFFC2F))

	if b {
		t0 += 0x1000003D1
		t1 += (t0 >> 52)
		t0 &= 0xFFFFFFFFFFFFF
		t2 += (t1 >> 52)
		t1 &= 0xFFFFFFFFFFFFF
		t3 += (t2 >> 52)
		t2 &= 0xFFFFFFFFFFFFF
		t4 += (t3 >> 52)
		t3 &= 0xFFFFFFFFFFFFF

		/* If t4 didn't carry to bit 48 already, then it should have after any final reduction */
		verifyCheck(t4>>48 == x)

		/* Mask off the possible multiple of 2^256 from the final reduction */
		t4 &= 0x0FFFFFFFFFFFF
	}

	r.n[0] = t0
	r.n[1] = t1
	r.n[2] = t2
	r.n[3] = t3
	r.n[4] = t4
}

func (r *FpOld) secp256k1_fe_normalizes_to_zero() bool {
	t0, t1, t2, t3, t4 := r.n[0], r.n[1], r.n[2], r.n[3], r.n[4]

	/* z0 tracks a possible raw value of 0, z1 tracks a possible raw value of P */
	var z0, z1 uint64

	/* Reduce t4 at the start so there will be at most a single carry from the first pass */
	x := t4 >> 48
	t4 &= 0x0FFFFFFFFFFFF

	/* The first pass ensures the magnitude is 1, ... */
	t0 += x * 0x1000003D1
	t1 += (t0 >> 52)
	t0 &= 0xFFFFFFFFFFFFF
	z0 = t0
	z1 = t0 ^ 0x1000003D0
	t2 += (t1 >> 52)
	t1 &= 0xFFFFFFFFFFFFF
	z0 |= t1
	z1 &= t1
	t3 += (t2 >> 52)
	t2 &= 0xFFFFFFFFFFFFF
	z0 |= t2
	z1 &= t2
	t4 += (t3 >> 52)
	t3 &= 0xFFFFFFFFFFFFF
	z0 |= t3
	z1 &= t3
	z0 |= t4
	z1 &= t4 ^ 0xF000000000000

	/* ... except for a possible carry at bit 48 of t4 (i.e. bit 256 of the field element) */
	verifyCheck(t4>>49 == 0)

	return (z0 == 0) || (z1 == 0xFFFFFFFFFFFFF)
}

func (r *FpOld) secp256k1_fe_normalizes_to_zero_var() bool {
	var t0, t1, t2, t3, t4 uint64
	var z0, z1 uint64
	var x uint64

	t0 = r.n[0]
	t4 = r.n[4]

	/* Reduce t4 at the start so there will be at most a single carry from the first pass */
	x = t4 >> 48

	/* The first pass ensures the magnitude is 1, ... */
	t0 += x * 0x1000003D1

	/* z0 tracks a possible raw value of 0, z1 tracks a possible raw value of P */
	z0 = t0 & 0xFFFFFFFFFFFFF
	z1 = z0 ^ 0x1000003D0

	/* Fast return path should catch the majority of cases */
	if (z0 != 0) && (z1 != 0xFFFFFFFFFFFFF) {
		return false
	}

	t1 = r.n[1]
	t2 = r.n[2]
	t3 = r.n[3]

	t4 &= 0x0FFFFFFFFFFFF

	t1 += (t0 >> 52)
	t2 += (t1 >> 52)
	t1 &= 0xFFFFFFFFFFFFF
	z0 |= t1
	z1 &= t1
	t3 += (t2 >> 52)
	t2 &= 0xFFFFFFFFFFFFF
	z0 |= t2
	z1 &= t2
	t4 += (t3 >> 52)
	t3 &= 0xFFFFFFFFFFFFF
	z0 |= t3
	z1 &= t3
	z0 |= t4
	z1 &= t4 ^ 0xF000000000000

	/* ... except for a possible carry at bit 48 of t4 (i.e. bit 256 of the field element) */
	verifyCheck(t4>>49 == 0)

	return (z0 == 0) || (z1 == 0xFFFFFFFFFFFFF)
}

// NOTE: Was inlined.
func (r *FpOld) secp256k1_fe_set_int(a int) {
	r.n[0] = uint64(a)
	r.n[1] = 0
	r.n[2] = 0
	r.n[3] = 0
	r.n[4] = 0
}

// NOTE: Was inlined.
func (r *FpOld) secp256k1_fe_is_zero() bool {
	t := r.n
	return (t[0] | t[1] | t[2] | t[3] | t[4]) == 0
}

// NOTE: Was inlined.
func (r *FpOld) secp256k1_fe_is_odd() bool {
	return r.n[0]&1 == 1
}

// NOTE: Was inlined.
func (r *FpOld) secp256k1_fe_clear() {
	var i int
	for i = 0; i < 5; i++ {
		r.n[i] = 0
	}
}

func (r *FpOld) secp256k1_fe_cmp_var(other *FpOld) int {
	var i int
	for i = 4; i >= 0; i-- {
		if r.n[i] > other.n[i] {
			return 1
		}
		if r.n[i] < other.n[i] {
			return -1
		}
	}
	return 0
}

func (r *FpOld) secp256k1_fe_set_b32(a [32]byte) int {
	r.n[0] = uint64(a[31]) | (uint64(a[30]) << 8) | (uint64(a[29]) << 16) | (uint64(a[28]) << 24) | (uint64(a[27]) << 32) | (uint64(a[26]) << 40) | (uint64((a[25])&0xF) << 48)
	r.n[1] = uint64(((a[25] >> 4) & 0xF)) | (uint64(a[24]) << 4) | (uint64(a[23]) << 12) | (uint64(a[22]) << 20) | (uint64(a[21]) << 28) | (uint64(a[20]) << 36) | (uint64(a[19]) << 44)
	r.n[2] = uint64(a[18]) | (uint64(a[17]) << 8) | (uint64(a[16]) << 16) | (uint64(a[15]) << 24) | (uint64(a[14]) << 32) | (uint64(a[13]) << 40) | (uint64((a[12])&0xF) << 48)
	r.n[3] = uint64(((a[12] >> 4) & 0xF)) | (uint64(a[11]) << 4) | (uint64(a[10]) << 12) | (uint64(a[9]) << 20) | (uint64(a[8]) << 28) | (uint64(a[7]) << 36) | (uint64(a[6]) << 44)
	r.n[4] = uint64(a[5]) | (uint64(a[4]) << 8) | (uint64(a[3]) << 16) | (uint64(a[2]) << 24) | (uint64(a[1]) << 32) | (uint64(a[0]) << 40)
	if r.n[4] == 0x0FFFFFFFFFFFF && (r.n[3]&r.n[2]&r.n[1]) == 0xFFFFFFFFFFFFF && r.n[0] >= 0xFFFFEFFFFFC2F {
		return 0
	}
	return 1
}

/** Convert a field element to a 32-byte big endian value. Requires the input to be normalized */
func (r *FpOld) secp256k1_fe_get_b32(b []byte) {
	b[0] = byte(r.n[4]>>40) & 0xFF
	b[1] = byte(r.n[4]>>32) & 0xFF
	b[2] = byte(r.n[4]>>24) & 0xFF
	b[3] = byte(r.n[4]>>16) & 0xFF
	b[4] = byte(r.n[4]>>8) & 0xFF
	b[5] = byte(r.n[4]) & 0xFF
	b[6] = byte(r.n[3]>>44) & 0xFF
	b[7] = byte(r.n[3]>>36) & 0xFF
	b[8] = byte(r.n[3]>>28) & 0xFF
	b[9] = byte(r.n[3]>>20) & 0xFF
	b[10] = byte(r.n[3]>>12) & 0xFF
	b[11] = byte(r.n[3]>>4) & 0xFF
	b[12] = byte((r.n[2]>>48)&0xF) | byte((r.n[3]&0xF)<<4)
	b[13] = byte(r.n[2]>>40) & 0xFF
	b[14] = byte(r.n[2]>>32) & 0xFF
	b[15] = byte(r.n[2]>>24) & 0xFF
	b[16] = byte(r.n[2]>>16) & 0xFF
	b[17] = byte(r.n[2]>>8) & 0xFF
	b[18] = byte(r.n[2]) & 0xFF
	b[19] = byte(r.n[1]>>44) & 0xFF
	b[20] = byte(r.n[1]>>36) & 0xFF
	b[21] = byte(r.n[1]>>28) & 0xFF
	b[22] = byte(r.n[1]>>20) & 0xFF
	b[23] = byte(r.n[1]>>12) & 0xFF
	b[24] = byte(r.n[1]>>4) & 0xFF
	b[25] = byte((r.n[0]>>48)&0xF) | byte((r.n[1]&0xF)<<4)
	b[26] = byte(r.n[0]>>40) & 0xFF
	b[27] = byte(r.n[0]>>32) & 0xFF
	b[28] = byte(r.n[0]>>24) & 0xFF
	b[29] = byte(r.n[0]>>16) & 0xFF
	b[30] = byte(r.n[0]>>8) & 0xFF
	b[31] = byte(r.n[0]) & 0xFF
}

// NOTE: Was inlined.
func (r *FpOld) secp256k1_fe_negate(a *FpOld, m uint64) {
	r.n[0] = 0xFFFFEFFFFFC2F*2*(m+1) - a.n[0]
	r.n[1] = 0xFFFFFFFFFFFFF*2*(m+1) - a.n[1]
	r.n[2] = 0xFFFFFFFFFFFFF*2*(m+1) - a.n[2]
	r.n[3] = 0xFFFFFFFFFFFFF*2*(m+1) - a.n[3]
	r.n[4] = 0x0FFFFFFFFFFFF*2*(m+1) - a.n[4]
}

// NOTE: Was inlined.
func (r *FpOld) secp256k1_fe_mul_int(a uint64) {
	r.n[0] *= a
	r.n[1] *= a
	r.n[2] *= a
	r.n[3] *= a
	r.n[4] *= a
}

// NOTE: Was inlined.
func (r *FpOld) secp256k1_fe_add(a *FpOld) {
	r.n[0] += a.n[0]
	r.n[1] += a.n[1]
	r.n[2] += a.n[2]
	r.n[3] += a.n[3]
	r.n[4] += a.n[4]
}

func (r *FpOld) secp256k1_fe_mul(a, b *FpOld) {
	// secp256k1_fe_mul_inner(r.n, a.n, b.n)
}

func (r *FpOld) secp256k1_fe_sqr(a *FpOld) {
	// secp256k1_fe_sqr_inner(r.n, a.n)
}
