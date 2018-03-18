package funcs

import (
	"github.com/robertkrimen/otto"
	"strings"
	"fmt"
	"time"
	"strconv"
	"github.com/qjsoftcn/gutils"
)

const FUNC_FMT = "fmt"

func fmt_check(str string) string {
	funs := find(str, FUNC_FMT)
	eps := []int{1}
	for _, fun := range funs {
		nfun := editFun(fun, eps)
		str = strings.Replace(str, fun, nfun, 1)
	}
	return str
}

//fmt(val,fmt)
func fmt_fun(call otto.FunctionCall) otto.Value {
	context := getContext(call.Otto)
	switch len(call.ArgumentList) {
	case 2:
		//pexp(exp,time)
		return fmt_func(context, call.Argument(0).String(),
			call.Argument(1).String())
	default:
		return fmt_func(context, call.Argument(0).String(), "")
	}

	return createErr(FUNC_FMT + " err")
}

func fmtFloat(f float64, valFmt string) string {
	_, err := strconv.Atoi(valFmt)
	if err == nil {
		return fmt.Sprintf("%."+valFmt+"f", f)
	} else {
		if len(valFmt) == 2 {
			if valFmt[1] == '%' {
				n, err := strconv.Atoi(string(valFmt[0]))
				if err != nil {
					n = 2
				}

				return fmt.Sprintf("%."+strconv.Itoa(n)+"f", f) + "%"
			}
		}

		return fmt.Sprintf("%.2f", f)

	}
}



func fmtString(s string, valFmt string) string {
	//float
	f, err := strconv.ParseFloat(s, -1)
	if err == nil {
		return fmtFloat(f, valFmt)
	}

	//time
	t, err :=  time.ParseInLocation(gutils.ToGoTimeFmt(valFmt), s, time.Local)
	if err == nil {
		return gutils.Format(t, valFmt)
	}

	return s
}

func fmtResult(ctx *FuncContext, val interface{}, valFmt string) string {
	fr := ""
	switch val.(type) {
	case float64:
		f, _ := val.(float64)
		fr = fmtFloat(f, valFmt)
	case int64:
		i, _ := val.(float64)
		fr = fmt.Sprint(i)
	case string:
		s, _ := val.(string)
		fr = fmtString(s, valFmt)
	case time.Time:
		t, _ := val.(time.Time)
		fr = gutils.Format(t, valFmt)
	default:
		fr = fmt.Sprint(val)
	}

	return fr
}

func fmt_func(ctx *FuncContext, val interface{}, valFmt string) otto.Value {
	result := fmtResult(ctx, val, valFmt)
	r, _ := otto.ToValue(result)
	return r
}
