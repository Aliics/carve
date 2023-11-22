package main

import "strconv"

type value interface {
	show(*runtime) string
	unwrap(*runtime) value
}

type varValue struct {
	name string
}

func (v varValue) show(r *runtime) string {
	return r.varDefs[v.name].show(r)
}

func (v varValue) unwrap(r *runtime) value {
	return r.varDefs[v.name].unwrap(r)
}

type strValue string

func (s strValue) show(*runtime) string {
	return string(s)
}

func (s strValue) unwrap(*runtime) value {
	return s
}

type boolValue bool

func (b boolValue) show(*runtime) string {
	if b {
		return keywordTokenTrue
	} else {
		return keywordTokenFalse
	}
}

func (b boolValue) unwrap(*runtime) value {
	return b
}

type intValue int

func (i intValue) show(*runtime) string {
	return strconv.Itoa(int(i))
}

func (i intValue) unwrap(*runtime) value {
	return i
}
