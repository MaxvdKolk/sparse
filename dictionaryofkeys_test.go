package sparse

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestDOKConversion(t *testing.T) {
	var tests = []struct {
		m, n   int
		data   map[key]float64
		output []float64
	}{
		{
			m: 11, n: 11,
			data: map[key]float64{
				{0, 3}:   1,
				{1, 1}:   2,
				{2, 2}:   3,
				{5, 8}:   4,
				{10, 10}: 5,
				{1, 5}:   6,
				{3, 5}:   7,
			},
			output: []float64{
				0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0,
				0, 2, 0, 0, 0, 6, 0, 0, 0, 0, 0,
				0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 7, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5,
			},
		},
		{
			m: 5, n: 4,
			data: map[key]float64{
				{0, 3}: 1,
				{1, 1}: 2,
				{2, 2}: 3,
				{4, 2}: 4,
				{0, 0}: 5,
				{1, 3}: 6,
				{3, 3}: 7,
			},
			output: []float64{
				5, 0, 0, 1,
				0, 2, 0, 6,
				0, 0, 3, 0,
				0, 0, 0, 7,
				0, 0, 4, 0,
			},
		},
	}

	for ti, test := range tests {
		t.Logf("**** Test Run %d.\n", ti+1)
		expected := mat.NewDense(test.m, test.n, test.output)

		dok := NewDOK(test.m, test.n)
		for k, v := range test.data {
			dok.Set(k.i, k.j, v)
		}

		dok2 := dok.ToDOK()
		if !(mat.Equal(expected, dok2)) {
			t.Logf("Expected:\n%v \nbut found DOK matrix:\n%v\n", mat.Formatted(expected), mat.Formatted(dok2))
			t.Fail()
		}

		coo := dok.ToCOO()
		if !(mat.Equal(expected, coo)) {
			t.Logf("Expected:\n%v \nbut found COO matrix:\n%v\n", mat.Formatted(expected), mat.Formatted(coo))
			t.Fail()
		}

		csr := dok.ToCSR()
		if !(mat.Equal(expected, csr)) {
			t.Logf("Expected:\n%v \nbut found CSR matrix:\n%v\n", mat.Formatted(expected), mat.Formatted(csr))
			t.Fail()
		}

		csc := dok.ToCSC()
		if !(mat.Equal(expected, csc)) {
			t.Logf("Expected:\n%v \nbut found CSC matrix:\n%v\n", mat.Formatted(expected), mat.Formatted(csc))
			t.Fail()
		}

		dense := dok.ToDense()
		if !(mat.Equal(expected, dense)) {
			t.Logf("Expected:\n%v \nbut found Dense matrix:\n%v\n", mat.Formatted(expected), mat.Formatted(dense))
			t.Fail()
		}
	}
}

func TestDOKTranspose(t *testing.T) {
	var tests = []struct {
		r, c   int
		data   []float64
		er, ec int
		result []float64
	}{
		{
			r: 3, c: 4,
			data: []float64{
				1, 0, 0, 0,
				0, 2, 0, 0,
				0, 0, 3, 6,
			},
			er: 4, ec: 3,
			result: []float64{
				1, 0, 0,
				0, 2, 0,
				0, 0, 3,
				0, 0, 6,
			},
		},
	}

	for ti, test := range tests {
		t.Logf("**** Test Run %d.\n", ti+1)

		expected := mat.NewDense(test.er, test.ec, test.result)

		dok := CreateDOK(test.r, test.c, test.data)

		if !mat.Equal(expected, dok.T()) {
			t.Logf("Expected:\n %v\n but received:\n %v\n", mat.Formatted(expected), mat.Formatted(dok.T()))
			t.Fail()
		}
	}
}

func TestDOKDoNonZero(t *testing.T) {
	var tests = []struct {
		r, c int
		data []float64
	}{
		{
			r: 3, c: 4,
			data: []float64{
				1, 0, 6, 0,
				0, 2, 0, 0,
				0, 0, 3, 6,
			},
		},
		{
			r: 3, c: 4,
			data: []float64{
				1, 0, 7, 0,
				0, 0, 0, 0,
				6, 0, 3, 0,
			},
		},
	}

	for ti, test := range tests {
		t.Logf("**** Test Run %d.\n", ti+1)

		dok := CreateDOK(test.r, test.c, test.data).(*DOK)

		var nnz int
		dok.DoNonZero(func(i, j int, v float64) {
			if testv := test.data[i*test.c+j]; testv == 0 || testv != v {
				t.Logf("Expected %f at (%d, %d) but received %f\n", v, i, j, testv)
				t.Fail()
			}
			nnz++
		})

		if nnz != dok.NNZ() {
			t.Logf("Expected %d Non Zero elements but found %d", nnz, dok.NNZ())
			t.Fail()
		}
	}
}

type MatrixCreator func(m, n int, data []float64) mat.Matrix

func CreateDOK(m, n int, data []float64) mat.Matrix {
	dok := NewDOK(m, n)
	if data != nil {
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				if data[i*n+j] != 0 {
					dok.Set(i, j, data[i*n+j])
				}
			}
		}
	}

	return dok
}

func CreateCOO(m, n int, data []float64) mat.Matrix {
	return CreateDOK(m, n, data).(*DOK).ToCOO()
}

func CreateCSR(m, n int, data []float64) mat.Matrix {
	return CreateDOK(m, n, data).(*DOK).ToCSR()
}

func CreateCSC(m, n int, data []float64) mat.Matrix {
	return CreateDOK(m, n, data).(*DOK).ToCSC()
}

func CreateDIA(m, n int, data []float64) mat.Matrix {
	min := n
	if m <= n {
		min = m
	}

	c := make([]float64, min)
	for i := 0; i < min; i++ {
		c[i] = data[i*n+i]
	}
	return NewDIA(m, n, c)
}

func CreateDense(m, n int, data []float64) mat.Matrix {
	return mat.NewDense(m, n, data)
}

func TestFailDokInitialization(t *testing.T) {
	tcs := []struct {
		r, c int
	}{
		{-1, 1},
		{1, -1},
		{-1, -1},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Haven`t panic for wrong size: %v", tc)
				}
			}()
			_ = NewDOK(tc.r, tc.c)
		})
	}
}

func TestFailDokIndexes(t *testing.T) {
	d := NewDOK(1, 1)

	tcs := []struct {
		r, c int
	}{
		{-1, 1},
		{1, -1},
		{-1, -1},
		{0, 2},
		{2, 0},
		{100, 100},
	}

	for _, tc := range tcs {
		t.Run(fmt.Sprintf("At:%v", tc), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Haven`t panic for wrong index: %v", tc)
				}
			}()
			_ = d.At(tc.r, tc.c)
		})
		t.Run(fmt.Sprintf("Set:%v", tc), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Haven`t panic for wrong index: %v", tc)
				}
			}()
			d.Set(tc.r, tc.c, -1.0)
		})
	}
}
