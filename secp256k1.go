package secp256k1

/*
#include <stdlib.h>

#include "c-field/field_impl.h"
#include "secp256k1n.h"
*/
import "C"
import (
	"crypto/rand"
	"encoding/binary"
	"math/big"
	mrand "math/rand"
	"reflect"
	"unsafe"
)

const r0 uint64 = 0xda1732fc9bebf
const r1 uint64 = 0x1950b75fc4402
const r2 uint64 = 0x1455123

type Secp256k1P struct {
	limbs [5]uint64
}

func NewSecp256k1P(x uint64) Secp256k1P {
	return Secp256k1P{limbs: [5]uint64{x, 0, 0, 0, 0}}
}

// Secp256k1PFromInt creates a new field element from a big.Int. The big.Int is
// truncated to 32 bytes.
func Secp256k1PFromInt(n *big.Int) Secp256k1P {
	var ret Secp256k1P
	nBytes := n.Bytes()
	var bs [32]byte
	copy(bs[32-len(nBytes):], nBytes)
	ret.SetB32(bs[:])
	return ret
}

// RandomSecp256k1P returns a random field element.
func RandomSecp256k1P() Secp256k1P {
	val := [32]byte{}
	// Always returns nil error
	rand.Read(val[:])
	ret := NewSecp256k1P(0)
	ret.SetB32(val[:])

	return ret
}

// ZeroSecp256k1P returns the zero element (additive identity) of the field.
func ZeroSecp256k1P() Secp256k1P {
	return NewSecp256k1P(0)
}

// OneSecp256k1P returns the one element (multiplicative identity) of the
// field.
func OneSecp256k1P() Secp256k1P {
	return NewSecp256k1P(1)
}

// Cast converts sets x (P field) to have the value of y (N field). Note that
// if x is desired to have the same numerical value of y then y needs to be
// normalised first.
func (x *Secp256k1P) Cast(y *Secp256k1N) {
	x.limbs = y.limbs
}

// Generate implements the quick.Generator interface.
func (x Secp256k1P) Generate(r *mrand.Rand, size int) reflect.Value {
	// TODO: We don't use the provided rng here. Does this matter?
	ret := RandomSecp256k1P()
	return reflect.ValueOf(ret)
}

func (r *Secp256k1P) Set(a *Secp256k1P) {
	r.limbs = a.limbs
}

func (r *Secp256k1P) Int() *big.Int {
	return big.NewInt(0).SetBytes(r.GetB32())
}

func (r *Secp256k1P) Normalize() {
	C.secp256k1_fe_normalize((*C.secp256k1_fe)(unsafe.Pointer(r)))
}

func (r *Secp256k1P) NormalizeWeak() {
	C.secp256k1_fe_normalize_weak((*C.secp256k1_fe)(unsafe.Pointer(r)))
}

func (r *Secp256k1P) NormalizeVar() {
	C.secp256k1_fe_normalize_var((*C.secp256k1_fe)(unsafe.Pointer(r)))
}

func (r *Secp256k1P) NormalizesToZero() bool {
	return C.secp256k1_fe_normalizes_to_zero((*C.secp256k1_fe)(unsafe.Pointer(r))) != 0
}

func (r *Secp256k1P) NormalizesToZeroVar() bool {
	return C.secp256k1_fe_normalizes_to_zero_var((*C.secp256k1_fe)(unsafe.Pointer(r))) != 0
}

func (r *Secp256k1P) SetUint64(a uint64) {
	C.secp256k1_fe_set_int((*C.secp256k1_fe)(unsafe.Pointer(r)), C.uint64_t(a))
}

func (r *Secp256k1P) Clear() {
	C.secp256k1_fe_clear((*C.secp256k1_fe)(unsafe.Pointer(r)))
}

func (r *Secp256k1P) IsZero() bool {
	return C.secp256k1_fe_is_zero((*C.secp256k1_fe)(unsafe.Pointer(r))) != 0
}

func (r *Secp256k1P) IsOdd() bool {
	return C.secp256k1_fe_is_odd((*C.secp256k1_fe)(unsafe.Pointer(r))) != 0
}

func (r *Secp256k1P) Eq(a *Secp256k1P) bool {
	return C.secp256k1_fe_equal((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a))) != 0
}

