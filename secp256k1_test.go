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
	TRIALS := 10000

	randomBig256 := func() *big.Int {
		b := make([]byte, 32)
		rand.Read(b)
		return big.NewInt(0).SetBytes(b)
	}

	Context("When doing arithmetic in Fp", func() {
		It("Should add correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x, y := randomBig256(), randomBig256()
				sum := big.NewInt(0).Add(x, y)
				sum.Mod(sum, P)
				a := secp256k1.NewFp(0)
				b := secp256k1.NewFp(0)
				c := secp256k1.NewFp(0)
				a.SetBytes(x.Bytes())
				b.SetBytes(y.Bytes())
				c.Add(&a, &b)
				Expect(c.BigInt().Cmp(sum)).To(Equal(0))
			}
		})

		It("Should negate correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x := randomBig256()
				neg := big.NewInt(0).Sub(P, x)
				a := secp256k1.NewFp(0)
				b := secp256k1.NewFp(0)
				a.SetBytes(x.Bytes())
				b.Neg(&a)
				Expect(b.BigInt().Cmp(neg)).To(Equal(0))
			}
		})

		It("Should multiply correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x, y := randomBig256(), randomBig256()
				prod := big.NewInt(0).Mul(x, y)
				prod.Mod(prod, P)
				a := secp256k1.NewFp(0)
				b := secp256k1.NewFp(0)
				c := secp256k1.NewFp(0)
				a.SetBytes(x.Bytes())
				b.SetBytes(y.Bytes())
				c.Mul(&a, &b)
				Expect(c.BigInt().Cmp(prod)).To(Equal(0))
			}
		})

		It("Should invert correctly", func() {
			for i := 0; i < TRIALS; i++ {
				x := randomBig256()
				inv := big.NewInt(0).ModInverse(x, P)
				a := secp256k1.NewFp(0)
				b := secp256k1.NewFp(0)
				a.SetBytes(x.Bytes())
				b.Inv(&a)
				Expect(b.BigInt().Cmp(inv)).To(Equal(0))
			}
		})
	})
})
