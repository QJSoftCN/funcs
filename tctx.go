package funcs

import (
	"time"
	"github.com/qjsoftcn/texp"
)

type TimeContext struct {
	Start  time.Time
	End    time.Time
	parser *texp.TimeExpParser
}

func(this *TimeContext)Parse(timeExp string)(*time.Time,error){
	t,err:=this.parser.Parse(timeExp)
	return t,err
}

func NewTimeContext(start, end time.Time) *TimeContext {
	var tCtx = &TimeContext{}
	tCtx.Start = start
	tCtx.End = end
	tCtx.parser = texp.NewParser(start, end)
	return tCtx
}