func (r *Secp256k1P) EqVar(a *Secp256k1P) bool {
	return C.secp256k1_fe_equal_var((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a))) != 0
}

func (r *Secp256k1P) CmpVar(a *Secp256k1P) int {
	return int(C.secp256k1_fe_cmp_var((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a))))
}

func (r *Secp256k1P) SetB32(b []byte) {
	cBytes := (*C.uchar)(C.CBytes(b))
	C.secp256k1_fe_set_b32((*C.secp256k1_fe)(unsafe.Pointer(r)), cBytes)
	C.free(unsafe.Pointer(cBytes))
}

func (r *Secp256k1P) GetB32() []byte {
	cBytes := C.malloc(C.sizeof_char * 32)
	C.secp256k1_fe_get_b32((*C.uchar)(cBytes), (*C.secp256k1_fe)(unsafe.Pointer(r)))
	b := C.GoBytes(cBytes, 32)
	C.free(unsafe.Pointer(cBytes))
	return b
}

func (r *Secp256k1P) Uint64() uint64 {
	return r.limbs[0] + (r.limbs[1]&0xfff)<<52
}

func (r *Secp256k1P) Neg(a *Secp256k1P, m int) {
	C.secp256k1_fe_negate((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)), C.int(m))
}

func (r *Secp256k1P) MulInt(a int) {
	C.secp256k1_fe_mul_int((*C.secp256k1_fe)(unsafe.Pointer(r)), C.int(a))
}

func (r *Secp256k1P) Add(a *Secp256k1P) {
	C.secp256k1_fe_add((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)))
}

func (r *Secp256k1P) Mul(a, b *Secp256k1P) {
	C.secp256k1_fe_mul((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)), (*C.secp256k1_fe)(unsafe.Pointer(b)))
}

func (r *Secp256k1P) Sqr(a *Secp256k1P) {
	C.secp256k1_fe_sqr((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)))
}

func (r *Secp256k1P) Sqrt(a *Secp256k1P) bool {
	return C.secp256k1_fe_sqrt((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a))) != 0
}

func (r *Secp256k1P) IsQuadVar() bool {
	return C.secp256k1_fe_is_quad_var((*C.secp256k1_fe)(unsafe.Pointer(r))) != 0
}

func (r *Secp256k1P) Inv(a *Secp256k1P) {
	C.secp256k1_fe_inv((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)))
}

func (r *Secp256k1P) InvVar(a *Secp256k1P) {
	C.secp256k1_fe_inv_var((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)))
}

func (r *Secp256k1P) Cmov(a *Secp256k1P, flag bool) {
	if flag {
		C.secp256k1_fe_cmov((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)), 1)
	} else {
		C.secp256k1_fe_cmov((*C.secp256k1_fe)(unsafe.Pointer(r)), (*C.secp256k1_fe)(unsafe.Pointer(a)), 0)
	}
}

// Secp256k1N represents an element of the field defined by the prime that is
// the order of the elliptic curve group Secp256k1.
type Secp256k1N struct {
	limbs [5]uint64
}

// NewSecp256k1N returns a new field element with value equal to x.
func NewSecp256k1N(x uint64) Secp256k1N {
	return Secp256k1N{limbs: [5]uint64{x, 0, 0, 0, 0}}
}

