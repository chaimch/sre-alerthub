package common

import "strings"

// 逆波兰表达式实现, 主要是检验表达式是否语法正确
func CalculateReversePolishNotation(labelmap map[string]string, expression string) bool {
	stack := []bool{}
	exp := strings.Split(expression, ` `)
	for i := 0; i < len(exp); i++ {
		if exp[i] == "&" {
			var v1, v2 bool
			v2, stack = stack[len(stack)-1], stack[:len(stack)-1]
			v1, stack = stack[len(stack)-1], stack[:len(stack)-1]
			stack = append(stack, v1 && v2)
		} else if exp[i] == "|" {
			var v1, v2 bool
			v2, stack = stack[len(stack)-1], stack[:len(stack)-1]
			v1, stack = stack[len(stack)-1], stack[:len(stack)-1]
			stack = append(stack, v1 || v2)
		} else {
			if strings.Contains(exp[i], "!") {
				val := strings.Split(exp[i], "!=")
				if _, ok := labelmap[val[0]]; ok {
					if labelmap[val[0]] == val[1] {
						stack = append(stack, false)
					} else {
						stack = append(stack, true)
					}
				} else {
					return false
				}
			} else {
				val := strings.Split(exp[i], "=")
				if _, ok := labelmap[val[0]]; ok {
					if labelmap[val[0]] == val[1] {
						stack = append(stack, true)
					} else {
						stack = append(stack, false)
					}
				} else {
					return false
				}
			}
		}
	}
	return stack[len(stack)-1]
}

var Recover2Send = map[string]map[[2]int64]*Ready2Send{
	"LANXIN": map[[2]int64]*Ready2Send{},
	//"HOOK":   map[[2]int64]*Ready2Send{},
}
