package secp256k1

/*
#include <stdlib.h>

#include "c-field/field_impl.h"

void mulqq(size_t a, size_t b, size_t* dst) {
    __int128 result = (__int128)a * (__int128)b;
	dst[0] = (size_t)result;
	dst[1] = result >> 64;
}
*/
import "C"
import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	"unsafe"
)

const r0 uint64 = 0xda1732fc9bebf
const r1 uint64 = 0x1950b75fc4402
const r2 uint64 = 0x1455123

func mulQQ(x, y uint64) (uint64, uint64) {
	var z [2]uint64
	C.mulqq(C.uint64_t(x), C.uint64_t(y), (*C.uint64_t)(&z[0]))
	return z[0], z[1]
}

type Secp256k1P struct {
	inner C.struct_secp256k1_fe
}

func NewSecp256k1P(a uint64) Secp256k1P {
	inner := C.struct_secp256k1_fe{}
	C.secp256k1_fe_set_int(&inner, C.uint64_t(a))
	return Secp256k1P{inner}
}

func (r *Secp256k1P) Set(a *Secp256k1P) {
	r.inner = a.inner
}

func (r *Secp256k1P) BigInt() *big.Int {
	return big.NewInt(0).SetBytes(r.GetB32())
}

func (r *Secp256k1P) Normalize() {
	C.secp256k1_fe_normalize(&r.inner)
}

func (r *Secp256k1P) NormalizeWeak() {
	C.secp256k1_fe_normalize_weak(&r.inner)
}

func (r *Secp256k1P) NormalizeVar() {
	C.secp256k1_fe_normalize_var(&r.inner)
}

func (r *Secp256k1P) NormalizesToZero() bool {
	return C.secp256k1_fe_normalizes_to_zero(&r.inner) != 0
}

func (r *Secp256k1P) NormalizesToZeroVar() bool {
	return C.secp256k1_fe_normalizes_to_zero_var(&r.inner) != 0
}

func (r *Secp256k1P) SetUint64(a uint64) {
	C.secp256k1_fe_set_int(&r.inner, C.uint64_t(a))
}

func (r *Secp256k1P) Clear() {
	C.secp256k1_fe_clear(&r.inner)
}

func (r *Secp256k1P) IsZero() bool {
	return C.secp256k1_fe_is_zero(&r.inner) != 0
}

func (r *Secp256k1P) IsOdd() bool {
	return C.secp256k1_fe_is_odd(&r.inner) != 0
}

func (r *Secp256k1P) Equal(a *Secp256k1P) bool {
	return C.secp256k1_fe_equal(&r.inner, &a.inner) != 0
}

func (r *Secp256k1P) EqualVar(a *Secp256k1P) bool {
	return C.secp256k1_fe_equal_var(&r.inner, &a.inner) != 0
}

func (r *Secp256k1P) CmpVar(a *Secp256k1P) int {
	return int(C.secp256k1_fe_cmp_var(&r.inner, &a.inner))
}

func (r *Secp256k1P) SetB32(b []byte) {
	cBytes := (*C.uchar)(C.CBytes(b))
	C.secp256k1_fe_set_b32(&r.inner, cBytes)
	C.free(unsafe.Pointer(cBytes))
}

func (r *Secp256k1P) GetB32() []byte {
	cBytes := C.malloc(C.sizeof_char * 32)
	C.secp256k1_fe_get_b32((*C.uchar)(cBytes), &r.inner)
	b := C.GoBytes(cBytes, 32)
	C.free(unsafe.Pointer(cBytes))
	return b
}

func (r *Secp256k1P) Negate(a *Secp256k1P, m int) {
	C.secp256k1_fe_negate(&r.inner, &a.inner, C.int(m))
}

func (r *Secp256k1P) MulInt(a int) {
	C.secp256k1_fe_mul_int(&r.inner, C.int(a))
}

func (r *Secp256k1P) Add(a *Secp256k1P) {
	C.secp256k1_fe_add(&r.inner, &a.inner)
}

func (r *Secp256k1P) Mul(a, b *Secp256k1P) {
	C.secp256k1_fe_mul(&r.inner, &a.inner, &b.inner)
}

func (r *Secp256k1P) Sqr(a *Secp256k1P) {
	C.secp256k1_fe_sqr(&r.inner, &a.inner)
}

func (r *Secp256k1P) Sqrt(a *Secp256k1P) bool {
	return C.secp256k1_fe_sqrt(&r.inner, &a.inner) != 0
}

func (r *Secp256k1P) IsQuadVar() bool {
	return C.secp256k1_fe_is_quad_var(&r.inner) != 0
}

func (r *Secp256k1P) Inv(a *Secp256k1P) {
	C.secp256k1_fe_inv(&r.inner, &a.inner)
}

