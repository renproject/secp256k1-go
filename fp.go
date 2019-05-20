package secp256k1

/*
#include <stdlib.h>

#include "c-field/field_impl.h"
#include "c-field/field_5x52.h"
#include "c-field/field_5x52_impl.h"
*/
import "C"
import (
	"math/big"
	"unsafe"
)

type Fp struct {
	inner C.struct_secp256k1_fe
}

func NewFp(a uint64) Fp {
	inner := C.struct_secp256k1_fe{}
	C.secp256k1_fe_set_int(&inner, C.ulong(a))
	return Fp{inner}
}

func (r *Fp) Set(a *Fp) {
	r.inner = a.inner
}

func (r *Fp) SetBytes(b []byte) {
	b32 := make([]byte, 32)
	copy(b32[32-len(b):], b)
	cBytes := (*C.uchar)(C.CBytes(b32))
	C.secp256k1_fe_set_b32(&r.inner, cBytes)
	C.free(unsafe.Pointer(cBytes))
}

func (r *Fp) Bytes() []byte {
	cBytes := C.malloc(C.sizeof_char * 32)
	C.secp256k1_fe_get_b32((*C.uchar)(cBytes), &r.inner)
	b := C.GoBytes(cBytes, 32)
	C.free(unsafe.Pointer(cBytes))
	return b
}

func (r *Fp) BigInt() *big.Int {
	return big.NewInt(0).SetBytes(r.Bytes())
}

func (r *Fp) Add(a, b *Fp) {
	r.Set(a)
	C.secp256k1_fe_add(&r.inner, &b.inner)
	C.secp256k1_fe_normalize_var(&r.inner)
}

func (r *Fp) Neg(a *Fp) {
	C.secp256k1_fe_negate(&r.inner, &a.inner, 0)
	C.secp256k1_fe_normalize_var(&r.inner)
}

func (r *Fp) Mul(a, b *Fp) {
	C.secp256k1_fe_mul(&r.inner, &a.inner, &b.inner)
	C.secp256k1_fe_normalize_var(&r.inner)
}

func (r *Fp) Inv(a *Fp) {
	C.secp256k1_fe_inv_var(&r.inner, &a.inner)
	C.secp256k1_fe_normalize_var(&r.inner)
}
