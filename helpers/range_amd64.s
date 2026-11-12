//go:build !purego && !race

#include "textflag.h"

TEXT ·indexByteRangeSIMD(SB), NOSPLIT, $0-32
	MOVQ in+0(FP), SI
	MOVQ n+8(FP), CX
	MOVL first+16(FP), AX
	MOVL last+20(FP), BX
	XORQ DI, DI
	IMULL $0x01010101, AX
	IMULL $0x01010101, BX
	MOVD AX, X0
	PSHUFD $0, X0, X0
	MOVD BX, X1
	PSHUFD $0, X1, X1
vector:
	CMPQ CX, $16
	JL scalar
	MOVOU (SI), X2
	MOVOU X2, X3
	PMAXUB X0, X2
	PCMPEQB X3, X2
	PMINUB X1, X3
	MOVOU (SI), X4
	PCMPEQB X4, X3
	PAND X3, X2
	PMOVMSKB X2, DX
	TESTL DX, DX
	JNZ scalar
	ADDQ $16, SI
	ADDQ $16, DI
	SUBQ $16, CX
	JMP vector
scalar:
	TESTQ CX, CX
	JZ miss
	MOVBLZX (SI), DX
	CMPL DX, first+16(FP)
	JL next
	CMPL DX, last+20(FP)
	JG next
	MOVQ DI, ret+24(FP)
	RET
next:
	ADDQ $1, SI
	INCQ DI
	DECQ CX
	JMP scalar
miss:
	MOVQ $-1, ret+24(FP)
	RET

TEXT ·indexRuneRangeSIMD(SB), NOSPLIT, $0-32
	MOVQ in+0(FP), SI
	MOVQ n+8(FP), CX
	MOVL first+16(FP), AX
	MOVL last+20(FP), BX
	XORQ DI, DI
	MOVD AX, X0
	PSHUFD $0, X0, X0
	MOVD BX, X1
	PSHUFD $0, X1, X1
vector:
	CMPQ CX, $4
	JL scalar
	MOVOU (SI), X2
	MOVOU X2, X3
	MOVOU X0, X4
	PCMPGTL X2, X4
	PCMPGTL X1, X3
	POR X3, X4
	PMOVMSKB X4, DX
	CMPL DX, $65535
	JNE scalar
	ADDQ $16, SI
	ADDQ $4, DI
	SUBQ $4, CX
	JMP vector
scalar:
	TESTQ CX, CX
	JZ miss
	MOVL (SI), DX
	CMPL DX, first+16(FP)
	JL next
	CMPL DX, last+20(FP)
	JG next
	MOVQ DI, ret+24(FP)
	RET
next:
	ADDQ $4, SI
	INCQ DI
	DECQ CX
	JMP scalar
miss:
	MOVQ $-1, ret+24(FP)
	RET
