package mip

import (
	"fmt"
)

// Term is the building block of a constraint and an objective. A term consist
// of a coefficient and a var and should be interpreted as the product
// of coefficient and the var in the context of the constraint or
// objective.
type Term interface {
	// Coefficient returns the coefficient value of the invoking term.
	Coefficient() float64
	// Var returns the variable of the term.
	Var() Var
}

// Terms is a slice of Term instances.
type Terms []Term

// QuadraticTerm consists of a coefficient and two vars. It should be
// interpreted as the product of a coefficient and the two vars.
type QuadraticTerm interface {
	// Coefficient returns the coefficient value of the invoking term.
	Coefficient() float64
	// Var1 returns the first variable.
	Var1() Var
	// Var2 returns the second variable.
	Var2() Var
}

// QuadraticTerms is a slice of QuadraticTerm instances.
type QuadraticTerms []QuadraticTerm

func makeLinearTermsUnique(terms Terms) Terms {
	uniqueTerms := make(map[Var]Term)

	for _, t := range terms {
		if t.Coefficient() == 0 {
			continue
		}
		if _, ok := uniqueTerms[t.Var()]; !ok {
			uniqueTerms[t.Var()] = t
		} else {
			coefficient := uniqueTerms[t.Var()].Coefficient() +
				t.Coefficient()

			if coefficient == 0.0 {
				delete(uniqueTerms, t.Var())
			} else {
				uniqueTerms[t.Var()] = &term{
					coefficient: coefficient,
					variable:    t.Var(),
				}
			}
		}
	}

	uniqueTermsAsSlice := make(Terms, 0, len(uniqueTerms))

	for _, t := range uniqueTerms {
		uniqueTermsAsSlice = append(uniqueTermsAsSlice, t)
	}

	return uniqueTermsAsSlice
}

type term struct {
	variable    Var
	coefficient float64
}

func (t *term) Coefficient() float64 {
	return t.coefficient
}

func (t *term) Var() Var {
	return t.variable
}

func (t *term) String() string {
	return fmt.Sprintf("%v %v",
		t.coefficient,
		t.variable)
}

// quadraticTerm maintains the invariant that variable1.Index() <=
// variable2.Index(). use newQuadraticTerm to ensure you don't forget that.
type quadraticTerm struct {
	variable1   Var
	variable2   Var
	coefficient float64
}

func (t *quadraticTerm) Coefficient() float64 {
	return t.coefficient
}

func (t *quadraticTerm) Var1() Var {
	return t.variable1
}

func (t *quadraticTerm) Var2() Var {
	return t.variable2
}

func (t *quadraticTerm) String() string {
	if t.variable1.Index() == t.variable2.Index() {
		return fmt.Sprintf("%v %v^2",
			t.coefficient,
			t.variable1)
	}
	return fmt.Sprintf("%v %v*%v",
		t.coefficient,
		t.variable1,
		t.variable2)
}

func makeQuadraticTermsUnique(
	quadraticTerms QuadraticTerms,
) QuadraticTerms {
	terms := make(map[Var]map[Var]QuadraticTerm)
	termCounter := 0
	for _, v := range quadraticTerms {
		v1 := v.Var1()
		v2 := v.Var2()
		if v.Coefficient() == 0.0 {
			continue
		}
		mqt, ok := terms[v1]
		if !ok {
			terms[v1] = make(map[Var]QuadraticTerm)
			terms[v1][v2] = v
			termCounter++
			continue
		}
		qt, ok := mqt[v2]
		if !ok {
			mqt[v2] = v
			termCounter++
			continue
		}
		newCoef := v.Coefficient() + qt.Coefficient()
		if newCoef == 0.0 {
			delete(mqt, v2)
			termCounter--
			continue
		}
		mqt[v2] = &quadraticTerm{
			coefficient: newCoef,
			variable1:   v1,
			variable2:   v2,
		}
	}

	rTerms := make(QuadraticTerms, 0, termCounter)
	for k := range terms {
		for _, v := range terms[k] {
			rTerms = append(rTerms, v)
		}
	}
	return rTerms
}

func newQuadraticTerm(
	coefficient float64,
	variable1 Var,
	variable2 Var,
) QuadraticTerm {
	if variable1.Index() <= variable2.Index() {
		return &quadraticTerm{
			coefficient: coefficient,
			variable1:   variable1,
			variable2:   variable2,
		}
	}
	return &quadraticTerm{
		coefficient: coefficient,
		variable1:   variable2,
		variable2:   variable1,
	}
}
