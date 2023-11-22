package main

import (
	"fmt"
	"strconv"
)

type value interface {
	show(*runtime) string
	unwrap(*runtime) (value, error)
}

type varValue struct {
	name string
}

func (v varValue) show(r *runtime) string {
	unwrap, err := v.unwrap(r)
	if err != nil {
		return "null"
	}
	return unwrap.show(r)
}

func (v varValue) unwrap(r *runtime) (value, error) {
	var foundDef varDef
	for i := len(r.varDefs) - 1; i >= 0; i-- {
		def := r.varDefs[i]
		if def.name == v.name {
			foundDef = def
			break
		}
	}

	if foundDef == (varDef{}) {
		return nil, fmt.Errorf("%s is not defined", v.name)
	}

	if foundDef.level > r.level {
		return nil, fmt.Errorf("%s is not reachable on this stack", v.name)
	}

	return foundDef.value.unwrap(r)
}

type strValue string

func (s strValue) show(*runtime) string {
	return string(s)
}

func (s strValue) unwrap(*runtime) (value, error) {
	return s, nil
}

type boolValue bool

func (b boolValue) show(*runtime) string {
	if b {
		return keywordTokenTrue
	} else {
		return keywordTokenFalse
	}
}

func (b boolValue) unwrap(*runtime) (value, error) {
	return b, nil
}

type intValue int

func (i intValue) show(*runtime) string {
	return strconv.Itoa(int(i))
}

func (i intValue) unwrap(*runtime) (value, error) {
	return i, nil
}
