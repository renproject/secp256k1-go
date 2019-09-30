#include <stdint.h>

static void secp256k1n_mul(uint64_t *r, const uint64_t *a, const uint64_t *b) {
/**
 * Registers: rdx:rax = multiplication accumulator
 *            r9:r8   = c
 *            r15:rcx = d
 *            r10-r14 = a0-a4
 *            rbx     = b
 *            rdi     = r
 *            rsi     = a / t?
 */
  uint64_t tmp1, tmp2, tmp3, cl5, cu5, cl6, cu6, cl7, cu7, cl8, cu8;
__asm__ __volatile__(
    "movq 0(%%rsi),%%r10\n"
    "movq 8(%%rsi),%%r11\n"
    "movq 16(%%rsi),%%r12\n"
    "movq 24(%%rsi),%%r13\n"
    "movq 32(%%rsi),%%r14\n"

	/* c = a1 * b4 */
    "movq 32(%%rbx),%%rax\n"
    "mulq %%r11\n"
    "movq %%rax,%%r8\n"
    "movq %%rdx,%%r9\n"
	/* c += a2 * b3 */
    "movq 24(%%rbx),%%rax\n"
    "mulq %%r12\n"
    "addq %%rax,%%r8\n"
    "adcq %%rdx,%%r9\n"
	/* c += a3 * b2 */
    "movq 16(%%rbx),%%rax\n"
    "mulq %%r13\n"
    "addq %%rax,%%r8\n"
    "adcq %%rdx,%%r9\n"
	/* c += a4 * b1 */
    "movq 8(%%rbx),%%rax\n"
    "mulq %%r14\n"
    "addq %%rax,%%r8\n"
    "adcq %%rdx,%%r9\n"
	/* c5 = c */
	"movq %%r8,%q4\n"
	"movq %%r9,%q5\n"
	/* c = a2 * b4 */
    "movq 32(%%rbx),%%rax\n"
    "mulq %%r12\n"
    "movq %%rax,%%r8\n"
    "movq %%rdx,%%r9\n"
	/* c += a3 * b3 */
    "movq 24(%%rbx),%%rax\n"
    "mulq %%r13\n"
    "addq %%rax,%%r8\n"
    "adcq %%rdx,%%r9\n"
	/* c += a4 * b2 */
    "movq 16(%%rbx),%%rax\n"
    "mulq %%r14\n"
    "addq %%rax,%%r8\n"
    "adcq %%rdx,%%r9\n"
	/* c6 = c */
	"movq %%r8,%q6\n"
	"movq %%r9,%q7\n"
	/* c = a3 * b4 */
    "movq 32(%%rbx),%%rax\n"
    "mulq %%r13\n"
    "movq %%rax,%%r8\n"
    "movq %%rdx,%%r9\n"
	/* c += a4 * b3 */
    "movq 24(%%rbx),%%rax\n"
    "mulq %%r14\n"
    "addq %%rax,%%r8\n"
    "adcq %%rdx,%%r9\n"
	/* c7 = c */
	"movq %%r8,%q8\n"
	"movq %%r9,%q9\n"
	/* c8 = a4 * b4 */
    "movq 32(%%rbx),%%rax\n"
    "mulq %%r13\n"
    "movq %%rax,%q10\n"
    "movq %%rdx,%q11\n"

    /* d = a0 * b3 */
    "movq 24(%%rbx),%%rax\n"
    "mulq %%r10\n"
    "movq %%rax,%%rcx\n"
    "movq %%rdx,%%r15\n"
    /* d += a1 * b2 */
    "movq 16(%%rbx),%%rax\n"
    "mulq %%r11\n"
    "movq %%rax,%%rcx\n"
    "movq %%rdx,%%r15\n"
    /* d += a2 * b1 */
    "movq 8(%%rbx),%%rax\n"
    "mulq %%r12\n"
    "movq %%rax,%%rcx\n"
    "movq %%rdx,%%r15\n"
    /* d += a3 * b0 */
    "movq 0(%%rbx),%%rax\n"
    "mulq %%r13\n"
    "movq %%rax,%%rcx\n"
    "movq %%rdx,%%r15\n"
    /* d += (c6 & M) * (r2<<4) */
    "movq %q6,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rax\n"
    "movq $0x14551230,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c7 & M) * (r1<<4) */
    "movq %q8,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rax\n"
    "movq $0x1950b75fc44020,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c8 & M) * (r0<<4) */
    "movq %q10,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rax\n"
    "movq $0xda1732fc9bebf0,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* r[3] = d & M (partial result) */
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rax\n"
	"movq %%rax,24(%%rdi)\n"

    /* d >>= 52 */
    "shrdq $52,%%r15,%%rcx\n"
    "xorq %%r15,%%r15\n"
	/* c = c6 >> 52 */
	"movq %q6,%%r8\n"
	"movq %q7,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (r2<<4) */
    "movq %%r8,%%rax\n"
    "movq $0x14551230,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c7 >> 52 */
	"movq %q8,%%r8\n"
	"movq %q9,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (r1<<4) */
    "movq %%r8,%%rax\n"
    "movq $0x1950b75fc44020,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (r0<<4) */
    "movq %%r8,%%rax\n"
    "movq $0xda1732fc9bebf0,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a4 * b0 */
    "movq 0(%%rbx),%%rax\n"
    "mulq %%r14\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a3 * b1 */
    "movq 8(%%rbx),%%rax\n"
    "mulq %%r13\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a2 * b2 */
    "movq 16(%%rbx),%%rax\n"
    "mulq %%r12\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a1 * b3 */
    "movq 24(%%rbx),%%rax\n"
    "mulq %%r11\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a0 * b4 */
    "movq 32(%%rbx),%%rax\n"
    "mulq %%r10\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c7 & M) * (r2<<4) */
    "movq %q8,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rax\n"
    "movq $0x14551230,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c8 & M) * (r1<<4) */
    "movq %q10,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rax\n"
    "movq $0x1950b75fc44020,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* r[4] = d & M (partial result) */
	"movq %%rcx,%%rsi\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rsi\n"
    "movq %%rsi,32(%%rdi)\n"

    /* d >>= 52 */
    "shrdq $52,%%r15,%%rcx\n"
    "xorq %%r15,%%r15\n"
	/* d *= r0 */
	"movq $0xda1732fc9bebf,%%rax\n"
	"mulq %%rcx\n"
    "movq %%rax,%%rcx\n"
    "movq %%rdx,%%r15\n"
	/* c = c7 >> 52 */
	"movq %q8,%%r8\n"
	"movq %q9,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (((r0<<4)*(r2<<4)) & M) */
    "movq %%r8,%%rax\n"
    "movq $0xe2ffd866a831d00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (((r0<<4)*(r1<<4)) & M) */
    "movq %%r8,%%rax\n"
    "movq $0x777920542397e00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a0 * b0 */
    "movq 0(%%rbx),%%rax\n"
    "mulq %%r10\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c5 & M) * (r0<<4) */
    "movq %q4,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rdx,%%rax\n"
    "movq $0xda1732fc9bebf0,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c8 & M) * ((r0<<4)*(r2<<4)) & M) */
    "movq %q10,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0xe2ffd866a831d00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* r[0] = d & M (complete) */
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rcx,%%rdx\n"
	"movq %%rdx,0(%%rdi)\n"

	/* d <<= 52 */
    "shrdq $52,%%r15,%%rcx\n"
    "xorq %%r15,%%r15\n"
	/* c = c7 >> 52 */
	"movq %q8,%%r8\n"
	"movq %q9,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (((r0<<4)*(r2<<4)) >> 52) */
    "movq %%r8,%%rax\n"
    "movq $0x115249200,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += c * (((r1<<4)*(r2<<4)) & M) */
    "movq %%r8,%%rax\n"
    "movq $0xcca28498bee4600,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (((r0<<4)*(r1<<4)) >> 52) */
    "movq %%r8,%%rax\n"
    "movq $0x15910772c569a00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += c * (((r1<<4)*(r1<<4)) & M) */
    "movq %%r8,%%rax\n"
    "movq $0xcbaebca01100400,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c5 >> 52 */
	"movq %q4,%%r8\n"
	"movq %q5,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * ((r0<<4) & M) */
    "movq %%r8,%%rax\n"
    "movq $0xda1732fc9bebf0,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c8 & M) * ((r0<<4)*(r2<<4)) >> 52) */
    "movq %q10,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0x115249200,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * ((r0<<4)*(r2<<4) & M) */
    "movq %%r8,%%rax\n"
    "movq $0xe2ffd866a831d00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a0 * b1 */
    "movq 8(%%rbx),%%rax\n"
    "mulq %%r10\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a1 * b0 */
    "movq 0(%%rbx),%%rax\n"
    "mulq %%r11\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* d += (c5 & M) * (r1<<4) */
    "movq %q4,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0x1950b75fc44020,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* d += (c6 & M) * (r0<<4) */
    "movq %q6,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0xda1732fc9bebf0,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* d += (c8 & M) * ((r1<<4)*(r2<<4) & M) */
    "movq %q10,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0xcca28498bee4600,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* r[1] = d & M (complete) */
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rcx,%%rdx\n"
	"movq %%rdx,8(%%rdi)\n"

	/* d <<= 52 */
    "shrdq $52,%%r15,%%rcx\n"
    "xorq %%r15,%%r15\n"
	/* c = c7 >> 52 */
	"movq %q8,%%r8\n"
	"movq %q9,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (((r1<<4)*(r2<<4)) >> 52) */
    "movq %%r8,%%rax\n"
    "movq $0x202b7e00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += c * (r2<<4)*(r2<<4) */
    "movq %%r8,%%rax\n"
    "movq $0x19d671c952ac900,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (((r1<<4)*(r1<<4)) >> 52) */
    "movq %%r8,%%rax\n"
    "movq $0x280dd43d389300,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += c * (((r1<<4)*(r2<<4)) & M) */
    "movq %%r8,%%rax\n"
    "movq $0xcca28498bee4600,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * ((r0<<4)*(r2<<4) >> 52) */
    "movq %%r8,%%rax\n"
    "movq $0x115249200,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c5 >> 52 */
	"movq %q4,%%r8\n"
	"movq %q5,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (r1<<4) */
    "movq %%r8,%%rax\n"
    "movq $0x1950b75fc44020,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c6 >> 52 */
	"movq %q6,%%r8\n"
	"movq %q7,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (r0<<4) */
    "movq %%r8,%%rax\n"
    "movq $0xda1732fc9bebf0,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c8 & M) * ((r1<<4)*(r2<<4)) >> 52) */
    "movq %q10,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0x202b7e00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * ((r1<<4)*(r2<<4) & M) */
    "movq %%r8,%%rax\n"
    "movq $0xcca28498bee4600,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a0 * b2 */
    "movq 16(%%rbx),%%rax\n"
    "mulq %%r10\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a1 * b1 */
    "movq 8(%%rbx),%%rax\n"
    "mulq %%r11\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += a2 * b0 */
    "movq 0(%%rbx),%%rax\n"
    "mulq %%r12\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c5 & M) * (r2<<4) */
    "movq %q4,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0x14551230,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c6 & M) * (r1<<4) */
    "movq %q6,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0x1950b75fc44020,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c7 & M) * (r0<<4) */
    "movq %q8,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0xda1732fc9bebf0,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += (c8 & M) * ((r2*r2)<<8) */
    "movq %q10,%%rax\n"
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rax,%%rdx\n"
    "movq $0x19d671c952ac900,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* r[2] = d & M (complete) */
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rcx,%%rdx\n"
	"movq %%rdx,16(%%rdi)\n"

	/* d <<= 52 */
    "shrdq $52,%%r15,%%rcx\n"
    "xorq %%r15,%%r15\n"
	/* c = c8 >> 52 */
	"movq %q10,%%r8\n"
	"movq %q11,%%r9\n"
    "shrdq $52,%%r9,%%r8\n"
    /* d += c * (((r1<<4)*(r2<<4)) >> 52) */
    "movq %%r8,%%rax\n"
    "movq $0x202b7e00,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
    /* d += c * (((r2<<4)*(r2<<4)) >> 52) */
    "movq %%r8,%%rax\n"
    "movq $0x19d671c952ac900,%%rdx\n"
    "mulq %%rdx\n"
    "addq %%rax,%%rcx\n"
    "adcq %%rdx,%%r15\n"
	/* r[3] += d & M (complete) */
    "movq $0xfffffffffffff,%%rdx\n"
    "andq %%rcx,%%rdx\n"
    "addq %%rdx,24(%%rdi)\n"

	/* d <<= 52 */
    "shrdq $52,%%r15,%%rcx\n"
    "xorq %%r15,%%r15\n"
	/* r[4] += d (complete) */
    "addq %%rcx,24(%%rdi)\n"
: "+S"(a), "=m"(tmp1), "=m"(tmp2), "=m"(tmp3), "=m"(cl5), "=m"(cu5), "=m"(cl6), "=m"(cu6), "=m"(cl7), "=m"(cu7), "=m"(cl8), "=m"(cl8)
: "b"(b), "D"(r)
: "%rax", "%rcx", "%rdx", "%r8", "%r9", "%r10", "%r11", "%r12", "%r13", "%r14", "%r15", "cc", "memory"
);
}
