package secp256k1_test

import (
	"math/big"
	"testing/quick"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/renproject/secp256k1-go"
)

var _ = Describe("Wrapped field elements", func() {
	P, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	N, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)

	Context("When using Fp field elements", func() {
		It("Should normalise correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				// Make sure normalisation won't be trivial
				x.MulInt(10)

				z := secp256k1.NewSecp256k1P(0)
				z.Set(&x)
				z.Normalize()

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should normalise (weak) correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				// Make sure normalisation won't be trivial
				x.MulInt(10)

				z := secp256k1.NewSecp256k1P(0)
				z.Set(&x)
				z.NormalizeWeak()

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should normalise (var) correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				// Make sure normalisation won't be trivial
				x.MulInt(10)

				z := secp256k1.NewSecp256k1P(0)
				z.Set(&x)
				z.NormalizeVar()

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should check equality to zero correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				// Make sure normalisation won't be trivial
				z := secp256k1.NewSecp256k1P(0)
				z.Neg(&x, 2)
				z.Add(&x)

				return z.NormalizesToZero()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should check equality to zero (var) correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				// Make sure normalisation won't be trivial
				z := secp256k1.NewSecp256k1P(0)
				z.Neg(&x, 2)
				z.Add(&x)

				return z.NormalizesToZeroVar()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should set a uint64 correctly", func() {
			err := quick.Check(func(x uint64) bool {
				// Make sure normalisation won't be trivial
				z := secp256k1.NewSecp256k1P(0)
				z.SetUint64(x)
				z.NormalizeVar()

				return z.Int().Uint64() == x
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should clear correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				x.Clear()

				return x.NormalizesToZeroVar()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should check that a field element is zero correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				x.Clear()

				return x.IsZero()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly identify odd elements", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				return (x.Int().Uint64()%2 == 1) == x.IsOdd()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should check equality correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				var z, inv secp256k1.Secp256k1P
				z.Set(&x)

				// Try to get a different limbed representation
				inv.Inv(&z)
				z.Sqr(&z)
				z.Mul(&z, &inv)

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should check equality (var) correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				var z, inv secp256k1.Secp256k1P
				z.Set(&x)

				// Try to get a different limbed representation
				inv.Inv(&z)
				z.Sqr(&z)
				z.Mul(&z, &inv)

				return z.EqVar(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should perform comparisons correctly", func() {
			err := quick.Check(func(x, y secp256k1.Secp256k1P) bool {
				x.NormalizeVar()
				y.NormalizeVar()

				return x.Int().Cmp(y.Int()) == x.CmpVar(&y)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly square field elements", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				var y, z secp256k1.Secp256k1P
				y.Sqr(&x)
				z.Mul(&x, &x)

				return z.Eq(&y)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly compute square roots", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				var z secp256k1.Secp256k1P

				// If the square root exists, it should have been correctly
				// computed
				if z.Sqrt(&x) {
					z.Sqr(&z)
					return z.Eq(&x)
				}

				return true
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly determine if a field element is a quadratic residue", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				x.Sqr(&x)

				return x.IsQuadVar()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly perform a conditional copy", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P, flag bool) bool {
				var z secp256k1.Secp256k1P
				z.Clear()

				z.Cmov(&x, flag)
				if flag {
					return z.Eq(&x)
				} else {
					return z.IsZero()
				}
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should add correctly", func() {
			err := quick.Check(func(x, y secp256k1.Secp256k1P) bool {
				sum := big.NewInt(0).Add(x.Int(), y.Int())
				sum.Mod(sum, P)
				z := secp256k1.NewSecp256k1P(0)

				z.Set(&x)
				z.Add(&y)
				z.NormalizeVar()

				return z.Int().Cmp(sum) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should negate correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				neg := big.NewInt(0).Sub(P, big.NewInt(0).Mod(x.Int(), P))
				z := secp256k1.NewSecp256k1P(0)

				z.Neg(&x, 0)
				z.NormalizeVar()

				return z.Int().Cmp(neg) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should multiply correctly", func() {
			err := quick.Check(func(x, y secp256k1.Secp256k1P) bool {
				prod := big.NewInt(0).Mul(x.Int(), y.Int())
				prod.Mod(prod, P)
				z := secp256k1.NewSecp256k1P(0)

				z.Mul(&x, &y)
				z.NormalizeVar()

				return z.Int().Cmp(prod) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should invert correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				inv := big.NewInt(0).ModInverse(x.Int(), P)
				z := secp256k1.NewSecp256k1P(0)

				z.Inv(&x)
				z.NormalizeVar()

				return z.Int().Cmp(inv) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should invert (var) correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				inv := big.NewInt(0).ModInverse(x.Int(), P)
				z := secp256k1.NewSecp256k1P(0)

				z.InvVar(&x)
				z.NormalizeVar()

				return z.Int().Cmp(inv) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly construct the zero element", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				var z secp256k1.Secp256k1P
				zero := secp256k1.ZeroSecp256k1P()
				z.Set(&x)

				// The defining property of the zero element is that it is the
				// additive identity
				z.Add(&zero)

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly construct the one element", func() {
			err := quick.Check(func(x secp256k1.Secp256k1P) bool {
				var z secp256k1.Secp256k1P
				one := secp256k1.OneSecp256k1P()

				// The defining property of the zero element is that it is the
				// multiplicative identity
				z.Mul(&x, &one)

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should convert from a big.Int correctly", func() {
			err := quick.Check(func(b [32]byte) bool {
				y := big.NewInt(0).SetBytes(b[:])
				z := secp256k1.Secp256k1PFromInt(y)
				return z.Int().Cmp(y) == 0
			}, nil)
			Expect(err).To(BeNil())
		})
	})

	Context("When using Fn field elements", func() {
		It("Should add correctly", func() {
			err := quick.Check(func(x, y secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				sum := big.NewInt(0).Add(x.Int(), y.Int())
				sum.Mod(sum, N)
				z.Add(&x, &y)
				z.Normalize()
				return z.Int().Cmp(sum) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should negate correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				var y secp256k1.Secp256k1N
				neg := x.Int()
				neg.Mod(neg, N)
				neg.Sub(N, neg)

				y.Neg(&x, 0)
				y.Normalize()

				return y.Int().Cmp(neg) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should multiply correctly", func() {
			err := quick.Check(func(x, y secp256k1.Secp256k1N) bool {
				z := secp256k1.NewSecp256k1N(0)
				prod := big.NewInt(0).Mul(x.Int(), y.Int())
				prod.Mod(prod, N)

				z.Mul(&x, &y)
				z.Normalize()

				return z.Int().Cmp(prod) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should invert correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				inv := big.NewInt(0).ModInverse(x.Int(), N)

				z.Inv(&x)
				z.Normalize()

				return z.Int().Cmp(inv) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should check equality correctly", func() {
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z, inv secp256k1.Secp256k1N
				z.Set(&x)

				// Try to get a different limbed representation
				inv.Inv(&z)
				z.Sqr(&z)
				z.Mul(&z, &inv)

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly construct the zero element", func() {
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				zero := secp256k1.ZeroSecp256k1N()

				// The defining property of the zero element is that it is the
				// additive identity
				z.Add(&x, &zero)

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly construct the one element", func() {
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				var z secp256k1.Secp256k1N
				one := secp256k1.OneSecp256k1N()

				// The defining property of the zero element is that it is the
				// multiplicative identity
				z.Mul(&x, &one)

				return z.Eq(&x)
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly identify the zero element", func() {
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				z := secp256k1.ZeroSecp256k1N()
				return z.IsZero()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should correctly identify the one element", func() {
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				z := secp256k1.OneSecp256k1N()
				return z.IsOne()
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should set a value from bytes correctly", func() {
			err := quick.Check(func(b [32]byte) bool {
				var z secp256k1.Secp256k1N
				y := big.NewInt(0).SetBytes(b[:])
				z.SetB32(b[:])
				return z.Int().Cmp(y) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should return the bytes for a value correctly", func() {
			var b [32]byte
			err := quick.Check(func(x secp256k1.Secp256k1N) bool {
				x.GetB32(b[:])
				y := big.NewInt(0).SetBytes(b[:])
				return x.Int().Cmp(y) == 0
			}, nil)
			Expect(err).To(BeNil())
		})

		It("Should convert from a big.Int correctly", func() {
			err := quick.Check(func(b [32]byte) bool {
				y := big.NewInt(0).SetBytes(b[:])
				z := secp256k1.Secp256k1NFromInt(y)
				return z.Int().Cmp(y) == 0
			}, nil)
			Expect(err).To(BeNil())
		})
	})
})
