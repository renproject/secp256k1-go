package secp256k1_test

import (
	"math/big"
	"testing/quick"

	. "github.com/onsi/ginkgo"
	"github.com/renproject/secp256k1-go"
)

var _ = Describe("Wrapped field elements", func() {
	P, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	N, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	Context("When doing arithmetic in Fp", func() {
		It("Should add correctly", func() {
			quick.Check(func(x, y secp256k1.Secp256k1P) bool {
				sum := big.NewInt(0).Add(x.Int(), y.Int())
				sum.Mod(sum, P)
				z := secp256k1.NewSecp256k1P(0)

				z.Set(&x)
				z.Add(&y)
				z.NormalizeVar()

				return z.Int().Cmp(sum) == 0
			}, nil)
		})

		It("Should negate correctly", func() {
			quick.Check(func(x secp256k1.Secp256k1P) bool {
				neg := big.NewInt(0).Sub(P, big.NewInt(0).Mod(x.Int(), P))
				z := secp256k1.NewSecp256k1P(0)

				z.Negate(&x, 0)
				z.NormalizeVar()

				return z.Int().Cmp(neg) == 0
			}, nil)
		})

		It("Should multiply correctly", func() {
			quick.Check(func(x, y secp256k1.Secp256k1P) bool {
				prod := big.NewInt(0).Mul(x.Int(), y.Int())
				prod.Mod(prod, P)
				z := secp256k1.NewSecp256k1P(0)

				z.Mul(&x, &y)
				z.NormalizeVar()

				return z.Int().Cmp(prod) == 0
			}, nil)
		})

		It("Should invert correctly", func() {
			quick.Check(func(x secp256k1.Secp256k1P) bool {
				inv := big.NewInt(0).ModInverse(x.Int(), P)
				z := secp256k1.NewSecp256k1P(0)

				z.Inv(&x)
				z.NormalizeVar()

				return z.Int().Cmp(inv) == 0
			}, nil)
		})
	})

	Context("When doing arithmetic in Fn", func() {
		It("Should add correctly", func() {
			quick.Check(func(x, y secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				sum := big.NewInt(0).Add(x.Int(), y.Int())
				sum.Mod(sum, N)
				z.Add(&x, &y)
				z.Normalize()
				return z.Int().Cmp(sum) == 0
			}, nil)
		})

		It("Should negate correctly", func() {
			quick.Check(func(x secp256k1.Secp256k1N) bool {
				var y secp256k1.Secp256k1N
				neg := x.Int()
				neg.Mod(neg, N)
				neg.Sub(N, neg)

				y.Neg(&x, 0)
				y.Normalize()

				return y.Int().Cmp(neg) == 0
			}, nil)
		})

		It("Should multiply correctly", func() {
			quick.Check(func(x, y secp256k1.Secp256k1N) bool {
				z := secp256k1.NewSecp256k1N(0)
				prod := big.NewInt(0).Mul(x.Int(), y.Int())
				prod.Mod(prod, N)

				z.Mul(&x, &y)
				z.Normalize()

				return z.Int().Cmp(prod) == 0
			}, nil)
		})

		It("Should invert correctly", func() {
			quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				inv := big.NewInt(0).ModInverse(x.Int(), N)

				z.Inv(&x)
				z.Normalize()

				return z.Int().Cmp(inv) == 0
			}, nil)
		})

		It("Should check equaty correctly", func() {
			quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z, inv secp256k1.Secp256k1N
				z.Set(&x)

				// Try to get a different limbed representation
				inv.Inv(&z)
				z.Sqr(&z)
				z.Mul(&z, &inv)

				return z.Eq(&x)
			}, nil)
		})

		It("Should correctly construct the zero element", func() {
			quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				zero := secp256k1.ZeroSecp256k1N()

				// The defining property of the zero element is that it is the
				// additive identity
				z.Add(&x, &zero)

				return z.Eq(&x)
			}, nil)
		})

		It("Should correctly construct the one element", func() {
			quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				one := secp256k1.OneSecp256k1N()

				// The defining property of the zero element is that it is the
				// multiplicative identity
				z.Mul(&x, &one)

				return z.Eq(&x)
			}, nil)
		})

		It("Should correctly identify the zero element", func() {
			quick.Check(func(x secp256k1.Secp256k1N) bool {
				z := secp256k1.ZeroSecp256k1N()
				return z.IsZero()
			}, nil)
		})

		It("Should correctly identify the one element", func() {
			quick.Check(func(x secp256k1.Secp256k1N) bool {
				z := secp256k1.OneSecp256k1N()
				return z.IsOne()
			}, nil)
		})
	})
})
