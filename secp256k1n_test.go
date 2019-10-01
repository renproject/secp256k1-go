package secp256k1

import "testing"

func BenchmarkAddN(b *testing.B) {
	lhs := make([]Secp256k1N, 100)
	rhs := make([]Secp256k1N, 100)
	res := NewSecp256k1N(0)

	for i := range lhs {
		lhs[i] = RandomSecp256k1N()
		rhs[i] = RandomSecp256k1N()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l := lhs[(i/100)%100]
		r := rhs[i%100]
		res.Add(&l, &r)
	}
}

func BenchmarkMulN(b *testing.B) {
	lhs := make([]Secp256k1N, 100)
	rhs := make([]Secp256k1N, 100)
	res := NewSecp256k1N(0)

	for i := range lhs {
		lhs[i] = RandomSecp256k1N()
		rhs[i] = RandomSecp256k1N()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		l := lhs[(i/100)%100]
		r := rhs[i%100]
		res.Mul(&l, &r)
	}
}

func BenchmarkInvN(b *testing.B) {
	x := make([]Secp256k1N, 100)
	res := NewSecp256k1N(0)

	for i := range x {
		x[i] = RandomSecp256k1N()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := x[i%100]
		res.Inv(&r)
	}
}
