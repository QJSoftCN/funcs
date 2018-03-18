package funcs

import (
	"github.com/robertkrimen/otto"
	"log"
	"regexp"
	"errors"
	"time"
)

const (
	Context     = "Context"
	var_pattern = "s[0-9]+_[a-z]+[0-9]+"
)

type CalcErr struct {
	Key     string
	Err     error
	Formula string
	Slove   string
}


type CalcResult interface {
	ValText() string
	Errs() []CalcErr
}

type FuncContext struct {
	timeCtx   *TimeContext
	vm        *otto.Otto
	exPeriods *ExcluedPeriods
	results   map[string]interface{}
	formulas  map[string]string
	calcErrs  []CalcErr
}

func (this *FuncContext) AddConst(key string, val interface{}) {
	this.vm.Set(key, val)
	this.results[key] = val
}

func (this *FuncContext) AddFormula(key, f string) {
	this.formulas[key] = f
}

func (this *FuncContext) RunFormula(key string) {
	f, ok := this.formulas[key]
	if ok {
		this.calcFormula(key, f)
	}
}

func setFuncs(vm *otto.Otto) {
	// load common funcs
	vm.Set("avg", avg)
	vm.Set("max", max)
	vm.Set("min", min)
	// load point funcs
	vm.Set(FUNC_PC, pc)
	vm.Set(FUNC_PE, pe)
	vm.Set(FUNC_FMT, fmt_fun)

}

func buildVM() *otto.Otto {
	var vm = otto.New()
	setFuncs(vm)
	return vm
}

func buidResults() map[string]interface{} {
	rs := make(map[string]interface{})
	return rs
}


func NewFuncContext(start, end time.Time) FuncContext {

	var ctx = FuncContext{}
	ctx.timeCtx = NewTimeContext(start, end)
	ctx.vm = buildVM()
	ctx.vm.Set(Context, &ctx)

	ctx.results = buidResults()

	fs := make(map[string]string)
	ctx.formulas = fs
	ces := make([]CalcErr, 0)
	ctx.calcErrs = ces
	return ctx
}


func NewCalcErr(key, f string, err error, slove string) CalcErr {
	ce := CalcErr{}
	ce.Key = key
	ce.Err = err
	ce.Formula = f
	ce.Slove = slove
	return ce
}

func (this *FuncContext) addErr(ce CalcErr) {
	this.calcErrs = append(this.calcErrs, ce)
}

var reg = regexp.MustCompile(var_pattern)

func (this *FuncContext) calcFormula(k, f string) {
	_, ok := this.results[k]
	if ok {
		return
	}
	vs := reg.FindAllString(f, -1)
	for _, v := range vs {
		_, ok := this.results[v]
		if ok {
			continue
		} else {
			vf, ok := this.formulas[v]
			if ok {
				this.calcFormula(v, vf)
			} else {
				//not exist
				this.vm.Set(v, otto.NullValue())
				this.results[v] = nil

				ce := NewCalcErr(k, f, errors.New(v+" is null"), "")
				this.addErr(ce)
			}
		}
	}

	ret, err := this.vm.Run(f)
	if err != nil {
		ce := NewCalcErr(k, f, err, "")
		this.addErr(ce)
	}

	this.vm.Set(k, ret)
	this.putResult(k, ret)

}

func (this *FuncContext) putResult(k string, ret otto.Value) {
	if ret.IsBoolean() {
		b, _ := ret.ToBoolean()
		this.results[k] = b
		return
	}

	if ret.IsNumber() {
		b, _ := ret.ToFloat()
		this.results[k] = b
		return
	}

	if ret.IsString() {
		b, _ := ret.ToString()
		this.results[k] = b
		return
	}

	//unkown type
	this.results[k] = ret
	log.Println("calc cell ", k, " unkown type,val is ", ret)

}

func (this *FuncContext) GetCalcErrs() *[]CalcErr {
	return &this.calcErrs
}

func (this *FuncContext) GetResult(key string) interface{} {
	return this.results[key]
}

func getContext(vm *otto.Otto) *FuncContext {
	ctx, err := vm.Get(Context)
	if err != nil {
		log.Println(err)
	}

	c, err := ctx.Export()
	if err != nil {
		log.Println(err)
	}

	if ctx, ok := c.(*FuncContext); ok {
		return ctx
	}

	//返回默认
	return nil

}

var need_checked_funs []string
var ncf_map map[string]func(fun string) string

func init() {
	ncf_map = make(map[string](func(fun string) string))
	need_checked_funs = append(need_checked_funs, FUNC_FMT, FUNC_PE)

	ncf_map[FUNC_FMT] = fmt_check
	ncf_map[FUNC_PE] = pe_check
	ncf_map[FUNC_PC] = pc_check
}
