// © 2019-present nextmv.io inc

/*
Package mip provides a general interface for solving mixed integer linear
optimization problems using a variety of back-end solvers. The base interface is
the Model which is a collection of variables, constraints and an objective. The
interface Solver is constructed by mip.NewSolver. The solver can be invoked
using Solver.Solve and returns a Solution.

A new Model is created:

	d := mip.NewModel()

Var instances are created and added to the model:

	x := d.NewFloat(0.0, 100.0)
	y := d.NewInt(0, 100)

Constraint instances are created and added to the model:

	c1 := d.NewConstraint(mip.GreaterThanOrEqual, 1.0)
	c1.NewTerm(-2.0, x)
	c1.NewTerm(2.0, y)

	c2 := d.NewConstraint(mip.LessThanOrEqual, 13.0)
	c2.NewTerm(-8.0, x)
	c2.NewTerm(10.0, y)

The Objective is specified:

	d.Objective().SetMaximize()
	d.Objective().NewTerm(1.0, x)
	d.Objective().NewTerm(1.0, y)

A Solver is created and invoked to produce a Solution:

	solver, _ := mip.NewSolver("backend_solver_identifier", mipModel)
	solution, _ := solver.Solve(mip.DefaultSolverOptions())
*/
package mip
