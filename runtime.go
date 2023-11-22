package main

type funcDef struct {
	paramNames []string
	instructs  []instruct
	level      int
}

type varDef struct {
	name  string
	value value
	level int
}

type runtime struct {
	funcDefs map[string]funcDef
	varDefs  []varDef
	level    int
}

func defaultRuntime() *runtime {
	return &runtime{
		funcDefs: map[string]funcDef{
			"print": {
				[]string{"value"},
				[]instruct{runBIFInstruct(
					bifPrint,
					[]value{varValue{"value"}},
				)},
				0,
			},
			"equals?":      makeCmpFunc(bifEquals),
			"greaterThan?": makeCmpFunc(bifGreaterThan),
			"lessThan?":    makeCmpFunc(bifLessThan),
			"contains?":    makeCmpFunc(bifContains),
			"concat":       makeCmpFunc(bifConcat),
			"plus":         makeCmpFunc(bifPlus),
			"minus":        makeCmpFunc(bifMinus),
		},
	}
}

func makeCmpFunc(bif bif) funcDef {
	return funcDef{
		[]string{"a", "b"},
		[]instruct{runBIFInstruct(
			bif,
			[]value{varValue{"a"}, varValue{"b"}},
		)},
		0,
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
