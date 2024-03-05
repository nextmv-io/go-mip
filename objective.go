// © 2019-present nextmv.io inc

package mip

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// Objective specifies the objective of the model. An objective
// consists out of terms and a specification if it should be maximized or
// minimized.
//
// For example:
//
//	maximize 2.5 * x + 3.5 * y
//
// 2.5 * x and 3.5 * y are 2 terms in this example.
type Objective interface {
	// IsLinear returns true if the invoking objective is a linear function.
	IsLinear() bool
	// IsMaximize returns true if the invoking objective is a maximization
	// objective.
	IsMaximize() bool
	// IsQuadratic returns true if the invoking objective is a quadratic function.
	IsQuadratic() bool
	// NewTerm adds a term to the invoking objective, invoking this API
	// multiple times for the same variable will take the sum of coefficients
	// of earlier added terms for that variable.
	//
	// 		m := NewModel()
	// 		x := m.NewFloat(10.0, 100.0)
	//
	// 		m.Objective().SetMaximize()			// results in: maximize -
	// 		m.Objective().NewTerm(1.0, x)		// results in: maximize 1.0 * x
	// 		m.Objective().NewTerm(2.0, x)		// results in: maximize 3.0 * x
	NewTerm(coefficient float64, variable Var) Term
	// NewQuadraticTerm adds a new quadratic term to the invoking objective,
	// invoking this API multiple times for the same variables will take the sum
	// of coefficients of earlier added terms for that variable.
	//
	//      m := NewModel()
	//      x1 := m.NewFloat(10.0, 100.0)
	//      x2 := m.NewFloat(10.0, 100.0)
	//
	//      m.Objective().SetMaximize()
	//      // results in: maximize -
	//      m.Objective().NewQuadraticTerm(1.0, x1, x1)
	//      // results in: maximize 1.0 * x1^2
	//      m.Objective().NewQuadraticTerm(1.0, x1, x2)
	//      // results in: maximize 1.0 * x1^2 + x1x2
	//      m.Objective().NewQuadraticTerm(1.0, x2, x1)
	//      // results in: maximize 1.0 * x1^2 + 2.0 * x1x2
	NewQuadraticTerm(coefficient float64, variable1, variable2 Var) QuadraticTerm
	// SetMaximize sets the invoking objective to be a maximization objective.
	SetMaximize()
	// SetMinimize sets the invoking objective to be a minimization objective.
	SetMinimize()
	// Term returns a term for a given variable together with the sum of the
	// coefficients of all terms referencing that variable. The second return
	// argument defines how many terms have been defined on the objective for
	// the given variable.
	Term(variable Var) (Term, int)
	// Terms returns a copy slice of terms of the invoking objective,
	// each variable is reported once. If the same variable has been
	// added multiple times the sum of coefficients is reported for that
	// variable. The order of the terms is not specified and is not guaranteed
	// to be the same from one invocation to the next.
	Terms() Terms
	// QuadraticTerm returns a quadratic term for a given pair of variables
	// together with the sum of the coefficients of all quadratic terms
	// referencing that variable. The second return argument defines how many
	// quadratic terms have been defined on the objective the given pair of
	// variables.
	QuadraticTerm(variable1, variable2 Var) (QuadraticTerm, int)
	// QuadraticTerms returns a copy slice of quadratic terms of the invoking
	// objective, each variable pair is reported once. If the same pair has been
	// added multiple times the sum of coefficients is reported for that
	// variable. The order of the terms is not specified and is not guaranteed
	// to be the same from one invocation to the next.
	QuadraticTerms() QuadraticTerms
}

type objective struct {
	terms          Terms
	quadraticTerms QuadraticTerms
	maximize       bool
}

func (o *objective) SetMaximize() {
	o.maximize = true
}

func (o *objective) SetMinimize() {
	o.maximize = false
}

func (o *objective) NewTerm(
	coefficient float64,
	variable Var,
) Term {
	if math.IsNaN(coefficient) {
		panic("objective term coefficient is NaN")
	}

	term := &term{
		coefficient: coefficient,
		variable:    variable,
	}

	o.terms = append(o.terms, term)

	return term
}

func (o *objective) NewQuadraticTerm(
	coefficient float64,
	variable1 Var,
	variable2 Var,
) QuadraticTerm {
	if math.IsNaN(coefficient) {
		panic("constraint quadratic term coefficient is NaN")
	}

	term := newQuadraticTerm(coefficient, variable1, variable2)

	o.quadraticTerms = append(o.quadraticTerms, term)

	return term
}

func (o *objective) IsMaximize() bool {
	return o.maximize
}

func (o *objective) IsLinear() bool {
	return !o.IsQuadratic()
}

func (o *objective) IsQuadratic() bool {
	return o.quadraticTerms != nil && len(o.quadraticTerms) > 0
}

func (o *objective) Term(variable Var) (Term, int) {
	coefficient := 0.0
	definitions := 0

	for _, t := range o.terms {
		if t.Var().Index() == variable.Index() {
			definitions++
			coefficient += t.Coefficient()
		}
	}

	return &term{
		coefficient: coefficient,
		variable:    variable,
	}, definitions
}

func (o *objective) Terms() Terms {
	return makeLinearTermsUnique(o.terms)
}

func (o *objective) QuadraticTerm(
	variable1,
	variable2 Var,
) (QuadraticTerm, int) {
	coefficient := 0.0
	definitions := 0
	var var1, var2 Var
	if variable1.Index() <= variable2.Index() {
		var1 = variable1
		var2 = variable2
	} else {
		var1 = variable2
		var2 = variable1
	}

	for _, t := range o.quadraticTerms {
		if t.Var1().Index() == var1.Index() && t.Var2().Index() == var2.Index() {
			definitions++
			coefficient += t.Coefficient()
		}
	}

	return newQuadraticTerm(coefficient, var1, var2), definitions
}

func (o *objective) QuadraticTerms() QuadraticTerms {
	return makeQuadraticTermsUnique(o.quadraticTerms)
}

func (o *objective) String() string {
	var sb strings.Builder

	if o.IsMaximize() {
		sb.WriteString("maximize")
	} else {
		sb.WriteString("minimize")
	}

	operator := " "

	terms := o.Terms()

	sort.SliceStable(terms, func(i, j int) bool {
		return terms[i].Var().Index() < terms[j].Var().Index()
	})

	for _, t := range terms {
		fmt.Fprintf(&sb, " %v %v",
			operator,
			t)
		operator = "+"
	}

	qTerms := o.QuadraticTerms()

	sort.SliceStable(qTerms, func(i, j int) bool {
		return qTerms[i].Var1().Index() < qTerms[j].Var1().Index() ||
			(qTerms[i].Var1().Index() == qTerms[j].Var1().Index() &&
				qTerms[i].Var2().Index() < qTerms[j].Var2().Index())
	})

	for _, t := range qTerms {
		fmt.Fprintf(&sb, " %v %v",
			operator,
			t)
		operator = "+"
	}

	return sb.String()
}
