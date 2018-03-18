package funcs

import (
	"github.com/robertkrimen/otto"
	"strings"
	"time"
	"../points"
)

const FUNC_PE = "pe"

func pe_check(fun string) string {
	fun = common_check(FUNC_PE, fun)
	return fun
}

type PEResult struct {
	Body     string
	Paras    []string
	Formula  string
	CalcErrs []CalcErr
	Val      interface{}
	CanUse   bool
}

//pe(exp)
//pe(exp,t)
func pe(call otto.FunctionCall) otto.Value {
	context := getContext(call.Otto)
	switch len(call.ArgumentList) {
	case 1:
		//pexp(exp)
		return pexp(context, call.Argument(0).String(), "")
	case 2:
		//pexp(exp,time)
		return pexp(context, call.Argument(0).String(),
			call.Argument(1).String())
	default:
		return pexp(context, call.Argument(0).String(), "")
	}

	return createErr(FUNC_PE + " err")
}

func calcPeResult(ctx *FuncContext, exp, timeExp string) PEResult {
	rKey := FuncBody(FUNC_PE, exp, timeExp)
	r, ok := ctx.results[rKey]
	if ok {
		per, _ := r.(PEResult)
		return per
	}

	per := PEResult{}
	per.Body = rKey
	per.Paras = make([]string, 1)
	per.Paras[0] = exp
	per.CalcErrs = make([]CalcErr, 0)

	var et *time.Time
	if len(timeExp) > 0 {
		per.Paras = append(per.Paras, timeExp)
		t, err := ctx.timeCtx.Parse(timeExp)
		if err != nil {
			ce := NewCalcErr(rKey, timeExp, err, "采用当前值替换计算")
			per.CalcErrs = append(per.CalcErrs, ce)
			ctx.addErr(ce)
		}

		et = t
	}

	vars := GetPointVars(exp)
	f := exp
	for _, v := range vars {
		val, err := calcVar(v, et)
		if err != nil {
			ce := NewCalcErr(rKey, v, err, "")
			per.CalcErrs = append(per.CalcErrs, ce)
			ctx.addErr(ce)
		} else {
			f = strings.Replace(f, v, val, 1)
		}
	}

	per.Formula = f
	fr, err := ctx.vm.Run(f)
	if err != nil {
		ce := NewCalcErr(rKey, f, err, "")
		per.CalcErrs = append(per.CalcErrs, ce)
		ctx.addErr(ce)
		per.CanUse = false
	} else {
		per.CanUse = true
		if fr.IsNumber() {
			per.Val, _ = fr.ToFloat()
		}
		if fr.IsString() {
			per.Val, _ = fr.ToString()
		}
		if fr.IsBoolean() {
			per.Val, _ = fr.ToBoolean()
		}
	}

	return per
}

func pexp(ctx *FuncContext, exp, timeExp string) otto.Value {
	result := calcPeResult(ctx, exp, timeExp)
	if result.CanUse {
		switch result.Val.(type) {
		case bool:
			f, _ := result.Val.(bool)
			ret, _ := otto.ToValue(f)
			return ret
		case float64:
			f, _ := result.Val.(float64)
			ret, _ := otto.ToValue(f)
			return ret
		default:
			f, _ := result.Val.(string)
			ret, _ := otto.ToValue(f)
			return ret
		}

	} else {
		return otto.NullValue()
	}
}

func calcVar(v string, t *time.Time) (string, error) {

	pvar := points.ParseVar(v)

	var pv *points.PointValue
	var err error

	if t == nil {
		pv, err = points.ReadSnapshot(pvar.PointName)
	} else {
		pv, err = points.InterVal(pvar.PointName, *t, pvar.GetInterWay())
	}

	if err != nil {
		return v, err
	}

	switch pvar.GetMethod() {
	case points.PVCM_Time:
		return pv.GetTimeString(), err
	default:
		return pv.ValToString(), err
	}

}
