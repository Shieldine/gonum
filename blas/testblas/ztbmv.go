// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
)

type Ztbmver interface {
	Ztbmv(uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, k int, ab []complex128, ldab int, x []complex128, incX int)

	Ztrmver
}

func ZtbmvTest(t *testing.T, impl Ztbmver) {
	rnd := rand.New(rand.NewPCG(1, 1))
	for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
		for _, trans := range []blas.Transpose{blas.NoTrans, blas.Trans, blas.ConjTrans} {
			for _, diag := range []blas.Diag{blas.NonUnit, blas.Unit} {
				for _, n := range []int{1, 2, 3, 5} {
					for k := 0; k < n; k++ {
						for _, ldab := range []int{k + 1, k + 1 + 10} {
							for _, incX := range []int{-4, 1, 5} {
								testZtbmv(t, impl, rnd, uplo, trans, diag, n, k, ldab, incX)
							}
						}
					}
				}
			}
		}
	}
}

// testZtbmv tests Ztbmv by comparing its output to that of Ztrmv.
func testZtbmv(t *testing.T, impl Ztbmver, rnd *rand.Rand, uplo blas.Uplo, trans blas.Transpose, diag blas.Diag, n, k, ldab, incX int) {
	const tol = 1e-13

	// Allocate a dense-storage triangular band matrix filled with NaNs that
	// will be used as the reference matrix for Ztrmv.
	lda := max(1, n)
	a := makeZGeneral(nil, n, n, lda)
	// Fill the referenced triangle with random data within the band and
	// with zeros outside.
	if uplo == blas.Upper {
		for i := 0; i < n; i++ {
			for j := i; j < min(n, i+k+1); j++ {
				re := rnd.NormFloat64()
				im := rnd.NormFloat64()
				a[i*lda+j] = complex(re, im)
			}
			for j := i + k + 1; j < n; j++ {
				a[i*lda+j] = 0
			}
		}
	} else {
		for i := 0; i < n; i++ {
			for j := 0; j < i-k; j++ {
				a[i*lda+j] = 0
			}
			for j := max(0, i-k); j <= i; j++ {
				re := rnd.NormFloat64()
				im := rnd.NormFloat64()
				a[i*lda+j] = complex(re, im)
			}
		}
	}
	if diag == blas.Unit {
		// The diagonal should not be referenced by Ztbmv and Ztrmv, so
		// invalidate it with NaNs.
		for i := 0; i < n; i++ {
			a[i*lda+i] = znan
		}
	}
	// Create the triangular band matrix.
	ab := zPackTriBand(k, ldab, uplo, n, a, lda)
	abCopy := make([]complex128, len(ab))
	copy(abCopy, ab)

	// Generate a random complex vector x.
	xtest := make([]complex128, n)
	for i := range xtest {
		re := rnd.NormFloat64()
		im := rnd.NormFloat64()
		xtest[i] = complex(re, im)
	}
	x := makeZVector(xtest, incX)
	xCopy := make([]complex128, len(x))
	copy(xCopy, x)

	want := make([]complex128, len(x))
	copy(want, x)

	// Compute the reference result of op(A)*x, storing it into want.
	impl.Ztrmv(uplo, trans, diag, n, a, lda, want, incX)
	// Compute op(A)*x, storing the result in-place into x.
	impl.Ztbmv(uplo, trans, diag, n, k, ab, ldab, x, incX)

	name := fmt.Sprintf("uplo=%v,trans=%v,diag=%v,n=%v,k=%v,ldab=%v,incX=%v", uplo, trans, diag, n, k, ldab, incX)
	if !zsame(ab, abCopy) {
		t.Errorf("%v: unexpected modification of ab", name)
	}
	if !zSameAtNonstrided(x, want, incX) {
		t.Errorf("%v: unexpected modification of x", name)
	}
	if !zEqualApproxAtStrided(x, want, incX, tol) {
		t.Errorf("%v: unexpected result\nwant %v\ngot  %v", name, want, x)
	}
}
