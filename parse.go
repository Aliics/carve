package main

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

type keywordToken = string

const (
	keywordTokenFunction keywordToken = "function"
	keywordTokenIf       keywordToken = "if"
	keywordTokenThen     keywordToken = "then"
	keywordTokenElse     keywordToken = "else"
	keywordTokenAnd      keywordToken = "and"
	keywordTokenOr       keywordToken = "or"
	keywordTokenEnd      keywordToken = "end"
	keywordTokenTrue     keywordToken = "true"
	keywordTokenFalse    keywordToken = "false"
	keywordTokenNot      keywordToken = "not"
)

func parse(lines []string) ([]instruct, error) {
	var trimmedLines []string
	for _, l := range lines {
		trimmedLines = append(trimmedLines, strings.TrimSpace(l))
	}

	var instructs []instruct
	var skipCount int
	for i, l := range trimmedLines {
		if skipCount > 0 {
			skipCount--
			continue
		}

		if l == "" || strings.HasPrefix(l, "--") {
			continue
		}

		if strings.HasPrefix(l, keywordTokenFunction) {
			funcInstructs, endLine, err := parseFuncDef(i, l, trimmedLines)
			if err != nil {
				return nil, err
			}

			instructs = append(instructs, funcInstructs...)
			skipCount = endLine

			continue
		}

		if strings.HasPrefix(l, keywordTokenIf) && strings.HasSuffix(l, keywordTokenThen) {
			ifInstruct, endLine, err := parseIfStatement(i, l, trimmedLines)
			if err != nil {
				return nil, err
			}
			instructs = append(instructs, ifInstruct)
			skipCount = endLine

			continue
		}

		if textLength, err := getTextTokenLength(l); err == nil {
			varExpr, isAssignment := strings.CutPrefix(
				strings.TrimSpace(l[textLength:]),
				"=",
			)

			if textLength == len(l) {
				valueInstruct, err := parseExpr(l)
				if err != nil {
					return nil, err
				}

				instructs = append(instructs, valueInstruct)
			} else if isAssignment {
				exprInstruct, err := parseExpr(varExpr)
				if err != nil {
					return nil, err
				}

				instructs = append(instructs, assignVarInstruct(l[:textLength], exprInstruct))
			} else {
				callInstruct, err := parseFuncCall(l, textLength)
				if err != nil {
					return nil, err
				}

				instructs = append(instructs, callInstruct)
			}

			continue
		}
	}

	return instructs, nil
}

func parseIfStatement(
	lineNum int,
	line string,
	lines []string,
) (instruct, int, error) {
	exprStr := line[len(keywordTokenIf)+1 : len(line)-len(keywordTokenThen)-1]
	endLine, err := getBodyLineLength(lines[lineNum:])
	if err != nil {
		return nil, -1, err
	}

	elseLine := slices.Index(lines[lineNum+1:lineNum+endLine], keywordTokenElse)
	var (
		ifInstructs   []instruct
		elseInstructs []instruct
	)
	if elseLine == -1 {
		ifInstructs, err = parse(lines[lineNum+1 : lineNum+endLine])
		if err != nil {
			return nil, -1, err
		}
	} else {
		ifInstructs, err = parse(lines[lineNum+1 : lineNum+elseLine+1])
		if err != nil {
			return nil, -1, err
		}
		elseInstructs, err = parse(lines[lineNum+elseLine+2 : lineNum+endLine])
		if err != nil {
			return nil, -1, err
		}
	}

	condExpr, err := parseExpr(exprStr)
	if err != nil {
		return nil, -1, err
	}

	return evaluateIfInstruct(condExpr, ifInstructs, elseInstructs), endLine, nil
}

func parseFuncDef(
	lineNum int,
	line string,
	lines []string,
) ([]instruct, int, error) {
	var instructs []instruct
	skip, err := getBodyLineLength(lines[lineNum:])
	if err != nil {
		return nil, -1, err
	}

	nameLength, err := getTextTokenLength(line[len(keywordTokenFunction)+1:])
	if err != nil {
		return nil, -1, err
	}
	nameEnd := len(keywordTokenFunction) + nameLength + 1
	funcName := line[len(keywordTokenFunction)+1 : nameEnd]

	paramNames := slices.DeleteFunc(
		strings.Split(
			strings.ReplaceAll(line[nameEnd+1:len(line)-1], " ", ""),
			",",
		),
		func(s string) bool {
			return s == ""
		},
	)

	funcBody, err := parse(lines[lineNum+1 : lineNum+skip])
	if err != nil {
		return nil, -1, err
	}

	def := funcDef{paramNames, funcBody, 1}
	instructs = append(
		instructs,
		defineFuncInstruct(funcName, def),
	)

	return instructs, skip, nil
}

