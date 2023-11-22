package main

import (
	"errors"
	"fmt"
)

type instruct func(*runtime) (value, error)

func defineFuncInstruct(name string, def funcDef) instruct {
	return func(r *runtime) (value, error) {
		r.funcDefs[name] = def
		return nil, nil
	}
}

func invokeFuncInstruct(name string, args []instruct) instruct {
	return func(r *runtime) (value, error) {
		def, ok := r.funcDefs[name]
		if !ok {
			return nil, fmt.Errorf("%s function is not defined", name)
		}

		if len(args) != len(def.paramNames) {
			return nil, fmt.Errorf(
				"arg count mismatch %d != %d",
				len(args), len(def.paramNames),
			)
		}

		for i, name := range def.paramNames {
			v, err := args[i](r)
			if err != nil {
				return nil, err
			}
			r.varDefs[name] = v
		}

		return exec(r, def.instructs)
	}
}

func evaluateIfInstruct(
	condition instruct,
	ifInstructs []instruct,
	elseInstructs []instruct,
) instruct {
	return func(r *runtime) (value, error) {
		cond, err := exec(r, []instruct{condition})
		if err != nil {
			return nil, err
		}

		boolValue, ok := cond.unwrap(r).(boolValue)
		if !ok {
			return nil, errors.New("condition expression was not a bool")
		}

		if boolValue {
			return exec(r, ifInstructs)
		} else if len(elseInstructs) > 0 {
			return exec(r, elseInstructs)
		}

		return nil, nil
	}
}

func assignVarInstruct(name string, i instruct) instruct {
	return func(r *runtime) (value, error) {
		v, err := i(r)
		if err != nil {
			return nil, err
		}

		r.varDefs[name] = v

		return nil, nil
	}
}

func valueExprInstruct(v value) instruct {
	return func(r *runtime) (value, error) {
		return v, nil
	}
}

func notExprInstruct(i instruct) instruct {
	return func(r *runtime) (value, error) {
		b, err := ensuringType[boolValue](r, i)
		if err != nil {
			return nil, err
		}

		return !b, nil
	}
}

func orExprInstruct(instructs []instruct) instruct {
	return func(r *runtime) (value, error) {
		for _, i := range instructs {
			b, err := ensuringType[boolValue](r, i)
			if err != nil {
				return nil, err
			}

			if b {
				return b, nil
			}
		}
		return boolValue(false), nil
	}
}

func andExprInstruct(instructs []instruct) instruct {
	return func(r *runtime) (value, error) {
		result := true
		for _, i := range instructs {
			b, err := ensuringType[boolValue](r, i)
			if err != nil {
				return nil, err
			}

			result = result && b == true
		}

		return boolValue(result), nil
	}
}

func ensuringType[V value](r *runtime, i instruct) (val V, err error) {
	v, err := i(r)
	if err != nil {
		return val, err
	}

	var ok bool
	val, ok = v.(V)
	if !ok {
		return val, errors.New("expected bool")
	}
	return val, nil
}
