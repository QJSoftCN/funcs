package funcs

import (
	"strings"
	"regexp"
	"strconv"
)

const (
	rev_pattern = "(s[0-9]+\\.)*[a-z]+[0-9]+"
	dot         = "."
	underline   = "_"
	g_pattern="'?[\\w\\.]+'?"
)

//excel var regexp
var reg_excel_var = regexp.MustCompile(rev_pattern)
var reg_g = regexp.MustCompile(g_pattern)

//check funs
func Check(sheetVar, str string) string {
	str = strings.ToLower(str)
	str = checkExcelVars(sheetVar, str)
	for _, f := range need_checked_funs {
		str = ncf_map[f](str)
	}
	return str
}

func valid(v string) bool {
	vs:= reg_excel_var.FindAllString(v,-1)
	if len(vs)==1{
		return strings.Compare(vs[0],v)==0
	}else{
		return false
	}
}

func checkExcelVars(sheetVar, str string) string {
	vars := reg_g.FindAllString(str, -1)
	gvars := make(map[string]string)
	for index, v := range vars {
		//valid v is real excel var
		if !valid(v) {
			continue
		}
		key := "GV" + strconv.Itoa(index)
		if strings.Contains(v, dot) {
			//this is global var
			gvars[key] = strings.Replace(v, dot, underline, 1)
		} else {
			gvar := sheetVar + underline + v
			gvars[key] = gvar
		}
		str = strings.Replace(str, v, key, 1)
	}

	for k, v := range gvars {
		str = strings.Replace(str, k, v, 1)
	}

	return str
}

// common fun check
func common_check(f, str string) string {
	funs := find(str, f)
	for _, fun := range funs {
		nfun := editFun(fun, nil)
		str = strings.Replace(str, fun, nfun, 1)
	}

	return str
}

func editFun(fun string, eps []int) string {
	ps := parsePara(fun)
	if len(eps) == 0 {
		for _, p := range ps {
			fun = strings.Replace(fun, p, "\""+p+"\"", 1)
		}
	} else {
		for _, i := range eps {
			fun = strings.Replace(fun, ps[i], "\""+ps[i]+"\"", 1)
		}
	}
	return fun
}

func parsePara(fun string) []string {
	ps := make([]string, 0)
	start := strings.Index(fun, "(")

	body := fun[start+1:len(fun)-1]
	bodyLen := len(body)

	num := 0
	pStart := 0
	for index, s := range body {
		if s == '(' || s == '[' || s == '{' {
			num++
		}

		if s == ')' || s == ']' || s == '}' {
			num--
		}

		if s == ',' {
			if num == 0 {
				ps = append(ps, body[pStart:pStart+index])
				pStart += index + 1
			}
		}
	}

	if pStart < bodyLen {
		ps = append(ps, body[pStart:])
	}

	return ps

}

func find(str, f string) []string {
	funs := make([]string, 0)
	funs = matchFuns(str, f, 0, funs)
	return funs
}

func matchFuns(str, f string, start int, funs []string) []string {

	l := len(str)
	if start >= l-1 {
		//end
		return funs
	}

	sub := f + "("
	startIndex := strings.Index(str[start:], sub)
	if startIndex == -1 {
		//not find
		return funs
	}

	startIndex += start

	num := 1
	endIndex := startIndex + len(sub)
	isFind := false
	for _, s := range str[startIndex+len(sub):] {
		switch s {
		case '(':
			num++
		case ')':
			num--
		default:

		}

		endIndex++
		if num == 0 {
			isFind = true
			break
		}
	}

	if isFind {
		//ok
		funs = append(funs, str[startIndex:endIndex])
		return matchFuns(str, f, endIndex, funs)
	} else {
		funs = append(funs, str[startIndex:endIndex])
		return matchFuns(str, f, endIndex, funs)
	}
}

func isConst(str string) bool {
	if str[0] == '"' && str[len(str)-1] == '"' {
		return true
	} else {
		return false
	}
}
