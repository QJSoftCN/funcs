package funcs

import (
	"github.com/robertkrimen/otto"
	"fmt"
)

const FUNC_PC  ="pc"

func pc_check(fun string)string{
	fun=common_check(FUNC_PC,fun)
	return fun
}
//pc('data','time','con')
//pc('data','time')
func pc(call otto.FunctionCall) otto.Value {

	context := getContext(call.Otto)

	switch len(call.ArgumentList) {
	case 2:
		//pc(data,time)
		return pointCalc2(context, call.Argument(0).String(), call.Argument(1).String())
	case 3:
		//pc(data,time,con)
		return pointCalc3(context, call.Argument(0).String(), call.Argument(1).String(), call.Argument(2).String())
	default:
		return pointCalc2(context, call.Argument(0).String(), call.Argument(1).String())
	}

	return createErr("pc err")
}

func pointCalc3(ctx *FuncContext, data, time, con string) otto.Value {

	fmt.Println(ctx, data, time, con)

	ret, _ := otto.ToValue(3)
	return ret
}

func pointCalc2(ctx *FuncContext, data, time string) otto.Value {
	fmt.Println(ctx, data, time)
	ret, _ := otto.ToValue(2)
	return ret
}
