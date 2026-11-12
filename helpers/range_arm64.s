//go:build !purego && !race

#include "textflag.h"

TEXT ·indexByteRangeSIMD(SB), NOSPLIT, $0-32
	MOVD in+0(FP), R0
	MOVD n+8(FP), R1
	MOVW first+16(FP), R2
	MOVW last+20(FP), R3
	MOVD $0, R4
	VDUP R2, V0.B16
	VDUP R3, V1.B16
wide:
	CMP $32, R1
	BLT vector
	VLD1 (R0), [V2.B16]
	VCMHS V0.B16, V2.B16, V3.B16
	VCMHS V2.B16, V1.B16, V4.B16
	VAND V3.B16, V4.B16, V5.B16
	ADD $16, R0, R8
	VLD1 (R8), [V2.B16]
	VCMHS V0.B16, V2.B16, V3.B16
	VCMHS V2.B16, V1.B16, V4.B16
	VAND V3.B16, V4.B16, V7.B16
	VORR V5.B16, V7.B16, V8.B16
	VUMAXV V8.B16, V6
	VMOV V6.S[0], R5
	CBNZ R5, widefound
	ADD $32, R0
	ADD $32, R4
	SUB $32, R1
	B wide
vector:
	CMP $16, R1
	BLT small
	VLD1 (R0), [V2.B16]
	VCMHS V0.B16, V2.B16, V3.B16
	VCMHS V2.B16, V1.B16, V4.B16
	VAND V3.B16, V4.B16, V5.B16
	VUMAXV V5.B16, V6
	VMOV V6.S[0], R5
	CBNZ R5, found
	ADD $16, R0
	ADD $16, R4
	SUB $16, R1
	B vector
small:
	CMP $8, R1
	BLT scalar
	VLD1 (R0), [V2.B8]
	VCMHS V0.B8, V2.B8, V3.B8
	VCMHS V2.B8, V1.B8, V4.B8
	VAND V3.B8, V4.B8, V5.B8
	VUMAXV V5.B8, V6
	VMOV V6.S[0], R5
	CBNZ R5, found
	ADD $8, R0
	ADD $8, R4
	SUB $8, R1
scalar:
	CBZ R1, miss
	MOVBU (R0), R5
	CMP R2, R5
	BLO next
	CMP R3, R5
	BHI next
	MOVD R4, ret+24(FP)
	RET
next:
	ADD $1, R0
	ADD $1, R4
	SUB $1, R1
	B scalar
miss:
	MOVD $-1, R4
	MOVD R4, ret+24(FP)
	RET

found:
	VMOV V5.D[0], R5
	VMOV V5.D[1], R6
	CBNZ R5, firsthalf
	MOVD R6, R5
	ADD $8, R4
firsthalf:
	RBIT R5, R5
	CLZ R5, R5
	LSR $3, R5
	ADD R5, R4
	MOVD R4, ret+24(FP)
	RET
widefound:
	VUMAXV V5.B16, V6
	VMOV V6.S[0], R5
	CBNZ R5, found
	VMOV V7.B16, V5.B16
	ADD $16, R4
	B found
TEXT ·indexRuneRangeSIMD(SB), NOSPLIT, $0-32
	MOVD in+0(FP), R0
	MOVD n+8(FP), R1
	MOVW first+16(FP), R2
	MOVW last+20(FP), R3
	MOVD $0, R4
	VDUP R2, V0.S4
	VDUP R3, V1.S4
wide:
	CMP $8, R1
	BLT vector
	VLD1 (R0), [V2.S4]
	VCMGE V0.S4, V2.S4, V3.S4
	VCMGE V2.S4, V1.S4, V4.S4
	VAND V3.B16, V4.B16, V5.B16
	ADD $16, R0, R8
	VLD1 (R8), [V2.S4]
	VCMGE V0.S4, V2.S4, V3.S4
	VCMGE V2.S4, V1.S4, V4.S4
	VAND V3.B16, V4.B16, V7.B16
	VORR V5.B16, V7.B16, V8.B16
	VUMAXV V8.S4, V6
	VMOV V6.S[0], R5
	CBNZ R5, widefound
	ADD $32, R0
	ADD $8, R4
	SUB $8, R1
	B wide
vector:
	CMP $4, R1
	BLT scalar
	VLD1 (R0), [V2.S4]
	VCMGE V0.S4, V2.S4, V3.S4
	VCMGE V2.S4, V1.S4, V4.S4
	VAND V3.B16, V4.B16, V5.B16
	VUMAXV V5.S4, V6
	VMOV V6.S[0], R5
	CBNZ R5, found
	ADD $16, R0
	ADD $4, R4
	SUB $4, R1
	B vector
scalar:
	CBZ R1, miss
	MOVW (R0), R5
	CMP R2, R5
	BLT next
	CMP R3, R5
	BGT next
	MOVD R4, ret+24(FP)
	RET
next:
	ADD $4, R0
	ADD $1, R4
	SUB $1, R1
	B scalar
miss:
	MOVD $-1, R4
	MOVD R4, ret+24(FP)
	RET
found:
	VMOV V5.D[0], R5
	VMOV V5.D[1], R6
	CBNZ R5, firsthalf
	MOVD R6, R5
	ADD $2, R4
firsthalf:
	RBIT R5, R5
	CLZ R5, R5
	LSR $5, R5
	ADD R5, R4
	MOVD R4, ret+24(FP)
	RET
widefound:
	VUMAXV V5.S4, V6
	VMOV V6.S[0], R5
	CBNZ R5, found
	VMOV V7.B16, V5.B16
	ADD $4, R4
	B found
