// Package projects allows managing Go projects.
package projects

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"
)

//go:embed projects.txt
var file []byte

// Project is a project name.
type Project string

// List is the list of projects parsed from projects.txt.
var List = MustParse(file)

// Parse parses and validates project names from data.
//
// Names are read one per line. Blank lines are ignored and surrounding whitespace is trimmed.
//
// It returns an error if the result is empty, if a name is duplicated, or if the names are not sorted in ascending order.
func Parse(data []byte) ([]Project, error) {
	lines := strings.Split(string(data), "\n")
	var ps []Project
	for i, line := range lines {
		name := strings.TrimSpace(line)
		if name == "" {
			continue
		}
		if len(ps) > 0 {
			prev := string(ps[len(ps)-1])
			if name == prev {
				return nil, fmt.Errorf("duplicate project %q at line %d", name, i+1)
			}
			if name < prev {
				return nil, fmt.Errorf("project %q at line %d is not sorted (previous %q)", name, i+1, prev)
			}
		}
		ps = append(ps, Project(name))
	}
	if len(ps) == 0 {
		return nil, errors.New("no projects")
	}
	return ps, nil
}

// MustParse calls [Parse] and panics if there is an error.
func MustParse(data []byte) []Project {
	ps, err := Parse(data)
	if err != nil {
		panic(err)
	}
	return ps
}
