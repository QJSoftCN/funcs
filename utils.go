package funcs

import (
	"regexp"
	"github.com/robertkrimen/otto"
)

const (
	pointVar_Pattern = `'[^']+'`
)

var (
	//exp
	reg_pvar = regexp.MustCompile(pointVar_Pattern)
)

func GetPointVars(para string) []string {
	//exp
	vars := reg_pvar.FindAllString(para, -1)
	return vars
}

func ToPointVar(para string) string {
	//exp
	return "'" + para + "'"
}

func FuncBody(funcName string, paras ... string) string {
	key := funcName + "("

	pLen := len(paras)

	for index, para := range paras {
		if para == "" {
			continue
		}
		key += para
		if index != pLen-1 {
			key += ","
		}

	}

	key += ")"
	return key
}

// 创建错误数据
func createErr(msg string) otto.Value {
	errValue, _ := otto.ToValue("Err:" + msg)
	return errValue
}
