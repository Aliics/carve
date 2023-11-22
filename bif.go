package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type bif uint

const (
	bifPrint bif = iota
	bifEquals
	bifGreaterThan
	bifLessThan
	bifContains
	bifConcat
	bifPlus
	bifMinus
)

func runBIFInstruct(n bif, args []value) instruct {
	switch n {
	case bifPrint:
		return func(r *runtime) (value, error) {
			var strs []string
			for _, arg := range args {
				strs = append(strs, arg.show(r))
			}

			fmt.Println(strings.Join(strs, " "))
			return nil, nil
		}
	case bifEquals:
		return func(r *runtime) (value, error) {
			first, second, err := ensureTwoSameType(r, args)
			if err != nil {
				return nil, err
			}

			return boolValue(first == second), nil
		}
	case bifGreaterThan, bifLessThan:
		return func(r *runtime) (value, error) {
			first, second, err := ensureTwoOfType[intValue](r, args)
			if err != nil {
				return nil, err
			}

			if n == bifGreaterThan {
				return boolValue(first > second), err
			} else {
				return boolValue(first < second), err
			}
		}
	case bifContains:
		return func(r *runtime) (value, error) {
			first, second, err := ensureTwoOfType[strValue](r, args)
			if err != nil {
				return nil, err
			}

			return boolValue(strings.Contains(string(first), string(second))), nil
		}
	case bifConcat:
		return func(r *runtime) (value, error) {
			first, second, err := ensureTwoOfType[strValue](r, args)
			if err != nil {
				return nil, err
			}

			return first + second, nil
		}
	case bifPlus, bifMinus:
		return func(r *runtime) (value, error) {
			first, second, err := ensureTwoOfType[intValue](r, args)
			if err != nil {
				return nil, err
			}

			if n == bifPlus {
				return first + second, nil
			} else {
				return first - second, nil
			}
		}
	default:
		return func(r *runtime) (value, error) {
			return nil, errors.New("unknown function")
		}
	}
}

func ensureTwoSameType(r *runtime, args []value) (first, second value, err error) {
	if len(args) != 2 {
		return nil, nil, errors.New("2 arguments expected for equals")
	}

	first, err = args[0].unwrap(r)
	if err != nil {
		return nil, nil, err
	}
	second, err = args[1].unwrap(r)
	if err != nil {
		return nil, nil, err
	}

	if reflect.TypeOf(first) != reflect.TypeOf(second) {
		return nil, nil, fmt.Errorf(
			"cannot compare %s to %v",
			args[0].show(r), args[1].show(r),
		)
	}

	return
}

func ensureTwoOfType[V value](r *runtime, args []value) (first, second V, err error) {
	if len(args) != 2 {
		err = errors.New("2 arguments expected")
		return
	}

	var ok bool

	val, err := args[0].unwrap(r)
	if err != nil {
		return
	}
	first, ok = val.(V)
	if !ok {
		err = fmt.Errorf("first argument must be a %T", first)
		return
	}

	val, err = args[1].unwrap(r)
	if err != nil {
		return
	}
	second, ok = val.(V)
	if !ok {
		err = fmt.Errorf("second argument must be a %T", second)
		return
	}

	return
}
