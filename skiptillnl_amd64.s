#include "textflag.h"

// Реализация на ассемблере для x86-64 с SSE4.2
TEXT ·skipTillNlSSE2(SB),NOSPLIT,$0-32
    MOVQ b_base+0(FP), SI
    MOVQ b_len+8(FP), CX
    XORQ AX, AX
    DECQ AX              // AX = -1 (default return)

    TESTQ CX, CX
    JZ done

    MOVQ SI, DI         // DI = pointer to data
    MOVQ CX, DX         // DX = remaining length

    // Check if we can use SSE
    CMPQ CX, $16
    JLT scalar

    // Align to 16 bytes
    MOVQ DI, BX
    ANDQ $15, BX
    JZ aligned
    MOVQ $16, AX
    SUBQ BX, AX

preloop:
    MOVB (DI), BL
    CMPB BL, $0x0A      // \n
    JE found
    CMPB BL, $0x0D      // \r
    JE found
    INCQ DI
    DECQ DX
    JZ notfound
    DECQ AX
    JNZ preloop

aligned:
    MOVQ DX, BX
    SHRQ $4, BX         // BX = number of 16-byte chunks
    JZ scalar

    // Prepare X0 = \n\n\n\n...
    MOVQ $0x0A0A0A0A0A0A0A0A, AX
    MOVQ AX, X0
    PUNPCKLBW X0, X0
    PUNPCKLBW X0, X0

    // Prepare X1 = \r\r\r\r...
    MOVQ $0x0D0D0D0D0D0D0D0D, AX
    MOVQ AX, X1
    PUNPCKLBW X1, X1
    PUNPCKLBW X1, X1

chunkloop:
    MOVOU (DI), X2      // Load 16 bytes
    MOVO X2, X3

    PCMPEQB X0, X2      // Compare with \n
    PCMPEQB X1, X3      // Compare with \r
    POR X2, X3          // Combine results
    PMOVMSKB X3, R8     // Get bitmask
    TESTL R8, R8
    JNZ found_chunk

    ADDQ $16, DI
    SUBQ $16, DX
    DECQ BX
    JNZ chunkloop

scalar:
    TESTQ DX, DX
    JZ notfound

scalar_loop:
    MOVB (DI), BL
    CMPB BL, $0x0A      // \n
    JE found
    CMPB BL, $0x0D      // \r
    JE found
    INCQ DI
    DECQ DX
    JNZ scalar_loop

notfound:
    XORQ AX, AX
    DECQ AX             // AX = -1
    JMP done

found_chunk:
    BSFQ R8, R8        // Find first set bit
    ADDQ R8, DI        // DI points to match
    SUBQ SI, DI        // DI = index
    MOVQ DI, AX
    JMP done

found:
    SUBQ SI, DI        // DI = index
    MOVQ DI, AX

done:
    MOVQ AX, ret+24(FP)
    RET

// func skipTillNlAVX2(b []byte) int
TEXT ·skipTillNlAVX2(SB),NOSPLIT,$0-32
    MOVQ b_base+0(FP), SI     // SI = pointer to slice data
    MOVQ b_len+8(FP), CX      // CX = length of slice
    XORQ AX, AX               // AX will hold result
    DECQ AX                   // Default result = -1

    TESTQ CX, CX
    JZ done                  // Empty slice case

    MOVQ SI, DI              // DI = current pointer
    MOVQ CX, DX              // DX = remaining bytes

    // Check if we can use AVX2 (need at least 32 bytes)
    CMPQ CX, $32
    JLT scalar

    // Align to 32 bytes
    MOVQ DI, BX
    ANDQ $31, BX
    JZ avx2_aligned
    MOVQ $32, AX
    SUBQ BX, AX

preloop:
    MOVB (DI), BL            // Load byte
    CMPB BL, $0x0A           // Check for \n
    JE found
    CMPB BL, $0x0D           // Check for \r
    JE found
    INCQ DI                  // Move to next byte
    DECQ DX                  // Decrement remaining count
    JZ notfound              // If no bytes left, return -1
    DECQ AX                  // Decrement alignment counter
    JNZ preloop              // Continue until aligned

avx2_aligned:
    MOVQ DX, BX
    SHRQ $5, BX              // BX = number of 32-byte chunks
    JZ avx2_tail

    // Prepare Y0 = \n\n\n\n... (32 bytes of 0x0A)
    MOVQ $0x0A0A0A0A0A0A0A0A, AX
    MOVQ AX, X0
    VPUNPCKLBW X0, X0, X0
    VPUNPCKLBW X0, X0, X0
    VINSERTI128 $1, X0, Y0, Y0

    // Prepare Y1 = \r\r\r\r... (32 bytes of 0x0D)
    MOVQ $0x0D0D0D0D0D0D0D0D, AX
    MOVQ AX, X1
    VPUNPCKLBW X1, X1, X1
    VPUNPCKLBW X1, X1, X1
    VINSERTI128 $1, X1, Y1, Y1

