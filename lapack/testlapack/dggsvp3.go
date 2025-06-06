// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testlapack

import (
	"math/rand/v2"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

type Dggsvp3er interface {
	Dlanger
	Dggsvp3(jobU, jobV, jobQ lapack.GSVDJob, m, p, n int, a []float64, lda int, b []float64, ldb int, tola, tolb float64, u []float64, ldu int, v []float64, ldv int, q []float64, ldq int, iwork []int, tau, work []float64, lwork int) (k, l int)
}

func Dggsvp3Test(t *testing.T, impl Dggsvp3er) {
	const tol = 1e-14

	rnd := rand.New(rand.NewPCG(1, 1))
	for cas, test := range []struct {
		m, p, n, lda, ldb, ldu, ldv, ldq int
	}{
		{m: 3, p: 3, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 5, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 5, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 5, p: 5, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 5, p: 5, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 5, p: 5, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 10, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 10, p: 5, n: 5, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 10, p: 10, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 10, p: 10, n: 10, lda: 0, ldb: 0, ldu: 0, ldv: 0, ldq: 0},
		{m: 5, p: 5, n: 5, lda: 10, ldb: 10, ldu: 10, ldv: 10, ldq: 10},
		{m: 5, p: 5, n: 5, lda: 10, ldb: 10, ldu: 10, ldv: 10, ldq: 10},
		{m: 5, p: 5, n: 10, lda: 20, ldb: 20, ldu: 10, ldv: 10, ldq: 20},
		{m: 5, p: 5, n: 10, lda: 20, ldb: 20, ldu: 10, ldv: 10, ldq: 20},
		{m: 5, p: 5, n: 10, lda: 20, ldb: 20, ldu: 10, ldv: 10, ldq: 20},
		{m: 10, p: 5, n: 5, lda: 10, ldb: 10, ldu: 20, ldv: 10, ldq: 10},
		{m: 10, p: 5, n: 5, lda: 10, ldb: 10, ldu: 20, ldv: 10, ldq: 10},
		{m: 10, p: 10, n: 10, lda: 20, ldb: 20, ldu: 20, ldv: 20, ldq: 20},
		{m: 10, p: 10, n: 10, lda: 20, ldb: 20, ldu: 20, ldv: 20, ldq: 20},
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

		tola := float64(max(m, n)) * impl.Dlange(lapack.Frobenius, m, n, a.Data, a.Stride, nil) * dlamchE
		tolb := float64(max(p, n)) * impl.Dlange(lapack.Frobenius, p, n, b.Data, b.Stride, nil) * dlamchE

		u := nanGeneral(m, m, ldu)
		v := nanGeneral(p, p, ldv)
		q := nanGeneral(n, n, ldq)

		iwork := make([]int, n)
		tau := make([]float64, n)

		work := []float64{0}
		impl.Dggsvp3(lapack.GSVDU, lapack.GSVDV, lapack.GSVDQ,
			m, p, n,
			a.Data, a.Stride,
			b.Data, b.Stride,
			tola, tolb,
			u.Data, u.Stride,
			v.Data, v.Stride,
			q.Data, q.Stride,
			iwork, tau,
			work, -1)

		lwork := int(work[0])
		work = make([]float64, lwork)

		k, l := impl.Dggsvp3(lapack.GSVDU, lapack.GSVDV, lapack.GSVDQ,
			m, p, n,
			a.Data, a.Stride,
			b.Data, b.Stride,
			tola, tolb,
			u.Data, u.Stride,
			v.Data, v.Stride,
			q.Data, q.Stride,
			iwork, tau,
			work, lwork)

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

		zeroA, zeroB := constructGSVPresults(n, p, m, k, l, a, b)

		// Check Uᵀ*A*Q = [ 0 RA ].
		uTmp := nanGeneral(m, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, u, aCopy, 0, uTmp)
		uAns := nanGeneral(m, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, uTmp, q, 0, uAns)

		if !equalApproxGeneral(uAns, zeroA, tol) {
			t.Errorf("test %d: Uᵀ*A*Q != [ 0 RA ]\nUᵀ*A*Q:\n%+v\n[ 0 RA ]:\n%+v",
				cas, uAns, zeroA)
		}

		// Check Vᵀ*B*Q = [ 0 RB ].
		vTmp := nanGeneral(p, n, n)
		blas64.Gemm(blas.Trans, blas.NoTrans, 1, v, bCopy, 0, vTmp)
		vAns := nanGeneral(p, n, n)
		blas64.Gemm(blas.NoTrans, blas.NoTrans, 1, vTmp, q, 0, vAns)

		if !equalApproxGeneral(vAns, zeroB, tol) {
			t.Errorf("test %d: Vᵀ*B*Q != [ 0 RB ]\nVᵀ*B*Q:\n%+v\n[ 0 RB ]:\n%+v",
				cas, vAns, zeroB)
		}
	}
}
