package mip

import (
	"fmt"
)

// Var represents the entities on which the solver has to make a decision
// without violating constraints and while optimizing the objective.
// Vars can be of a certain type, bool, float or int.
//
// Float vars can take a value of a float quantity
// Int vars are vars that must take an integer value
// (0, 1, 2, ...)
// Bool vars can take two values, zero or one.
type Var interface {
	// Index is a unique number assigned to the var. The index corresponds
	// to the location in the slice returned by Model.Variables().
	Index() int
	// IsBool returns true if the invoking variable is a bool variable,
	// otherwise it returns false.
	IsBool() bool
	// IsFloat returns true if the invoking variable is a float
	// variable otherwise false.
	IsFloat() bool
	// IsInt returns true if the invoking variable is an int variable
	// otherwise false.
	IsInt() bool
	// LowerBound returns the lowerBound of the invoking variable.
	//
	// Lower bounds of variables are limited by the lower bounds of the
	// underlying solver technology. The lower bound used will be the maximum
	// of the specification and the lower bound of the solver used.
	LowerBound() float64
	// Name returns assigned name. If no name has been set it will return
	// a unique auto-generated name.
	Name() string
	// SetName assigns name to invoking var
	SetName(name string)
	// UpperBound returns the upperBound of the invoking variable.
	//
	// Upper bounds of variables are limited by the upper bounds of the
	// underlying solver technology. The upper bound used will be the minimum
	// of the specification and the upper bound of the solver used.
	UpperBound() float64
}

// Vars is a slice of Var instances.
type Vars []Var

// Float a Var which can take any value in an interval.
type Float interface {
	Var
	ensureFloat() bool
}

// Int a Var which can take any integer value in an interval.
type Int interface {
	Var
	ensureInt() bool
}

// Bool a Var which can take two values, zero or one. A bool
// variable is also an int variable which can have two values zero and
// one.
type Bool interface {
	Int
	ensureBool() bool
}

type variable struct {
	model *model
	index int
}

type floatVariable struct {
	Float
	variable
	lowerBound float64
	upperBound float64
}

func (f *floatVariable) Index() int {
	return f.index
}

func (f *floatVariable) IsBool() bool {
	return false
}

func (f *floatVariable) IsFloat() bool {
	return true
}

func (f *floatVariable) IsInt() bool {
	return false
}

func (f *floatVariable) LowerBound() float64 {
	return f.lowerBound
}

func (f *floatVariable) Name() string {
	return f.model.getVarName(f)
}

func (f *floatVariable) SetName(name string) {
	f.model.setVarName(f, name)
}

func (f *floatVariable) UpperBound() float64 {
	return f.upperBound
}

func (f *floatVariable) String() string {
	name := f.Name()
	if name == "" {
		name = fmt.Sprintf("F%v", f.Index())
	}
	return name
}

type intVariable struct {
	Int
	variable
	lowerBound int64
	upperBound int64
}

func (i *intVariable) Index() int {
	return i.index
}

func (i *intVariable) IsBool() bool {
	return false
}

func (i *intVariable) IsFloat() bool {
	return false
}

func (i *intVariable) IsInt() bool {
	return true
}

func (i *intVariable) LowerBound() float64 {
	return float64(i.lowerBound)
}

func (i *intVariable) Name() string {
	return i.model.getVarName(i)
}

func (i *intVariable) SetName(name string) {
	i.model.setVarName(i, name)
}

func (i *intVariable) UpperBound() float64 {
	return float64(i.upperBound)
}

func (i *intVariable) String() string {
	name := i.Name()
	if name == "" {
		name = fmt.Sprintf("I%v", i.Index())
	}
	return name
}

type boolVariable struct {
	Bool
	variable
}

func (b *boolVariable) Index() int {
	return b.index
}

func (b *boolVariable) IsBool() bool {
	return true
}

func (b *boolVariable) IsFloat() bool {
	return false
}

func (b *boolVariable) IsInt() bool {
	return true
}

func (b *boolVariable) LowerBound() float64 {
	return 0.0
}

func (b *boolVariable) Name() string {
	return b.model.getVarName(b)
}

func (b *boolVariable) SetName(name string) {
	b.model.setVarName(b, name)
}

func (b *boolVariable) UpperBound() float64 {
	return 1.0
}

func (b *boolVariable) String() string {
	name := b.Name()
	if name == "" {
		name = fmt.Sprintf("B%v", b.Index())
	}
	return name
}
