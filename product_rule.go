package lfc

import (
	"errors"
	"log"
	"sync"
)

var (
	globalProductRule *ProductRule
)

type ProductRule struct {
	ProductName string
	RuleSets    []RuleSet
	redisClient *RedisClient
}

func newProductRule(productName string, ruleSets []RuleSet) *ProductRule {
	return &ProductRule{
		ProductName: productName,
		RuleSets:    ruleSets,
		redisClient: globalRedisClient,
	}
}

// FrequencyControl...
// frequency control entrance
func (pr *ProductRule) FrequencyControl(input *Input) (result *MatchResult, err error) {
	err = pr.writeRedis(input)
	if err != nil {
		log.Printf("write redis failed:%v\n", err)
	}
	result, err = pr.readRedis(input)
	return
}

// writeRedis...
func (pr *ProductRule) writeRedis(input *Input) (err error) {
	if pr == nil || input == nil {
		err = errors.New("product rule input is nil")
		return
	}
	// write concurrently
	wg := sync.WaitGroup{}
	wg.Add(len(pr.RuleSets))
	for _, item := range pr.RuleSets {
		ruleSet := item
		FastRecoverGoroutineFunc(func() {
			defer wg.Done()
			err = ruleSet.WriteRedis(pr.ProductName, input)
			if err != nil {
				log.Printf("write redis failed:%v\n", err)
			}
		})
	}
	wg.Wait()
	return
}

// readRedis...
func (pr *ProductRule) readRedis(input *Input) (result *MatchResult, err error) {
	if pr == nil || input == nil {
		err = errors.New("ProductRule or input nil")
		return
	}
	upperBound := input.Ts
	for _, ruleSet := range pr.RuleSets {
		lowerBound := getLowerBound(upperBound, ruleSet.maxPeriod)
		var fcDataList FcDataList
		fcDataList, err = ruleSet.ReadRedis(pr.ProductName, input.Data, upperBound, lowerBound)
		if err != nil {
			return
		}
		for _, rule := range ruleSet.Rules {
			lowerRuleBound := getLowerBound(upperBound, rule.Period)
			var ruleDataList FcDataList
			for _, fcData := range fcDataList {
				if lowerRuleBound <= fcData.Score && fcData.Score <= upperBound {
					ruleDataList = append(ruleDataList, fcData)

					// determine whether the threshold is reached
					if ruleDataList.IsReachThreshold(rule.Threshold) {
						matchDetail := MatchDetail{rule.RuleName, rule.Period, rule.Threshold, ruleDataList}
						result = &MatchResult{Code: rule.Code, Data: matchDetail}
						return
					}
				}
			}
		}
	}
	return
}

func (pr *ProductRule) formatByFrequencyConfig(fc *Fconfig) (productRule *ProductRule, err error) {
	if fc != nil {
		// fields字段组合哈希
		fHashMap := make(map[string][]string)
		// 规则哈希
		rHashMap := make(map[string]Rule)
		// 字段哈希到规则哈希数组
		fRuleHashMap := make(map[string][]string)
		for _, item := range fc.Rules {
			// 进行俩次排序
			fields := item.getSortFields()
			fieldsHash := getStrListHash(fields)
			if len(fields) > 0 {
				if _, exists := fHashMap[fieldsHash]; !exists {
					fHashMap[fieldsHash] = fields
				}
				if _, exists := fRuleHashMap[fieldsHash]; !exists {
					fRuleHashMap[fieldsHash] = make([]string, 0)
				}
			}

			// 将滑动窗口的时间字符串转为毫秒级时间 int64类型
			var period int64
			period, err = timeStringToMilliSecond(item.Period)
			if err != nil {
				return
			}

			rule := newRule(item.Name, period, int64(item.Threshold), item.Code)
			ruleHash := rule.getHash()
			if _, exists := rHashMap[ruleHash]; !exists {
				rHashMap[ruleHash] = rule
			}

			if _, ok := fRuleHashMap[fieldsHash]; ok {
				if !containInStringSlice(fRuleHashMap[fieldsHash], ruleHash) {
					fRuleHashMap[fieldsHash] = append(fRuleHashMap[fieldsHash], ruleHash)
				}
			}
		}

		var ruleSets []RuleSet
		for k, v := range fRuleHashMap {
			var rules []Rule
			if fields, ok := fHashMap[k]; ok {
				for _, rHash := range v {
					if rule, exists := rHashMap[rHash]; exists {
						rules = append(rules, rule)
					}
				}
				if len(rules) > 0 {
					ruleSet := newRuleSet(fields, k, rules)
					ruleSets = append(ruleSets, ruleSet)
				}
			}
		}

		globalProductRule = newProductRule(fc.AppName, ruleSets)

		productRule = globalProductRule
	}
	return
}
