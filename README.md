# go-mip

**go-mip** provides a Go interface for modeling and solving Mixed-Integer
Programming (MIP) problems.

## License

Please note that go-mip is provided as _source-available_ software (not
_open-source_). For further information, please refer to the
[LICENSE](./LICENSE.md) file.

## Usage

This library contains the types and interfaces to model and solve MIP problems.
The implementation is provided by other (solver-specific) packages, e.g.:
[go-highs](https://github.com/nextmv-io/go-highs).

These are the supported solvers and solver interfaces for MIP:

| Provider | Description |
| -------- | ----------- |
| [FICO Xpress][fico-xpress-section] | Commercial solver. |
| [HiGHS][highs-section] | Open-source solver.  |
| [OR-Tools][ortools-section] | Open-source solver interface for various solvers. Included open-source solvers: `CLP`, `GLOP`, `PDLP`, `SCIP`. |
| [Pyomo][pyomo-section] | Open-source interface for various solvers. Included open-source solvers in Cloud: `CBC`, `GLPK`. |

More information on supported solvers and solver interfaces is available in the
[Integrations section of Nextmv Docs][nextmv-docs-integrations].

## Options

All MIP applications are configurable through options.

* Runner options change how the runner outputs solutions and reads/writes data.
* Solve options change the behavior of the solver.

Options can be configured when running locally with the Nextmv CLI or when
running remotely on Nextmv Platform. Running remotely is possible with a custom
app or an app subscribed to a Marketplace app. Learn more about [Nextmv
apps][nextmv-docs-apps].

Note that if both an environment variable and its corresponding CLI flag are
defined, the CLI flag will overwrite the environment variable.

These are the default options that are available with MIP:

* `-runner.input.path` string
  * The input file path (env `RUNNER_INPUT_PATH`).
* `-runner.output.path` string
  * The output file path (env `RUNNER_OUTPUT_PATH`).
* `-runner.output.solutions` string
  * {all, last} (env `RUNNER_OUTPUT_SOLUTIONS`) (default "last").
* `-runner.profile.cpu` string
  * The CPU profile file path (env `RUNNER_PROFILE_CPU`).
* `-runner.profile.memory` string
  * The memory profile file path (env `RUNNER_PROFILE_MEMORY`).
* `-solve.control.bool` string
  * List of solver-specific control options (configurations) with bool values
    (env `SOLVE_CONTROL_BOOL`). Example: "name1=value1,name2=value2", where
    value1 and value2 are bool values.
* `-solve.control.float` string
  * List of solver-specific control options (configurations) with float values
    (env `SOLVE_CONTROL_FLOAT`). Example: "name1=value1,name2=value2", where
    value1 and value2 are float values.
* `-solve.control.int` string
  * List of solver-specific control options (configurations) with int values
    (env `SOLVE_CONTROL_INT`). Example: "name1=value1,name2=value2", where
    value1 and value2 are int values.
* `-solve.control.string` string
  * List of solver-specific control options (configurations) with string values
    (env `SOLVE_CONTROL_STRING`). Example: "name1=value1,name2=value2", where
    value1 and value2 are string values.
* `-solve.duration` duration
  * Maximum duration of the solver (env `SOLVE_DURATION`) (default 30s).
* `-solve.mip.gap.absolute` float
  * Absolute gap stopping value. If the problem is an integer problem the solver
    will stop if the gap between the relaxed problem and the best found integer
    problem is less than this value (env `SOLVE_MIP_GAP_ABSOLUTE`) (default
    1e-06).
* `-solve.mip.gap.relative` float
  * Relative gap stopping value. If the problem is an integer problem the solver
    will stop if the relative gap between the relaxed problem and the best found
    integer problem is less than this value (env `SOLVE_MIP_GAP_RELATIVE`)
    (default 0.0001).
* `-solve.verbosity` value
  * {off, low, medium, high} Verbosity of the solver in the console (env
    `SOLVE_VERBOSITY`).

Valid time units for `-solve.duration` are as follows, according to
[`time.ParseDuration` from Go's standard library][go-parse-duration]:

* `ns` (nanoseconds)
* `us/µs` (microseconds)
* `ms` (milliseconds)
* `s` (seconds)
* `m` (minutes)
* `h` (hours)

The options are marshalled to the output when running an app, under the
`options` key. Here is an example of how the options are displayed (solve
duration is displayed in `ns`):

```json
"options": {
  "solve": {
    "control": {
      "bool": [],
      "float": [],
      "int": [],
      "string": []
    },
    "duration": 10000000000,
    "mip": {
      "gap": {
        "absolute": 0.000001,
        "relative": 0.0001
      }
    },
    "verbosity": "off"
  }
}
```

### Running locally

CLI flags are added after the double-dash character (`--`):

```bash
go run . \
  -runner.input.path input.json \
  -runner.output.path output.json \
  -solve.duration 10s
```

To set an environment variable, convert its corresponding CLI flag to uppercase,
replacing each period (.) with an underscore (_) and removing the leading dash
(-). For example: `-solve.duration` is equivalent to `SOLVE_DURATION`.

```bash
RUNNER_INPUT_PATH=input.json \
RUNNER_OUTPUT_PATH=output.json \
SOLVE_DURATION=10s \
  go run .
```

(Note that there is no longer a need for the double-dash character.)

If both an environment variable and its corresponding CLI flag are defined, the
CLI flag will overwrite the environment variable.

### Running remotely
  
Use the `"options"` key in the `JSON` payload when executing a run on Nextmv
cloud. The leading dash (`-`) is removed from the option name, when compared to
running locally. For example, `-solve.duration` on a local run is equivalent to
`solve.duration` when running remotely.

:warning: Note that all values should be set as a `string`, regardless of the type.

Here is an example of how to use options when running remotely:

{% tabs %}
{% tab label="JSON" language="json" %}

```json
{
  "input": {},
  "options": {
    "solve.duration": "3s",
    "solve.mip.gap.relative": "0.15",
    "solve.verbosity": "high"
  }
}
```

{% /tab %}
{% /tabs %}

### Output

Most options are parsed to the output, under the `options` key:

```json
"options": {
  "solve": {
    "control": {
      "bool": [],
      "float": [],
      "int": [],
      "string": []
    },
    "duration": 10000000000,
    "mip": {
      "gap": {
        "absolute": 0.000001,
        "relative": 0.0001
      }
    },
    "verbosity": "off"
  }
},
```

You may also modify the options in your app code, which would override the ones
passed from the runner.

### Solver control parameters

The Nextmv SDK has options for controlling the solver that work regardless of
the provider being used. Some of these options include (using CLI-style flags):

* `-solve.duration`: maximum duration of the solve.
* `-solve.verbosity`: verbosity of the solver.
* `-solve.mip.gap.relative`: relative gap for the MIP solver.
* `-solve.mip.gap.absolute`: absolute gap for the MIP solver.

Each provider has its own set of parameters for controlling their solver.
Depending on the option’s data type, you can use one of the following flags to
set these parameters (the list showcases CLI-style flags):

* `-solve.control.bool`: parameters with `bool` values.
* `-solve.control.float`: parameters with `float` values.
* `-solve.control.int`: parameters with `int` values.
* `-solve.control.string` parameters with `string` values.

Each of these options receives a `name=value` pair. To set multiple options of
the same type, use a comma-separated list: `name1=value1,name2=value2`. Note
that not all providers support all data types.

Here is an example where Nextmv SDK options are set alongside [HiGHS-specific
options][highs-options] (using the Nextmv CLI). Notice that some options set
multiple solver parameters.

```bash
nextmv sdk run . -- \
  -runner.input.path input.json \
  -runner.output.path output.json \
  -solve.duration 10s \
  -solve.mip.gap.absolute 80 \
  -solve.mip.gap.relative 0.4 \
  -solve.verbosity high \
  -solve.control.float "mip_heuristic_effort=0.7" \
  -solve.control.int "mip_max_nodes=200,threads=1" \
  -solve.control.string "presolve=off"
```

Given that verbosity is set to `high`, the console will show the solver’s
log (exposing that the options were set):

```bash
Running HiGHS 1.3.1 [date: 2023-01-01, git hash: abcd1234]
Copyright (c) 2022 ERGO-Code under MIT licence terms
Running with 1 thread(s)

Presolve is switched off
Objective function is integral with scale 1

Solving MIP model with:
   1 rows
   11 cols (11 binary, 0 integer, 0 implied int., 0 continuous)
   11 nonzeros

        Nodes      |    B&B Tree     |            Objective Bounds              |  Dynamic Constraints |       Work      
     Proc. InQueue |  Leaves   Expl. | BestBound       BestSol              Gap |   Cuts   InLp Confl. | LpIters     Time

         0       0         0   0.00%   543             -inf                 inf        0      0      0         0     0.0s
 S       0       0         0   0.00%   543             444               22.30%        0      0      0         0     0.0s

Solving report
  Status            Optimal
  Primal bound      444
  Dual bound        451
  Gap               1.58% (tolerance: 40%)
  Solution status   feasible
                    444 (objective)
                    0 (bound viol.)
                    0 (int. viol.)
                    0 (row viol.)
  Timing            0.00 (total)
                    0.00 (presolve)
                    0.00 (postsolve)
  Nodes             1
  LP iterations     1 (total)
                    0 (strong br.)
                    0 (separation)
                    0 (heuristics)
```

The output also shows the options that were used:

```json
"options": {
  "solve": {
    "control": {
      "bool": [],
      "float": [
        {
          "name": "mip_heuristic_effort",
          "value": 0.7
        }
      ],
      "int": [
        {
          "name": "mip_max_nodes",
          "value": 200
        },
        {
          "name": "threads",
          "value": 1
        }
      ],
      "string": [
        {
          "name": "presolve",
          "value": "off"
        }
      ]
    },
    "duration": 10000000000,
    "mip": {
      "gap": {
        "absolute": 80,
        "relative": 0.4
      }
    },
    "verbosity": "off"
  }
},
```

[nextmv-docs-integrations]: https://www.nextmv.io/docs/integrations/overview
[ortools-section]: /mixed-integer-programming/supported-solvers#or-tools
[highs-section]: /mixed-integer-programming/supported-solvers#highs
[fico-xpress-section]: /mixed-integer-programming/supported-solvers#fico-xpress
[pyomo-section]: /mixed-integer-programming/supported-solvers#pyomo
[nextmv-docs-apps]: https://www.nextmv.io/docs/platform/topic-guides/apps
[highs-options]: https://ergo-code.github.io/HiGHS/dev/options/definitions/
[go-parse-duration]: https://pkg.go.dev/time#ParseDuration