// Secp256k1NFromInt creates a new field element from a big.Int. The big.Int is
// truncated to 32 bytes.
func Secp256k1NFromInt(n *big.Int) Secp256k1N {
	var ret Secp256k1N
	nBytes := n.Bytes()
	var bs [32]byte
	copy(bs[32-len(nBytes):], nBytes)
	ret.SetB32(bs[:])
	return ret
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

// Set copies the value of y into x.
func (x *Secp256k1N) Set(y *Secp256k1N) {
	x.limbs = y.limbs
}

// Cast converts sets x (N field) to have the value of y (P field). Note that
// if x is desired to have the same numerical value of y (possibly modulo N)
// then y needs to be normalised first.
func (x *Secp256k1N) Cast(y *Secp256k1P) {
	x.limbs = y.limbs
}

// RandomSecp256k1N returns a random field element.
func RandomSecp256k1N() Secp256k1N {
	val := [40]byte{}
	// Always returns nil error
	rand.Read(val[:])
	ret := NewSecp256k1N(0)
	ret.limbs[0] = binary.LittleEndian.Uint64(val[0:]) >> 12
	ret.limbs[1] = binary.LittleEndian.Uint64(val[8:]) >> 12
	ret.limbs[2] = binary.LittleEndian.Uint64(val[16:]) >> 12
	ret.limbs[3] = binary.LittleEndian.Uint64(val[24:]) >> 12
	ret.limbs[4] = binary.LittleEndian.Uint64(val[32:]) >> 16

	return ret
}

// Generate implements the quick.Generator interface.
func (x Secp256k1N) Generate(r *mrand.Rand, size int) reflect.Value {
	// TODO: We don't use the provided rng here. Does this matter?
	ret := RandomSecp256k1N()
	return reflect.ValueOf(ret)
}

// Int returns a big.Int that has the same value that x represents (reduced
// modulo N).
func (x *Secp256k1N) Int() *big.Int {
	ret := big.NewInt(0)
	n, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	for i := range x.limbs {
		ret.Lsh(ret, 52)
		ret.Add(ret, big.NewInt(0).SetUint64(x.limbs[4-i]))
	}
	ret.Mod(ret, n)

	return ret
}

// SetB32 sets the value of x to be the given big endian bytes. It is assumed
// that the byte slice has at least 32 elements, otherwise it will panic.
func (x *Secp256k1N) SetB32(a []byte) {
	x.limbs[0] = uint64(a[31]) | (uint64(a[30]) << 8) | (uint64(a[29]) << 16) | (uint64(a[28]) << 24) | (uint64(a[27]) << 32) | (uint64(a[26]) << 40) | (uint64(a[25]&0xF) << 48)
	x.limbs[1] = uint64((a[25]>>4)&0xF) | (uint64(a[24]) << 4) | (uint64(a[23]) << 12) | (uint64(a[22]) << 20) | (uint64(a[21]) << 28) | (uint64(a[20]) << 36) | (uint64(a[19]) << 44)
	x.limbs[2] = uint64(a[18]) | (uint64(a[17]) << 8) | (uint64(a[16]) << 16) | (uint64(a[15]) << 24) | (uint64(a[14]) << 32) | (uint64(a[13]) << 40) | (uint64(a[12]&0xF) << 48)
	x.limbs[3] = uint64((a[12]>>4)&0xF) | (uint64(a[11]) << 4) | (uint64(a[10]) << 12) | (uint64(a[9]) << 20) | (uint64(a[8]) << 28) | (uint64(a[7]) << 36) | (uint64(a[6]) << 44)
	x.limbs[4] = uint64(a[5]) | (uint64(a[4]) << 8) | (uint64(a[3]) << 16) | (uint64(a[2]) << 24) | (uint64(a[1]) << 32) | (uint64(a[0]) << 40)
}

// GetB32 writes the big endian bytes of x to the given slice a. It is assumed
// that the byte slice has length at least 32, otherwise it will panic.
func (x *Secp256k1N) GetB32(a []byte) {
	a[0] = byte((x.limbs[4] >> 40) & 0xFF)
	a[1] = byte((x.limbs[4] >> 32) & 0xFF)
	a[2] = byte((x.limbs[4] >> 24) & 0xFF)
	a[3] = byte((x.limbs[4] >> 16) & 0xFF)
	a[4] = byte((x.limbs[4] >> 8) & 0xFF)
	a[5] = byte(x.limbs[4] & 0xFF)
	a[6] = byte((x.limbs[3] >> 44) & 0xFF)
	a[7] = byte((x.limbs[3] >> 36) & 0xFF)
	a[8] = byte((x.limbs[3] >> 28) & 0xFF)
	a[9] = byte((x.limbs[3] >> 20) & 0xFF)
	a[10] = byte((x.limbs[3] >> 12) & 0xFF)
	a[11] = byte((x.limbs[3] >> 4) & 0xFF)
	a[12] = byte(((x.limbs[2] >> 48) & 0xF) | ((x.limbs[3] & 0xF) << 4))
	a[13] = byte((x.limbs[2] >> 40) & 0xFF)
	a[14] = byte((x.limbs[2] >> 32) & 0xFF)
	a[15] = byte((x.limbs[2] >> 24) & 0xFF)
	a[16] = byte((x.limbs[2] >> 16) & 0xFF)
	a[17] = byte((x.limbs[2] >> 8) & 0xFF)
	a[18] = byte(x.limbs[2] & 0xFF)
	a[19] = byte((x.limbs[1] >> 44) & 0xFF)
	a[20] = byte((x.limbs[1] >> 36) & 0xFF)
	a[21] = byte((x.limbs[1] >> 28) & 0xFF)
	a[22] = byte((x.limbs[1] >> 20) & 0xFF)
	a[23] = byte((x.limbs[1] >> 12) & 0xFF)
	a[24] = byte((x.limbs[1] >> 4) & 0xFF)
	a[25] = byte(((x.limbs[0] >> 48) & 0xF) | ((x.limbs[1] & 0xF) << 4))
	a[26] = byte((x.limbs[0] >> 40) & 0xFF)
	a[27] = byte((x.limbs[0] >> 32) & 0xFF)
	a[28] = byte((x.limbs[0] >> 24) & 0xFF)
	a[29] = byte((x.limbs[0] >> 16) & 0xFF)
	a[30] = byte((x.limbs[0] >> 8) & 0xFF)
	a[31] = byte(x.limbs[0] & 0xFF)
}

// Uint64 returns x truncated to 64 bits. Note that this will only return
// actual value of x if it is normalised.
func (x *Secp256k1N) Uint64() uint64 {
	return x.limbs[0] + (x.limbs[1]&0xfff)<<52
}

// Normalize reduces the limbed representation of x so that it is less than the
// prime and all of the limbs are in valid base 52 ranges.
//
// The magnitude of the caller should be no more than 2048.
func (x *Secp256k1N) Normalize() {
	t0, t1, t2, t3, t4 := x.limbs[0], x.limbs[1], x.limbs[2], x.limbs[3], x.limbs[4]

	y := t4 >> 48
	t4 &= 0x0ffffffffffff

	t0 += y * r0
	t1 += y*r1 + t0>>52
	t0 &= 0xfffffffffffff
	t2 += y*r2 + t1>>52
	t1 &= 0xfffffffffffff
	t3 += t2 >> 52
	t2 &= 0xfffffffffffff
	t4 += t3 >> 52
	t3 &= 0xfffffffffffff

	if t4>>48 != 0 || ((t4 == 0x0ffffffffffff) && (t3 == 0xfffffffffffff) && (t2 > 0xffffffebaaedc || (t2 == 0xffffffebaaedc && (t1 > 0xe6af48a03bbfd || (t1 == 0xe6af48a03bbfd && t0 >= 0x25e8cd0364141))))) {
		t0 += r0
		t1 += r1 + t0>>52
		t0 &= 0xfffffffffffff
		t2 += r2 + t1>>52
		t1 &= 0xfffffffffffff
		t3 += t2 >> 52
		t2 &= 0xfffffffffffff
		t4 += t3 >> 52
		t3 &= 0xfffffffffffff

		t4 &= 0x0ffffffffffff
	}

	x.limbs[0] = t0
	x.limbs[1] = t1
	x.limbs[2] = t2
	x.limbs[3] = t3
	x.limbs[4] = t4
}

// Add computes the field addition of y and z and stores it in x.
//
// The output will have magnitude equal to the sum of the magnitudes of the
// input values.
func (x *Secp256k1N) Add(y, z *Secp256k1N) {
	x.limbs[0] = y.limbs[0] + z.limbs[0]
	x.limbs[1] = y.limbs[1] + z.limbs[1]
	x.limbs[2] = y.limbs[2] + z.limbs[2]
	x.limbs[3] = y.limbs[3] + z.limbs[3]
	x.limbs[4] = y.limbs[4] + z.limbs[4]
}

// Neg computes the additive inverse of y and stores it in x. The second
// argument is the magnitude of y.
//
// The output will have magnitude equal to 7m.
func (x *Secp256k1N) Neg(y *Secp256k1N, m uint) {
	// 7 is the smallest number k such that
	// 		0x25e8cd0364141 * k > 2^52
	// and so this choice allows us to avoid underflow for a normalised input
	x.limbs[0] = 0x25e8cd0364141*7*uint64(m) - y.limbs[0]
	x.limbs[1] = 0xe6af48a03bbfd*7*uint64(m) - y.limbs[1]
	x.limbs[2] = 0xffffffebaaedc*7*uint64(m) - y.limbs[2]
	x.limbs[3] = 0xfffffffffffff*7*uint64(m) - y.limbs[3]
	x.limbs[4] = 0x0ffffffffffff*7*uint64(m) - y.limbs[4]
}

// Mul Performs the field multiplication of y and z and stores it in x.
// NOTE: Until the relevant fixme in the implementation of C.secp256k1_mul is
// addressed, the case where x == z (as pointers) will give incorrect results.
// In the meantime, if z == x then call the function with the arguments
// swapped; x.Mul(x, y) instead of x.Mul(y, x). If it also the case that x ==
// y, then x.Sqr(x) should be used instead.
func (x *Secp256k1N) Mul(y, z *Secp256k1N) {
	C.secp256k1n_mul((*C.secp256k1n)(unsafe.Pointer(x)), (*C.secp256k1n)(unsafe.Pointer(y)), (*C.secp256k1n)(unsafe.Pointer(z)))
}

// Sqr computes the square of y and stores it in x. That is, x = y*y.
func (x *Secp256k1N) Sqr(y *Secp256k1N) {
	C.secp256k1n_sqr((*C.secp256k1n)(unsafe.Pointer(x)), (*C.secp256k1n)(unsafe.Pointer(y)))
}

// Inv computes the multiplicative inverse of y and stores it in x. If y is
// equal to the zero element, then the result will also be the zero element.
func (x *Secp256k1N) Inv(y *Secp256k1N) {
	C.secp256k1n_inv((*C.secp256k1n)(unsafe.Pointer(x)), (*C.secp256k1n)(unsafe.Pointer(y)))
}

// MulInt multiplies x by a small int value
func (x *Secp256k1N) MulInt(a int) {
	x.limbs[0] *= uint64(a)
	x.limbs[1] *= uint64(a)
	x.limbs[2] *= uint64(a)
	x.limbs[3] *= uint64(a)
	x.limbs[4] *= uint64(a)
}

// Eq returns true if the two field elements are equal, and false otherwise.
func (x *Secp256k1N) Eq(y *Secp256k1N) bool {
	// TODO: More efficient implementation/
	var xNorm, yNorm = *x, *y
	xNorm.Normalize()
	yNorm.Normalize()
	return xNorm.limbs[0] == yNorm.limbs[0] &&
		xNorm.limbs[1] == yNorm.limbs[1] &&
		xNorm.limbs[2] == yNorm.limbs[2] &&
		xNorm.limbs[3] == yNorm.limbs[3] &&
		xNorm.limbs[4] == yNorm.limbs[4]
}

// Clear sets the value of x to zero by setting all of the limbs to zero.
func (x *Secp256k1N) Clear() {
	x.limbs[0] = 0
	x.limbs[1] = 0
	x.limbs[2] = 0
	x.limbs[3] = 0
	x.limbs[4] = 0
}

// IsZero returns true if the field element is equal to the zero element
// (additive identity), and false otherwise.
func (x *Secp256k1N) IsZero() bool {
	var z Secp256k1N
	z.Set(x)
	z.Normalize()
	return (z.limbs[0] | z.limbs[1] | z.limbs[2] | z.limbs[3] | z.limbs[4]) == 0
}

// IsOne returns true if the field element is equal to the one element
// (multiplicative identity), and false otherwise.
func (x *Secp256k1N) IsOne() bool {
	var z Secp256k1N
	z.Set(x)
	z.Normalize()
	return z.limbs[0] == 1 && (z.limbs[1]|z.limbs[2]|z.limbs[3]|z.limbs[4]) == 0
}

// IsOdd returns true if x is odd. Note that to ensure that this is the parity
// of the underlying field element, x should be normalised.
func (x *Secp256k1N) IsOdd() bool {
	return x.limbs[0]&1 == 1
}
