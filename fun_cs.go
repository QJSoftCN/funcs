 package funcs

import (
	"github.com/robertkrimen/otto"
	"log"
)

//max(a,b,c)
func min(call otto.FunctionCall)otto.Value{

	min:=0.0
	num:=0

	for index,arg:=range call.ArgumentList{

		if arg.IsNumber(){
			argValue,err:=arg.ToFloat()
			if err!=nil{
				log.Println(call,index,arg,err)
				continue
			}

			if num==0{
				min=argValue
			}else{
				if argValue<min{min=argValue}
			}

			num++
		}
	}

	if num>0{
		minValue,_:=otto.ToValue(min)
		return minValue
	}else{
		return createErr("no good para")
	}
}

//max(a,b,c)
func max(call otto.FunctionCall)otto.Value{

	max:=0.0
	num:=0

	for index,arg:=range call.ArgumentList{

		if arg.IsNumber(){
			argValue,err:=arg.ToFloat()
			if err!=nil{
				log.Println(call,index,arg,err)
				continue
			}

			if num==0{
				max=argValue
			}else{
				if argValue>max{max=argValue}
			}

			num++
		}
	}

	if num>0{
		maxValue,_:=otto.ToValue(max)
		return maxValue
	}else{
		return createErr("max no good number")
	}
}


//avg(a,b,c)
func avg(call otto.FunctionCall)otto.Value{

	sum:=0.0
	num:=0

	for index,arg:=range call.ArgumentList{

		if arg.IsNumber(){
			argValue,err:=arg.ToFloat()
			if err!=nil{
				log.Println(call,index,arg,err)
				continue
			}

			sum+=argValue
			num++
		}
	}

	if num>0{
		avg:=sum/float64(num)
		avgValue,_:=otto.ToValue(avg)
		return avgValue
	}else{
		return createErr("avg no good number")
	}
}