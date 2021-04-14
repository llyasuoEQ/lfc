package lfc

import "fmt"

type Rule struct {
	RuleName  string
	Period    int64
	Threshold int64
	Code      int64
}

func newRule(ruleName string, period, threshold, code int64) Rule {
	return Rule{
		RuleName:  ruleName,
		Period:    period,
		Threshold: threshold,
		Code:      code,
	}
}

func (r *Rule) toString() string {
	if r == nil {
		return ""
	}
	result := fmt.Sprintf("RuleName:%v", r.RuleName)
	result += fmt.Sprintf(" Period:%v", r.Period)
	result += fmt.Sprintf(" Threshold:%v", r.Threshold)
	result += fmt.Sprintf(" Code:%v", r.Code)
	return result
}

func (r *Rule) getHash() string {
	if r == nil {
		return ""
	}
	return getMd5Bytes([]byte(r.toString()))
}