func (r *Secp256k1P) InvVar(a *Secp256k1P) {
	C.secp256k1_fe_inv_var(&r.inner, &a.inner)
}

func (r *Secp256k1P) Cmov(a *Secp256k1P, flag bool) {
	if flag {
		C.secp256k1_fe_cmov(&r.inner, &a.inner, 1)
	} else {
		C.secp256k1_fe_cmov(&r.inner, &a.inner, 0)
	}
}

// Secp256k1N represents an element of the field defined by the prime that is
// the order of the elliptic curve group Secp256k1.
type Secp256k1N struct {
	limbs [5]uint64
}

func NewSecp256k1N(x uint64) Secp256k1N {
	return Secp256k1N{limbs: [5]uint64{x, 0, 0, 0, 0}}
}

// ZeroSecp256k1N returns the zero element (additive identity) of the field.
func ZeroSecp256k1N() Secp256k1N {
	return NewSecp256k1N(0)
}

// OneSecp256k1N returns the one element (multiplicative identity) of the
// field.
func OneSecp256k1N() Secp256k1N {
	return NewSecp256k1N(1)
}

// RandomSecp256k1N returns a random field element.
func RandomSecp256k1N() Secp256k1N {
	val := make([]byte, 40)
	_, err := rand.Read(val)
	if err != nil {
		panic("could not generate a random byte")
	}
	ret := NewSecp256k1N(0)
	ret.limbs[0] = binary.LittleEndian.Uint64(val[0:]) >> 12
	ret.limbs[1] = binary.LittleEndian.Uint64(val[8:]) >> 12
	ret.limbs[2] = binary.LittleEndian.Uint64(val[16:]) >> 12
	ret.limbs[3] = binary.LittleEndian.Uint64(val[24:]) >> 12
	ret.limbs[4] = binary.LittleEndian.Uint64(val[32:]) >> 16

	return ret
}

// TODO: Should this reduce modulo N?
func (x *Secp256k1N) Int() *big.Int {
	ret := big.NewInt(0)

	for i := range x.limbs {
		ret.Lsh(ret, 52)
		ret.Add(ret, big.NewInt(0).SetUint64(x.limbs[4-i]))
	}

	return ret
}

func (x *Secp256k1N) Normalize() {
	t0, t1, t2, t3, t4 := x.limbs[0], x.limbs[1], x.limbs[2], x.limbs[3], x.limbs[4]

	y := t4 >> 48
	t4 &= 0xffffffffffff

	t0 += y * r0
	t1 += y*r1 + t0>>52
	t0 &= 0xfffffffffffff
	t2 += y*r2 + t1>>52
	t1 &= 0xfffffffffffff
	t3 += t2 >> 52
	t2 &= 0xfffffffffffff
	t4 += t3 >> 52
	t3 &= 0xfffffffffffff

	// TODO: Double check the logic here.
	if t4>>48 != 0 || ((t4 == 0xffffffffffff) && (t3 == 0xfffffffffffff) && (t2 > r2 || (t2 == r2 && (t1 > r1 || (t1 == r1 && r0 >= r0))))) {
		t0 += r0
		t1 += r1 + t0>>52
		t0 &= 0xfffffffffffff
		t2 += r2 + t1>>52
		t1 &= 0xfffffffffffff
		t3 += t2 >> 52
		t2 &= 0xfffffffffffff
		t4 += t3 >> 52
		t3 &= 0xfffffffffffff

		t4 &= 0xffffffffffff
	}

	x.limbs[0] = t0
	x.limbs[1] = t1
	x.limbs[2] = t2
	x.limbs[3] = t3
	x.limbs[4] = t4
}

// Add returns a new field element that is the sum of the two field elements.
func (x *Secp256k1N) Add(y, z *Secp256k1N) {
	x.limbs[0] = y.limbs[0] + z.limbs[0]
	x.limbs[1] = y.limbs[1] + z.limbs[1]
	x.limbs[2] = y.limbs[2] + z.limbs[2]
	x.limbs[3] = y.limbs[3] + z.limbs[3]
	x.limbs[4] = y.limbs[4] + z.limbs[4]
}

// Neg returns the additive inverse of a field element.
func (x *Secp256k1N) Neg(y *Secp256k1N, m uint) {
	x.limbs[0] = 0x25e8cd0364141*2*(uint64(m)+1) - y.limbs[0]
	x.limbs[1] = 0xe6af48a03bbfd*2*(uint64(m)+1) - y.limbs[1]
	x.limbs[2] = 0xffffffebaaedc*2*(uint64(m)+1) - y.limbs[2]
	x.limbs[3] = 0xfffffffffffff*2*(uint64(m)+1) - y.limbs[3]
	x.limbs[4] = 0x0ffffffffffff*2*(uint64(m)+1) - y.limbs[4]
}

