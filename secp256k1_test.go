package secp256k1_test

import (
	"crypto/rand"
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/republicprotocol/secp256k1-go"
)

var _ = Describe("Wrapped field elements", func() {
	P, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F", 16)
	N, _ := big.NewInt(0).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	TRIALS := 10000

	randomBig256 := func() *big.Int {
		b := make([]byte, 32)
		rand.Read(b)
		return big.NewInt(0).SetBytes(b)
	}

	paddedBytes := func(b []byte) []byte {
		ret := make([]byte, 32)
		copy(ret[32-len(b):], b)
		return ret
	}

	Context("When doing arithmetic in Fp", func() {
		It("Should add correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x, y := randomBig256(), randomBig256()
				sum := big.NewInt(0).Add(x, y)
				sum.Mod(sum, P)
				a := secp256k1.NewSecp256k1P(0)
				b := secp256k1.NewSecp256k1P(0)
				c := secp256k1.NewSecp256k1P(0)
				a.SetB32(paddedBytes(x.Bytes()))
				b.SetB32(paddedBytes(y.Bytes()))

				c.Set(&a)
				c.Add(&b)
				c.NormalizeVar()

				Expect(c.BigInt().Cmp(sum)).To(Equal(0))
			}
		})

		It("Should negate correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x := randomBig256()
				neg := big.NewInt(0).Sub(P, x)
				a := secp256k1.NewSecp256k1P(0)
				b := secp256k1.NewSecp256k1P(0)
				a.SetB32(paddedBytes(x.Bytes()))

				b.Negate(&a, 0)
				b.NormalizeVar()

				Expect(b.BigInt().Cmp(neg)).To(Equal(0))
			}
		})

		It("Should multiply correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x, y := randomBig256(), randomBig256()
				prod := big.NewInt(0).Mul(x, y)
				prod.Mod(prod, P)
				a := secp256k1.NewSecp256k1P(0)
				b := secp256k1.NewSecp256k1P(0)
				c := secp256k1.NewSecp256k1P(0)

				a.SetB32(paddedBytes(x.Bytes()))
				b.SetB32(paddedBytes(y.Bytes()))
				c.Mul(&a, &b)
				c.NormalizeVar()

				Expect(c.BigInt().Cmp(prod)).To(Equal(0))
			}
		})

		It("Should invert correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x := randomBig256()
				inv := big.NewInt(0).ModInverse(x, P)
				a := secp256k1.NewSecp256k1P(0)
				b := secp256k1.NewSecp256k1P(0)

				a.SetB32(paddedBytes(x.Bytes()))
				b.Inv(&a)
				b.NormalizeVar()

				Expect(b.BigInt().Cmp(inv)).To(Equal(0))
			}
		})
	})

	Context("When doing arithmetic in Fn", func() {
		It("Should add correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x, y := secp256k1.RandomSecp256k1N(), secp256k1.RandomSecp256k1N()
				z := secp256k1.NewSecp256k1N(0)
				sum := big.NewInt(0).Add(x.Int(), y.Int())
				sum.Mod(sum, N)

				z.Add(&x, &y)
				z.Normalize()

				Expect(z.Int().Cmp(sum)).To(Equal(0))
			}
		})
	})
})
