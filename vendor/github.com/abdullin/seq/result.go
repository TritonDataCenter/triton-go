package seq

import "fmt"

type Result struct {
	Issues []Issue
}

type Issue struct {
	Path          string
	ExpectedValue string
	ActualValue   string
}

func (r *Result) Ok() bool {
	return len(r.Issues) == 0
}

func (d *Issue) String() string {
	return fmt.Sprintf("Expected '%s' to be '%v' but got '%s'",
		d.Path,
		d.ExpectedValue,
		d.ActualValue,
	)
}

func NewResult() *Result {
	return &Result{}
}

func (r *Result) AddIssue(key, expected, actual string) {
	r.Issues = append(r.Issues, Issue{key, expected, actual})
}
