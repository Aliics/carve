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
			"equals?": {
				[]string{"a", "b"},
				[]instruct{runBIFInstruct(
					bifEquals,
					[]value{varValue{"a"}, varValue{"b"}},
				)},
			},
			"greaterThan?": {
				[]string{"a", "b"},
				[]instruct{runBIFInstruct(
					bifGreaterThan,
					[]value{varValue{"a"}, varValue{"b"}},
				)},
			},
			"lessThan?": {
				[]string{"a", "b"},
				[]instruct{runBIFInstruct(
					bifLessThan,
					[]value{varValue{"a"}, varValue{"b"}},
				)},
			},
			"contains?": {
				[]string{"a", "b"},
				[]instruct{runBIFInstruct(
					bifContains,
					[]value{varValue{"a"}, varValue{"b"}},
				)},
			},
			"concat": {
				[]string{"a", "b"},
				[]instruct{runBIFInstruct(
					bifConcat,
					[]value{varValue{"a"}, varValue{"b"}},
				)},
			},
			"plus": {
				[]string{"a", "b"},
				[]instruct{runBIFInstruct(
					bifPlus,
					[]value{varValue{"a"}, varValue{"b"}},
				)},
			},
			"minus": {
				[]string{"a", "b"},
				[]instruct{runBIFInstruct(
					bifMinus,
					[]value{varValue{"a"}, varValue{"b"}},
				)},
			},
		},
		make(map[string]value),
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
