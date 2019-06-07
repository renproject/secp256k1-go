package secp256k1

/*
#include <stdlib.h>

#include "c-field/field_impl.h"
*/
import "C"
import (
	"math/big"
	"unsafe"
)

type Secp256k1 struct {
	inner C.struct_secp256k1_fe
}

func New(a uint64) Secp256k1 {
	inner := C.struct_secp256k1_fe{}
	C.secp256k1_fe_set_int(&inner, C.uint64_t(a))
	return Secp256k1{inner}
}

func (r *Secp256k1) Set(a *Secp256k1) {
	r.inner = a.inner
}

func (r *Secp256k1) BigInt() *big.Int {
	return big.NewInt(0).SetBytes(r.GetB32())
}

func (r *Secp256k1) Normalize() {
	C.secp256k1_fe_normalize(&r.inner)
}

func (r *Secp256k1) NormalizeWeak() {
	C.secp256k1_fe_normalize_weak(&r.inner)
}

func (r *Secp256k1) NormalizeVar() {
	C.secp256k1_fe_normalize_var(&r.inner)
}

func (r *Secp256k1) NormalizesToZero() bool {
	return C.secp256k1_fe_normalizes_to_zero(&r.inner) != 0
}

func (r *Secp256k1) NormalizesToZeroVar() bool {
	return C.secp256k1_fe_normalizes_to_zero_var(&r.inner) != 0
}

func (r *Secp256k1) SetUint64(a uint64) {
	C.secp256k1_fe_set_int(&r.inner, C.uint64_t(a))
}

func (r *Secp256k1) Clear() {
	C.secp256k1_fe_clear(&r.inner)
}

func (r *Secp256k1) IsZero() bool {
	return C.secp256k1_fe_is_zero(&r.inner) != 0
}

func (r *Secp256k1) IsOdd() bool {
	return C.secp256k1_fe_is_odd(&r.inner) != 0
}

func (r *Secp256k1) Equal(a *Secp256k1) bool {
	return C.secp256k1_fe_equal(&r.inner, &a.inner) != 0
}

func (r *Secp256k1) EqualVar(a *Secp256k1) bool {
	return C.secp256k1_fe_equal_var(&r.inner, &a.inner) != 0
}

func (r *Secp256k1) CmpVar(a *Secp256k1) int {
	return int(C.secp256k1_fe_cmp_var(&r.inner, &a.inner))
}

func (r *Secp256k1) SetB32(b []byte) {
	cBytes := (*C.uchar)(C.CBytes(b))
	C.secp256k1_fe_set_b32(&r.inner, cBytes)
	C.free(unsafe.Pointer(cBytes))
}

func (r *Secp256k1) GetB32() []byte {
	cBytes := C.malloc(C.sizeof_char * 32)
	C.secp256k1_fe_get_b32((*C.uchar)(cBytes), &r.inner)
	b := C.GoBytes(cBytes, 32)
	C.free(unsafe.Pointer(cBytes))
	return b
}

func (r *Secp256k1) Negate(a *Secp256k1, m int) {
	C.secp256k1_fe_negate(&r.inner, &a.inner, C.int(m))
}

func (r *Secp256k1) MulInt(a int) {
	C.secp256k1_fe_mul_int(&r.inner, C.int(a))
}

func (r *Secp256k1) Add(a *Secp256k1) {
	C.secp256k1_fe_add(&r.inner, &a.inner)
}

func (r *Secp256k1) Mul(a, b *Secp256k1) {
	C.secp256k1_fe_mul(&r.inner, &a.inner, &b.inner)
}

func (r *Secp256k1) Sqr(a *Secp256k1) {
	C.secp256k1_fe_sqr(&r.inner, &a.inner)
}

func (r *Secp256k1) Sqrt(a *Secp256k1) bool {
	return C.secp256k1_fe_sqrt(&r.inner, &a.inner) != 0
}

func (r *Secp256k1) IsQuadVar() bool {
	return C.secp256k1_fe_is_quad_var(&r.inner) != 0
}

func (r *Secp256k1) Inv(a *Secp256k1) {
	C.secp256k1_fe_inv(&r.inner, &a.inner)
}

func (r *Secp256k1) InvVar(a *Secp256k1) {
	C.secp256k1_fe_inv_var(&r.inner, &a.inner)
}

func (r *Secp256k1) Cmov(a *Secp256k1, flag bool) {
	if flag {
		C.secp256k1_fe_cmov(&r.inner, &a.inner, 1)
	} else {
		C.secp256k1_fe_cmov(&r.inner, &a.inner, 0)
	}
}