// Mul returns a new field element that is the product of the two field
// elements.
func (x *Secp256k1N) Mul(y, z *Secp256k1N) {
	var tmp0, tmp1 uint64

	c0l, c0u := mulQQ(y.limbs[0], z.limbs[0])
	c0u <<= 12
	c0u += c0l >> 52
	c0l &= 0xfffffffffffff

	c1l, c1u := mulQQ(y.limbs[0], z.limbs[1])
	c1u <<= 12
	c1u += c1l >> 52
	c1l &= 0xfffffffffffff
	tmp0, tmp1 = mulQQ(y.limbs[1], z.limbs[0])
	c1l += tmp0 & 0xfffffffffffff
	c1u += tmp1<<12 + tmp0>>52
	c1u += c1l >> 52
	c1l &= 0xfffffffffffff

	c2l, c2u := mulQQ(y.limbs[0], z.limbs[2])
	c2u <<= 12
	c2u += c2l >> 52
	c2l &= 0xfffffffffffff
	tmp0, tmp1 = mulQQ(y.limbs[1], z.limbs[1])
	c2l += tmp0 & 0xfffffffffffff
	c2u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[2], z.limbs[0])
	c2l += tmp0 & 0xfffffffffffff
	c2u += tmp1<<12 + tmp0>>52
	c2u += c2l >> 52
	c2l &= 0xfffffffffffff

	c3l, c3u := mulQQ(y.limbs[0], z.limbs[3])
	c3u <<= 12
	c3u += c3l >> 52
	c3l &= 0xfffffffffffff
	tmp0, tmp1 = mulQQ(y.limbs[1], z.limbs[2])
	c3l += tmp0 & 0xfffffffffffff
	c3u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[2], z.limbs[1])
	c3l += tmp0 & 0xfffffffffffff
	c3u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[3], z.limbs[0])
	c3l += tmp0 & 0xfffffffffffff
	c3u += tmp1<<12 + tmp0>>52
	c3u += c3l >> 52
	c3l &= 0xfffffffffffff

	c4l, c4u := mulQQ(y.limbs[0], z.limbs[4])
	c4u <<= 12
	c4u += c4l >> 52
	c4l &= 0xfffffffffffff
	tmp0, tmp1 = mulQQ(y.limbs[1], z.limbs[3])
	c4l += tmp0 & 0xfffffffffffff
	c4u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[2], z.limbs[2])
	c4l += tmp0 & 0xfffffffffffff
	c4u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[3], z.limbs[1])
	c4l += tmp0 & 0xfffffffffffff
	c4u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[4], z.limbs[0])
	c4l += tmp0 & 0xfffffffffffff
	c4u += tmp1<<12 + tmp0>>52
	c4u += c4l >> 52
	c4l &= 0xfffffffffffff

	c5l, c5u := mulQQ(y.limbs[1], z.limbs[4])
	c5u <<= 12
	c5u += c5l >> 52
	c5l &= 0xfffffffffffff
	tmp0, tmp1 = mulQQ(y.limbs[2], z.limbs[3])
	c5l += tmp0 & 0xfffffffffffff
	c5u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[3], z.limbs[2])
	c5l += tmp0 & 0xfffffffffffff
	c5u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[4], z.limbs[1])
	c5l += tmp0 & 0xfffffffffffff
	c5u += tmp1<<12 + tmp0>>52
	c5u += c5l >> 52
	c5l &= 0xfffffffffffff

	c6l, c6u := mulQQ(y.limbs[2], z.limbs[4])
	c6u <<= 12
	c6u += c6l >> 52
	c6l &= 0xfffffffffffff
	tmp0, tmp1 = mulQQ(y.limbs[3], z.limbs[3])
	c6l += tmp0 & 0xfffffffffffff
	c6u += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(y.limbs[4], z.limbs[2])
	c6l += tmp0 & 0xfffffffffffff
	c6u += tmp1<<12 + tmp0>>52
	c6u += c6l >> 52
	c6l &= 0xfffffffffffff

	c7l, c7u := mulQQ(y.limbs[3], z.limbs[4])
	c7u <<= 12
	c7u += c7l >> 52
	c7l &= 0xfffffffffffff
	tmp0, tmp1 = mulQQ(y.limbs[4], z.limbs[3])
	c7l += tmp0 & 0xfffffffffffff
	c7u += tmp1<<12 + tmp0>>52
	c7u += c7l >> 52
	c7l &= 0xfffffffffffff

	c8l, c8u := mulQQ(y.limbs[4], z.limbs[4])
	c8u <<= 12
	c8u += c8l >> 52
	c8l &= 0xfffffffffffff

	// TODO: These values are constants and should be hard coded.
	tmp0, tmp1 = mulQQ(r0, r1)
	s01l := tmp0 & 0xfffffffffffff
	s01u := tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r0, r2)
	s02l := tmp0 & 0xfffffffffffff
	s02u := tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r1, r1)
	s11l := tmp0 & 0xfffffffffffff
	s11u := tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r1, r2)
	s12l := tmp0 & 0xfffffffffffff
	s12u := tmp1<<12 + tmp0>>52
	tmp0, _ = mulQQ(r2, r2)
	s22l := tmp0 & 0xfffffffffffff

	tmp0, tmp1 = mulQQ(r0, c8u)
	f0l := tmp0 & 0xfffffffffffff
	f0u := tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r1, c7u+c8l)
	f1l := tmp0 & 0xfffffffffffff
	f1u := tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r2, c6u+c7l)
	f2l := tmp0 & 0xfffffffffffff
	f2u := tmp1<<12 + tmp0>>52

	s01u += s01l >> 52
	s01l &= 0xfffffffffffff
	s02u += s02l >> 52
	s02l &= 0xfffffffffffff
	s11u += s11l >> 52
	s11l &= 0xfffffffffffff
	s12u += s12l >> 52
	s12l &= 0xfffffffffffff
	s22l &= 0xfffffffffffff
	f0u += f0l >> 52
	f0l &= 0xfffffffffffff
	f1u += f1l >> 52
	f1l &= 0xfffffffffffff
	f2u += f2l >> 52
	f2l &= 0xfffffffffffff

	x.limbs[0] = c0l
	x.limbs[1] = c0u + c1l
	x.limbs[2] = c1u + c2l
	x.limbs[3] = c2u + c3l
	x.limbs[4] = c3u + c4l + (f0l+f1l+f2l)<<4

	tmp0, tmp1 = mulQQ(r0<<4, c4u+c5l)
	x.limbs[0] += tmp0 & 0xfffffffffffff
	x.limbs[1] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s01l<<8, c8u)
	x.limbs[0] += tmp0 & 0xfffffffffffff
	x.limbs[1] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s02l<<8, c7u+c8l)
	x.limbs[0] += tmp0 & 0xfffffffffffff
	x.limbs[1] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r0<<8, f0u+f1u+f2u)
	x.limbs[0] += tmp0 & 0xfffffffffffff
	x.limbs[1] += tmp1<<12 + tmp0>>52

	tmp0, tmp1 = mulQQ(r0<<4, c5u+c6l)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r1<<4, c4u+c5l)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s02l<<8, c8u)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s11l<<8, c8u)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s12l<<8, c7u+c8l)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s01u<<8, c8u)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s02u<<8, c7u+c8l)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r1<<8, f0u+f1u+f2u)
	x.limbs[1] += tmp0 & 0xfffffffffffff
	x.limbs[2] += tmp1<<12 + tmp0>>52

	tmp0, tmp1 = mulQQ(r0<<4, c6u+c7l)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r1<<4, c5u+c6l)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r2<<4, c4u+c5l)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s12l<<9, c8u)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s22l<<8, c7u+c8l)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s02u<<8, c8u)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s11u<<8, c8u)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s12u<<8, c7u+c8l)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r2<<8, f0u+f1u+f2u)
	x.limbs[2] += tmp0 & 0xfffffffffffff
	x.limbs[3] += tmp1<<12 + tmp0>>52

	tmp0, tmp1 = mulQQ(r0<<4, c7u+c8l)
	x.limbs[3] += tmp0 & 0xfffffffffffff
	x.limbs[4] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r1<<4, c6u+c7l)
	x.limbs[3] += tmp0 & 0xfffffffffffff
	x.limbs[4] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(r2<<4, c5u+c6l)
	x.limbs[3] += tmp0 & 0xfffffffffffff
	x.limbs[4] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s22l<<8, c8u)
	x.limbs[3] += tmp0 & 0xfffffffffffff
	x.limbs[4] += tmp1<<12 + tmp0>>52
	tmp0, tmp1 = mulQQ(s12u<<9, c8u)
	x.limbs[3] += tmp0 & 0xfffffffffffff
	x.limbs[4] += tmp1<<12 + tmp0>>52
}

// Inv returns a new field element that is the multiplicative inverse of the
// given field element. If the field element is the zero element, then the
// function will panic.
func (x *Secp256k1N) Inv() {
	panic("unimplemented")
}

// Eq returns true if the two field elements are equal, and false otherwise.
func (x *Secp256k1N) Eq(y *Secp256k1N) bool {
	panic("unimplemented")
}

// IsZero returns true if the field element is equal to the zero element
// (additive identity), and false otherwise.
func (x *Secp256k1N) IsZero() bool {
	panic("unimplemented")
}

// IsOne returns true if the field element is equal to the one element
// (multiplicative identity), and false otherwise.
func (x *Secp256k1N) IsOne() bool {
	panic("unimplemented")
}
