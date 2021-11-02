package lfc

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/redis.v5"
)

type RuleSet struct {
	Fields       []string
	fieldsHash   *string
	maxPeriod    int64
	maxThreshold int64
	Rules        []Rule
}

func newRuleSet(fields []string, fieldHash string, rules []Rule) RuleSet {
	object := RuleSet{
		Fields: fields,
		Rules:  rules,
	}
	if fieldHash != "" {
		object.fieldsHash = &fieldHash
	}
	(&object).format()
	return object
}

func (rs *RuleSet) format() {
	for _, rule := range rs.Rules {
		if rule.Period > rs.maxPeriod {
			rs.maxPeriod = rule.Period
		}
		if rule.Threshold > rs.maxThreshold {
			rs.maxThreshold = rule.Threshold
		}
	}
}

func (ruleSet *RuleSet) getCurrentFields() string {
	return strings.Join(ruleSet.Fields, redis_key_inner_split)
}

func (rs *RuleSet) isFit(inputData map[string]interface{}) bool {
	if rs == nil {
		return false
	}

	for _, v := range rs.Fields {
		if _, ok := inputData[v]; !ok {
			return false
		}
	}
	return true
}

func (rs *RuleSet) prepareRedisKey(productName string, inputData map[string]interface{}) string {
	var result string
	var buf bytes.Buffer

	// 拼凑fcData
	normalizedData := make([]string, len(rs.Fields))
	for index, field := range rs.Fields {
		normalizedData[index] = fmt.Sprintf("%v", inputData[field])
	}
	fcData := strings.Join(normalizedData, redis_key_inner_split)

	//pre:productName:fields:maxPeriod:inputData
	buf.WriteString(redis_key_prefix)
	buf.WriteString(redis_key_split)
	buf.WriteString(productName)
	buf.WriteString(redis_key_split)
	buf.WriteString(rs.getCurrentFields()) // sort field string
	buf.WriteString(redis_key_split)
	buf.WriteString(strconv.FormatInt(rs.maxPeriod, 10)) //最大保留时长
	buf.WriteString(redis_key_split)
	buf.WriteString(fcData)

	// hash the key
	result = getMd5Bytes(buf.Bytes())

	//result = buf.String()
	return result

}

func (rs *RuleSet) WriteRedis(produceName string, input *Input) (err error) {
	if rs == nil {
		err = errors.New("RuleSet is nil")
		return
	}
	// determine whether there are specified fields in the input
	if !rs.isFit(input.Data) {
		return
	}

	redisClient, err := redisInstance()
	if err != nil {
		return
	}

	redisKey := rs.prepareRedisKey(produceName, input.Data)

	if input.Ts == 0 {
		input.Ts = TimeStamp()
	}
	scoreNumber := redis.Z{
		Score:  float64(input.Ts),
		Member: input.Id,
	}

	err = rs.ZAdd(redisClient, redisKey, scoreNumber)
	fmt.Println(redisKey, scoreNumber)
	if err != nil {
		return
	}
	FastRecoverGoroutineFunc(func() {
		// set the expiration time of the key
		rs.Expire(redisClient, redisKey)
		if input.Ts % 100 < 20 {
			errT := rs.ZRemByScore(redisKey, input.Ts)
			if errT != nil {
				log.Print(errT)
			}
		}
	})
	return
}

func (rs *RuleSet) ZAdd(rdsClient *redis.Client, redisKey string, member redis.Z) (err error) {
	_, err = rdsClient.ZAdd(redisKey, member).Result()
	if err != nil {
		err = fmt.Errorf("ZAdd failed[err]%v,[key]%v", err, redisKey)
	}
	return
}

func (rs *RuleSet) Expire(rdsClient *redis.Client, redisKey string) {
	_, err := rdsClient.PExpire(redisKey, time.Duration(rs.maxPeriod)*time.Millisecond).Result()
	if err != nil {
		log.Printf("Expire failed,[key]%v [expiration]%v", redisKey, rs.maxPeriod)
	}
	return
}

func (rs *RuleSet) ReadRedis(productName string, inputData map[string]interface{}, upperBound, lowBound int64) (fcDataList FcDataList, err error) {
	if rs == nil {
		err = errors.New("RuleSet is nil")
		return
	}
	if !rs.isFit(inputData) {
		return
	}

	redisClient, err := redisInstance()
	if err != nil {
		return
	}

	redisKey := rs.prepareRedisKey(productName, inputData)
	upperBoundStr := strconv.FormatInt(upperBound, 10)
	lowerBoundStr := strconv.FormatInt(lowBound, 10)
	zRangeBy := redis.ZRangeBy{Min: lowerBoundStr, Max: upperBoundStr, Count: rs.maxThreshold}
	readRedisResult, err := redisClient.ZRevRangeByScoreWithScores(redisKey, zRangeBy).Result()
	fcDataList = redisToFcData(readRedisResult)
	return
}

func (rs *RuleSet) ZRemByScore(key string, ts int64) (err error) {
	if rs == nil {
		err = errors.New("RuleSet is nil")
		return
	}
	upperBound := getLowerBound(ts, rs.maxPeriod)
	redisClient, err := redisInstance()
	if err != nil {
		return
	}
	return redisClient.ZRangeByScore(key, redis.ZRangeBy{Min:"-1", Max:fmt.Sprint(upperBound)}).Err()
}
