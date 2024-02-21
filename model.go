// Package mip holds the implementation of the sdk/mip package.
package mip

import (
	"fmt"
	"math"
	"strings"
)

// Model manages the variables, constraints and objective.
type Model interface {
	// Constraints returns a copy slice of all constraints.
	Constraints() Constraints
	// Copy returns a copy of the model.
	Copy() Model
	// NewBool adds a bool variable to the invoking model,
	// returns the newly constructed variable.
	NewBool() Bool
	// NewFloat adds a float var with bounds [lowerBound,
	// upperBound] to the invoking model, returns the newly constructed
	// var.
	NewFloat(
		lowerBound float64,
		upperBound float64,
	) Float
	// NewInt adds an integer var with bounds [loweBound,
	// upperBound] to the invoking model, returns the newly constructed
	// var.
	NewInt(
		lowerBound int64,
		upperBound int64,
	) Int
	// NewConstraint adds a constraint with sense and right-hand-side value rhs
	// to the invoking model. All terms for existing and future variables
	// are initially zero. Returns the newly constructed constraint.
	// A constraint where all terms remain zero is ignored by the solver.
	NewConstraint(sense Sense, rhs float64) Constraint
	// Objective returns the objective of the model.
	Objective() Objective
	// Vars returns a copy slice of all vars.
	Vars() Vars
}

// NewModel SDK implementation.
func NewModel() Model {
	return &model{
		constraints:     make(Constraints, 0),
		constraintNames: make(map[Constraint]string),
		objective: &objective{
			maximize: false,
			terms:    make(Terms, 0),
		},
		vars:     make(Vars, 0),
		varNames: make(map[Var]string),
	}
}

type model struct {
	objective       Objective
	constraintNames map[Constraint]string
	varNames        map[Var]string
	constraints     Constraints
	vars            Vars
}

func (m *model) setConstraintName(constraint Constraint, name string) {
	m.constraintNames[constraint] = name
}

func (m *model) getConstraintName(constraint Constraint) string {
	if name, ok := m.constraintNames[constraint]; ok {
		return name
	}
	return ""
}

func (m *model) setVarName(variable Var, name string) {
	m.varNames[variable] = name
}

func (m *model) getVarName(variable Var) string {
	if name, ok := m.varNames[variable]; ok {
		return name
	}
	return ""
}

func (m *model) Constraints() Constraints {
	constraints := make(Constraints, len(m.constraints))

	copy(constraints, m.constraints)

	return constraints
}

func (m *model) Copy() Model {
	copyModel := NewModel()

	for _, v := range m.Vars() {
		switch {
		case v.IsFloat():
			{
				copyVar := copyModel.NewFloat(
					v.LowerBound(),
					v.UpperBound(),
				)
				copyVar.SetName(v.Name())
			}
		case v.IsBool():
			{
				copyVar := copyModel.NewBool()
				copyVar.SetName(v.Name())
			}
		case v.IsInt():
			copyVar := copyModel.NewInt(
				int64(v.LowerBound()),
				int64(v.UpperBound()),
			)
			copyVar.SetName(v.Name())
		}
	}

	if m.Objective().IsMaximize() {
		copyModel.Objective().SetMaximize()
	} else {
		copyModel.Objective().SetMinimize()
	}

	vars := copyModel.Vars()

	for _, t := range m.Objective().Terms() {
		copyModel.Objective().NewTerm(
			t.Coefficient(),
			vars[t.Var().Index()],
		)
	}
	for _, c := range m.Constraints() {
		copyConstraint := copyModel.NewConstraint(
			c.Sense(),
			c.RightHandSide(),
		)
		for _, t := range c.Terms() {
			copyConstraint.NewTerm(
				t.Coefficient(),
				vars[t.Var().Index()],
			)
		}
		copyConstraint.SetName(c.Name())
	}

	return copyModel
}

func (m *model) Objective() Objective {
	return m.objective
}

func (m *model) Vars() Vars {
	variables := make(Vars, len(m.vars))

	copy(variables, m.vars)

	return variables
}

func (m *model) NewBool() Bool {
	b := &boolVariable{
		variable: variable{
			index: len(m.vars),
			model: m,
		},
	}

	m.vars = append(m.vars, b)

	return b
}

func (m *model) NewFloat(
	lowerBound float64,
	upperBound float64,
) Float {
	if math.IsNaN(lowerBound) {
		panic("lower bound is NaN")
	}
	if math.IsNaN(upperBound) {
		panic("upper bound is NaN")
	}

	f := &floatVariable{
		variable: variable{
			index: len(m.vars),
			model: m,
		},
		lowerBound: lowerBound,
		upperBound: upperBound,
	}

	m.vars = append(m.vars, f)

	return f
}

func (m *model) NewInt(
	lowerBound int64,
	upperBound int64,
) Int {
	i := &intVariable{
		variable: variable{
			index: len(m.vars),
			model: m,
		},
		lowerBound: lowerBound,
		upperBound: upperBound,
	}

	m.vars = append(m.vars, i)

	return i
}

func (m *model) NewConstraint(
	sense Sense,
	rightHandSide float64,
) Constraint {
	if math.IsNaN(rightHandSide) {
		panic("right hand side is NaN")
	}
	constraint := &constraint{
		model:         m,
		rightHandSide: rightHandSide,
		sense:         sense,
		terms:         make([]Term, 0),
	}

	m.constraints = append(m.constraints, constraint)

	return constraint
}

func (m *model) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%v\n", m.objective)
	for i, c := range m.constraints {
		fmt.Fprintf(&sb, "%7d: %v\n", i, c)
	}
	for i, v := range m.vars {
		fmt.Fprintf(&sb, "%7d: %v [%v, %v]\n",
			i,
			v,
			v.LowerBound(),
			v.UpperBound(),
		)
	}
	return sb.String()
}
