#include "secp256k1n_asm.h"

typedef struct secp256k1n {
	uint64_t n[5];
} secp256k1n;

void secp256k1n_mul(secp256k1n *r, const secp256k1n *a, const secp256k1n * b) {
	secp256k1n_mul_inner(r->n, a->n, b->n);
}

void secp256k1n_sqr(secp256k1n *r, const secp256k1n *a) {
	secp256k1n_mul_inner(r->n, a->n, a->n);
}

void secp256k1n_inv(secp256k1n *r, const secp256k1n *a) {
    secp256k1n x2, x3, x4, x6, x8, x14, x28, x56, x112, x126, x127, t1;
    int j;

	/** The binary representation of (n - 2) has blocks of 1s with lengths in
	 * { 1, 2, 3, 4, 6, 8, 127 }. Use an addition chain to calculate 2^n - 1
	 * for each block: [1], [2], [3], [4], [6], [8], 14, 28, 56, 112, 126,
	 * [127]
     */

    secp256k1n_sqr(&x2, a);
    secp256k1n_mul(&x2, &x2, a);

    secp256k1n_sqr(&x3, &x2);
    secp256k1n_mul(&x3, &x3, a);

    secp256k1n_sqr(&x4, &x3);
    secp256k1n_mul(&x4, &x4, a);

    x6 = x4;
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&x6, &x6);
    }
    secp256k1n_mul(&x6, &x6, &x2);

    x8 = x6;
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&x8, &x8);
    }
    secp256k1n_mul(&x8, &x8, &x2);

    x14 = x8;
    for (j=0; j<6; j++) {
        secp256k1n_sqr(&x14, &x14);
    }
    secp256k1n_mul(&x14, &x14, &x6);

    x28 = x14;
    for (j=0; j<14; j++) {
        secp256k1n_sqr(&x28, &x28);
    }
    secp256k1n_mul(&x28, &x28, &x14);

    x56 = x28;
    for (j=0; j<28; j++) {
        secp256k1n_sqr(&x56, &x56);
    }
    secp256k1n_mul(&x56, &x56, &x28);

    x112 = x56;
    for (j=0; j<56; j++) {
        secp256k1n_sqr(&x112, &x112);
    }
    secp256k1n_mul(&x112, &x112, &x56);

    x126 = x112;
    for (j=0; j<14; j++) {
        secp256k1n_sqr(&x126, &x126);
    }
    secp256k1n_mul(&x126, &x126, &x14);

    secp256k1n_sqr(&x127, &x126);
    secp256k1n_mul(&x127, &x127, a);

    /* The final result is then assembled using a sliding window over the blocks. */

    t1 = x127;
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<4; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x3);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<4; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x3);
    for (j=0; j<3; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x2);
    for (j=0; j<4; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x3);
    for (j=0; j<5; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x3);
    for (j=0; j<4; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x2);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<5; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x4);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<3; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<4; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<10; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x3);
    for (j=0; j<4; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x3);
    for (j=0; j<9; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x8);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<3; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<3; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<5; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x4);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<5; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x2);
    for (j=0; j<4; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x2);
    for (j=0; j<2; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<8; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x2);
    for (j=0; j<3; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, &x2);
    for (j=0; j<3; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<6; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(&t1, &t1, a);
    for (j=0; j<8; j++) {
        secp256k1n_sqr(&t1, &t1);
    }
    secp256k1n_mul(r, &t1, &x6);
}