func getBodyLineLength(lines []string) (int, error) {
	skip := -1
	depth := 0
	for i, l := range lines {
		if strings.HasPrefix(l, keywordTokenFunction) || strings.HasPrefix(l, keywordTokenIf) {
			depth++
		}

		if l == keywordTokenEnd {
			depth--

			if depth == 0 {
				skip = i
				break
			}
		}
	}
	if skip == -1 {
		return 0, errors.New(`function was not closed by "end"`)
	}

	return skip, nil
}

func parseFuncCall(line string, funcNameLength int) (instruct, error) {
	funcName := line[:funcNameLength]
	if len(line) == funcNameLength || line[funcNameLength] != '(' {
		return nil, fmt.Errorf(`expected "(" after %s`, funcName)
	}

	exprInParens, hadClose := strings.CutSuffix(line, ")")
	if !hadClose {
		return nil, fmt.Errorf(`expected ")"`)
	}

	var instructs []instruct
	exprStrs := splitWithWrappingContext(exprInParens[funcNameLength+1:], ",")
	for _, exprStr := range exprStrs {
		i, err := parseExpr(exprStr)
		if err != nil {
			return nil, err
		}

		instructs = append(instructs, i)
	}

	return invokeFuncInstruct(funcName, instructs), nil
}

func parseExpr(exprStr string) (instruct, error) {
	var ors []instruct
	orStrs := splitWithWrappingContext(exprStr, keywordTokenOr)
	for _, orExpr := range orStrs {
		var ands []instruct
		andStrs := splitWithWrappingContext(orExpr, keywordTokenAnd)
		for _, andStr := range andStrs {
			expr, isNegated := strings.CutPrefix(
				strings.TrimSpace(andStr),
				keywordTokenNot+" ",
			)

			var exprInstruct instruct
			if expr == keywordTokenTrue {
				exprInstruct = valueExprInstruct(boolValue(true))
			} else if expr == keywordTokenFalse {
				exprInstruct = valueExprInstruct(boolValue(false))
			} else if unicode.IsDigit(rune(expr[0])) {
				n, err := strconv.Atoi(expr)
				if err != nil {
					return nil, err
				}
				exprInstruct = valueExprInstruct(intValue(n))
			} else if textLength, err := getTextTokenLength(expr); err == nil {
				if textLength == len(expr) {
					exprInstruct = valueExprInstruct(varValue{expr})
				} else {
					callInstruct, err := parseFuncCall(expr, textLength)
					if err != nil {
						return nil, err
					}

					exprInstruct = callInstruct
				}
			} else if strings.HasPrefix(expr, `"`) &&
				strings.HasSuffix(expr, `"`) {
				exprInstruct =
					valueExprInstruct(strValue(strings.Trim(expr, `"`)))
			} else {
				return nil, errors.New("expression expected")
			}

			if isNegated {
				exprInstruct = notExprInstruct(exprInstruct)
			}

			ands = append(ands, exprInstruct)
		}

		if len(ands) < 2 {
			ors = append(ors, ands...)
		} else {
			ors = append(ors, andExprInstruct(ands))
		}
	}

	if len(ors) < 2 {
		return ors[0], nil
	}

	return orExprInstruct(ors), nil
}

func splitWithWrappingContext(s, sep string) []string {
	var strs []string
	var strOpen bool
	var parenOpen bool
	var acc string
	for i := 0; i < len(s); i++ {
		if s[i] == '"' {
			strOpen = !strOpen
		}
		if !strOpen && s[i] == '(' || s[i] == ')' {
			parenOpen = !parenOpen
		}

		if !strOpen && !parenOpen && i < len(s)-len(sep) && s[i:i+len(sep)] == sep {
			strs = append(strs, acc)
			acc = ""
			i += len(sep)
		}

		acc += string(s[i])

		if i == len(s)-1 {
			strs = append(strs, strings.TrimSuffix(acc, sep))
		}
	}

	return strs
}

func getTextTokenLength(s string) (int, error) {
	var length int

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '?' {
			break
		}
		length++
	}

	if length == 0 {
		return -1, errors.New("token cannot be empty")
	}

	return length, nil
}
