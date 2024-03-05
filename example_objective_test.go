// © 2019-present nextmv.io inc

package mip_test

import (
	"fmt"
	"testing"

	mip "github.com/nextmv-io/go-mip"
)

func ExampleObjective_sense() {
	model := mip.NewModel()

	model.Objective().SetMaximize()

	fmt.Println(model.Objective().IsMaximize())
	fmt.Println(model.Objective())

	model.Objective().SetMinimize()

	fmt.Println(model.Objective().IsMaximize())
	fmt.Println(model.Objective())
	// Output:
	// true
	// maximize
	// false
	// minimize
}

func ExampleObjective_terms() {
	model := mip.NewModel()

	v1 := model.NewBool()
	v2 := model.NewBool()

	fmt.Println(len(model.Objective().Terms()))

	t1 := model.Objective().NewTerm(2.0, v1)
	fmt.Println(t1)
	t2 := model.Objective().NewTerm(1.0, v1)
	fmt.Println(t2)
	t3 := model.Objective().NewTerm(3.0, v2)
	fmt.Println(t3)

	fmt.Println(t1.Var().Index())
	fmt.Println(t1.Coefficient())

	fmt.Println(t2.Var().Index())
	fmt.Println(t2.Coefficient())

	fmt.Println(t3.Var().Index())
	fmt.Println(t3.Coefficient())

	terms := model.Objective().Terms()
	fmt.Println(len(terms))
	for _, term := range terms {
		fmt.Println(term.Var(), term.Coefficient())
	}
	fmt.Println("isMaximize: ", model.Objective().IsMaximize())
	// Unordered output:
	// 0
	// 2 B0
	// 1 B0
	// 3 B1
	// 0
	// 2
	// 0
	// 1
	// 1
	// 3
	// 2
	// B0 3
	// B1 3
	// isMaximize:  false
}

func ExampleObjective_termsToString() {
	m := mip.NewModel()
	x0 := m.NewBool()
	x1 := m.NewBool()
	x2 := m.NewBool()
	x3 := m.NewBool()
	z := m.Objective()
	z.NewTerm(3, x2)
	z.NewTerm(2, x1)
	z.NewTerm(1, x0)
	fmt.Println(z)
	fmt.Println(z.Term(x0))
	z.NewTerm(1, x0)
	fmt.Println(z.Term(x0))
	fmt.Println(z.Term(x3))
	// Output:
	// minimize   1 B0 + 2 B1 + 3 B2
	// 1 B0 1
	// 2 B0 2
	// 0 B3 0
}

func benchmarkObjectiveNewTerms(nrTerms int, b *testing.B) {
	model := mip.NewModel()
	v := model.NewFloat(1.0, 2.0)

	for i := 0; i < b.N; i++ {
		for i := 0; i < nrTerms; i++ {
			model.Objective().NewTerm(1.0, v)
		}
	}
}

func BenchmarkObjectiveNewTerms1(b *testing.B) {
	benchmarkObjectiveNewTerms(1, b)
}

func BenchmarkObjectiveNewTerms2(b *testing.B) {
	benchmarkObjectiveNewTerms(2, b)
}

func BenchmarkObjectiveNewTerms4(b *testing.B) {
	benchmarkObjectiveNewTerms(4, b)
}

func BenchmarkObjectiveNewTerms8(b *testing.B) {
	benchmarkObjectiveNewTerms(8, b)
}

func BenchmarkObjectiveNewTerms16(b *testing.B) {
	benchmarkObjectiveNewTerms(16, b)
}

func BenchmarkObjectiveNewTerms32(b *testing.B) {
	benchmarkObjectiveNewTerms(32, b)
}
