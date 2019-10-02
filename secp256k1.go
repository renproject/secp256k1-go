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

func (x *Secp256k1N) Set(y *Secp256k1N) {
	x.limbs = y.limbs
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

func (x Secp256k1N) Generate(r *mrand.Rand, size int) reflect.Value {
	// TODO: We don't use the provided rng here. Does this matter?
	ret := RandomSecp256k1N()
	return reflect.ValueOf(ret)
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
	C.secp256k1n_mul((*C.secp256k1n)(unsafe.Pointer(x)), (*C.secp256k1n)(unsafe.Pointer(y)), (*C.secp256k1n)(unsafe.Pointer(z)))
}

func (x *Secp256k1N) Sqr(y *Secp256k1N) {
	C.secp256k1n_sqr((*C.secp256k1n)(unsafe.Pointer(x)), (*C.secp256k1n)(unsafe.Pointer(y)))
}

// Inv returns a new field element that is the multiplicative inverse of the
// given field element. If the field element is the zero element, then the
// function will panic.
func (x *Secp256k1N) Inv(y *Secp256k1N) {
	C.secp256k1n_inv((*C.secp256k1n)(unsafe.Pointer(x)), (*C.secp256k1n)(unsafe.Pointer(y)))
}

// Eq returns true if the two field elements are equal, and false otherwise.
func (x *Secp256k1N) Eq(y *Secp256k1N) bool {
	// TODO: More efficient implementation/
	var z Secp256k1N
	z.Neg(x, 1)
	z.Add(&z, y)
	z.Normalize()
	return (z.limbs[0] | z.limbs[1] | z.limbs[2] | z.limbs[3] | z.limbs[4]) == 0
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
	return (z.limbs[0]|z.limbs[1]|z.limbs[2]|z.limbs[3]) == 0 && z.limbs[4] == 1
}
