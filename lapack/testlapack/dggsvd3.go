// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/floats/scalar"
	"gonum.org/v1/gonum/lapack"
)

type Dggsvd3er interface {
	Dggsvd3(jobU, jobV, jobQ lapack.GSVDJob, m, n, p int, a []float64, lda int, b []float64, ldb int, alpha, beta, u []float64, ldu int, v []float64, ldv int, q []float64, ldq int, work []float64, lwork int, iwork []int) (k, l int, ok bool)
}

func Dggsvd3Test(t *testing.T, impl Dggsvd3er) {
	const tol = 1e-13

	rnd := rand.New(rand.NewPCG(1, 1))
	for cas, test := range []struct {
		m, p, n, lda, ldb, ldu, ldv, ldq int

		ok bool
	}{
		{m: 3, p: 3, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 5, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 5, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 5, p: 5, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 5, p: 5, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 5, p: 5, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 10, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 10, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 10, p: 10, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 10, p: 10, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0, ok: true},
		{m: 5, p: 5, n: 5, lda: 10, ldb: 10, ldu: 10, ldv: 10, ldq: 10, ok: true},
		{m: 5, p: 5, n: 5, lda: 10, ldb: 10, ldu: 10, ldv: 10, ldq: 10, ok: true},
		{m: 5, p: 5, n: 10, lda: 20, ldb: 20, ldu: 10, ldv: 10, ldq: 20, ok: true},
		{m: 5, p: 5, n: 10, lda: 20, ldb: 20, ldu: 10, ldv: 10, ldq: 20, ok: true},
		{m: 5, p: 5, n: 10, lda: 20, ldb: 20, ldu: 10, ldv: 10, ldq: 20, ok: true},
		{m: 10, p: 5, n: 5, lda: 10, ldb: 10, ldu: 20, ldv: 10, ldq: 10, ok: true},
		{m: 10, p: 5, n: 5, lda: 10, ldb: 10, ldu: 20, ldv: 10, ldq: 10, ok: true},
		{m: 10, p: 10, n: 10, lda: 20, ldb: 20, ldu: 20, ldv: 20, ldq: 20, ok: true},
		{m: 10, p: 10, n: 10, lda: 20, ldb: 20, ldu: 20, ldv: 20, ldq: 20, ok: true},
	} {
		m := test.m
		p := test.p
		n := test.n
		lda := test.lda
		if lda == 0 {
			lda = n
		}
		ldb := test.ldb
		if ldb == 0 {
			ldb = n
		}
		ldu := test.ldu
		if ldu == 0 {
			ldu = m
		}
		ldv := test.ldv
		if ldv == 0 {
			ldv = p
		}
		ldq := test.ldq
		if ldq == 0 {
			ldq = n
		}

		a := randomGeneral(m, n, lda, rnd)
		aCopy := cloneGeneral(a)
		b := randomGeneral(p, n, ldb, rnd)
		bCopy := cloneGeneral(b)

		alpha := make([]float64, n)
		beta := make([]float64, n)

		u := nanGeneral(m, m, ldu)
		v := nanGeneral(p, p, ldv)
		q := nanGeneral(n, n, ldq)

		iwork := make([]int, n)

		work := []float64{0}
		impl.Dggsvd3(lapack.GSVDU, lapack.GSVDV, lapack.GSVDQ,
			m, n, p,
			a.Data, a.Stride,
			b.Data, b.Stride,
			alpha, beta,
			u.Data, u.Stride,
			v.Data, v.Stride,
			q.Data, q.Stride,
			work, -1, iwork)

		lwork := int(work[0])
		work = make([]float64, lwork)

		k, l, ok := impl.Dggsvd3(lapack.GSVDU, lapack.GSVDV, lapack.GSVDQ,
			m, n, p,
			a.Data, a.Stride,
			b.Data, b.Stride,
			alpha, beta,
			u.Data, u.Stride,
			v.Data, v.Stride,
			q.Data, q.Stride,
			work, lwork, iwork)

		if !ok {
			if test.ok {
				t.Errorf("test %d unexpectedly did not converge", cas)
			}
			continue
		}

		// Check orthogonality of U, V and Q.
		if resid := residualOrthogonal(u, false); resid > tol {
			t.Errorf("Case %v: U is not orthogonal; resid=%v, want<=%v", cas, resid, tol)
		}
		if resid := residualOrthogonal(v, false); resid > tol {
			t.Errorf("Case %v: V is not orthogonal; resid=%v, want<=%v", cas, resid, tol)
		}
		if resid := residualOrthogonal(q, false); resid > tol {
			t.Errorf("Case %v: Q is not orthogonal; resid=%v, want<=%v", cas, resid, tol)
		}

		// Check C^2 + S^2 = I.
		var elements []float64
		if m-k-l >= 0 {
			elements = alpha[k : k+l]
		} else {
			elements = alpha[k:m]
		}
		for i := range elements {
			i += k
			d := alpha[i]*alpha[i] + beta[i]*beta[i]
			if !scalar.EqualWithinAbsOrRel(d, 1, tol, tol) {
				t.Errorf("test %d: alpha_%d^2 + beta_%d^2 != 1: got: %v", cas, i, i, d)
			}
		}

		zeroR, d1, d2 := constructGSVDresults(n, p, m, k, l, a, b, alpha, beta)

		// Check Uᵀ*A*Q = D1*[ 0 R ].
		uTmp := nanGeneral(m, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, u, aCopy, 0, uTmp)
		uAns := nanGeneral(m, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, uTmp, q, 0, uAns)

		d10r := nanGeneral(m, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, d1, zeroR, 0, d10r)

		if !equalApproxGeneral(uAns, d10r, tol) {
			t.Errorf("test %d: Uᵀ*A*Q != D1*[ 0 R ]\nUᵀ*A*Q:\n%+v\nD1*[ 0 R ]:\n%+v",
				cas, uAns, d10r)
		}

		// Check Vᵀ*B*Q = D2*[ 0 R ].
		vTmp := nanGeneral(p, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, v, bCopy, 0, vTmp)
		vAns := nanGeneral(p, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, vTmp, q, 0, vAns)

		d20r := nanGeneral(p, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, d2, zeroR, 0, d20r)

		if !equalApproxGeneral(vAns, d20r, tol) {
			t.Errorf("test %d: Vᵀ*B*Q != D2*[ 0 R ]\nVᵀ*B*Q:\n%+v\nD2*[ 0 R ]:\n%+v",
				cas, vAns, d20r)
		}
	}
}