avx2_loop:
    VMOVDQU (DI), Y2         // Load 32 bytes
    VMOVDQU Y2, Y3

    VPCMPEQB Y0, Y2, Y2      // Compare with \n
    VPCMPEQB Y1, Y3, Y3      // Compare with \r
    VPOR Y2, Y3, Y3          // Combine results
    VPMOVMSKB Y3, R8         // Get bitmask
    TESTL R8, R8
    JNZ avx2_found

    ADDQ $32, DI             // Move pointer
    SUBQ $32, DX             // Decrement remaining count
    DECQ BX                  // Decrement chunk counter
    JNZ avx2_loop

avx2_tail:
    // Handle remaining 16-31 bytes with AVX2
    CMPQ DX, $16
    JB scalar

    VMOVDQU (DI), X2
    VMOVDQU X2, X3

    VPCMPEQB X0, X2, X2      // Compare with \n (lower 128 bits)
    VPCMPEQB X1, X3, X3      // Compare with \r (lower 128 bits)
    VPOR X2, X3, X3
    VPMOVMSKB X3, R8
    TESTL R8, R8
    JNZ avx2_found_lower

    ADDQ $16, DI
    SUBQ $16, DX

scalar:
    TESTQ DX, DX
    JZ notfound

scalar_loop:
    MOVB (DI), BL
    CMPB BL, $0x0A           // \n
    JE found
    CMPB BL, $0x0D           // \r
    JE found
    INCQ DI
    DECQ DX
    JNZ scalar_loop

notfound:
    XORQ AX, AX
    DECQ AX                  // AX = -1
    JMP done

avx2_found_lower:
    BSFQ R8, R8              // Find first set bit in lower 128 bits
    JMP avx2_calc_pos

avx2_found:
    VZEROUPPER               // Clear upper bits of YMM registers
    BSFQ R8, R8              // Find first set bit

avx2_calc_pos:
    ADDQ R8, DI              // DI points to match
    SUBQ SI, DI              // DI = index
    MOVQ DI, AX
    JMP done

found:
    SUBQ SI, DI              // DI = index
    MOVQ DI, AX

done:
    MOVQ AX, ret+24(FP)      // Store result
    RET

// func skipTillNlAVX512(b []byte) int
TEXT ·skipTillNlAVX512(SB),NOSPLIT,$0-32
    MOVQ b_base+0(FP), SI     // SI = pointer to slice data
    MOVQ b_len+8(FP), CX      // CX = length of slice
    XORQ AX, AX               // AX will hold result
    DECQ AX                   // Default result = -1

    TESTQ CX, CX
    JZ done                  // Empty slice case

    MOVQ SI, DI              // DI = current pointer
    MOVQ CX, DX              // DX = remaining bytes

    // Check if we can use AVX512 (need at least 64 bytes)
    CMPQ CX, $64
    JLT scalar

    // Prepare constants for comparison
    MOVQ $0x0A0A0A0A0A0A0A0A, AX
    MOVQ AX, X0
    VPUNPCKLBW X0, X0, X0
    VPUNPCKLBW X0, X0, X0
    VPUNPCKLBW X0, X0, X0
    VPBROADCASTB X0, Z0      // Z0 = 64 bytes of 0x0A (\n)

    MOVQ $0x0D0D0D0D0D0D0D0D, AX
    MOVQ AX, X1
    VPUNPCKLBW X1, X1, X1
    VPUNPCKLBW X1, X1, X1
    VPUNPCKLBW X1, X1, X1
    VPBROADCASTB X1, Z1      // Z1 = 64 bytes of 0x0D (\r)

avx512_loop:
    VMOVDQU64 (DI), Z2       // Load 64 bytes
    VMOVDQU64 Z2, Z3

    VPCMPEQB Z0, Z2, K1      // Compare with \n
    VPCMPEQB Z1, Z3, K2      // Compare with \r
    KORQ K1, K2, K3          // Combine results
    KTESTQ K3, K3
    JNZ avx512_found

    ADDQ $64, DI             // Advance pointer
    SUBQ $64, DX             // Decrement remaining count
    CMPQ DX, $64
    JGE avx512_loop          // Continue if enough bytes remain

scalar:
    TESTQ DX, DX
    JZ notfound

scalar_loop:
    MOVB (DI), BL
    CMPB BL, $0x0A           // \n
    JE found
    CMPB BL, $0x0D           // \r
    JE found
    INCQ DI
    DECQ DX
    JNZ scalar_loop

notfound:
    XORQ AX, AX
    DECQ AX                  // AX = -1
    JMP done

avx512_found:
    KMOVQ K3, R8             // Move mask to general register
    BSFQ R8, R8              // Find first set bit
    ADDQ R8, DI              // DI points to match
    SUBQ SI, DI              // DI = index
    MOVQ DI, AX
    JMP done

found:
    SUBQ SI, DI              // DI = index
    MOVQ DI, AX

done:
    VZEROUPPER               // Clear upper bits
    MOVQ AX, ret+24(FP)      // Store result
    RET
