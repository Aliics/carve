package main

type funcDef struct {
	paramNames []string
	instructs  []instruct
}

type runtime struct {
	funcDefs map[string]funcDef
	varDefs  map[string]value
}

func defaultRuntime() *runtime {
	return &runtime{
		map[string]funcDef{
			"print": {
				[]string{"value"},
				[]instruct{runBIFInstruct(
					bifPrint,
					[]value{varValue{"value"}},
				)},
			},
			"equals?":      makeCmpFunc(bifEquals),
			"greaterThan?": makeCmpFunc(bifGreaterThan),
			"lessThan?":    makeCmpFunc(bifLessThan),
			"contains?":    makeCmpFunc(bifContains),
			"concat":       makeCmpFunc(bifConcat),
			"plus":         makeCmpFunc(bifPlus),
			"minus":        makeCmpFunc(bifMinus),
		},
		make(map[string]value),
	}
}

func makeCmpFunc(bif bif) funcDef {
	return funcDef{
		[]string{"a", "b"},
		[]instruct{runBIFInstruct(
			bif,
			[]value{varValue{"a"}, varValue{"b"}},
		)},
	}
}

func exec(r *runtime, instructs []instruct) (last value, err error) {
	for _, i := range instructs {
		last, err = i(r)
		if err != nil {
			return nil, err
		}
	}

	return last, nil
}
