package lfc

import (
	redis "gopkg.in/redis.v5"
)

type FcData struct {
	Score  int64       `json:"ts"`
	Detail interface{} `json:"detail"`
}

type FcDataList []FcData

func (p FcDataList) Len() int           { return len(p) }
func (p FcDataList) Less(i, j int) bool { return p[i].Score < p[j].Score }
func (p FcDataList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (fcDataList *FcDataList) IsReachThreshold(threshold int64) (result bool) {
	if int64(len(*fcDataList)) >= threshold {
		result = true
	}
	return
}

// rankByScoreDesc...
// func rankByScoreDesc(data map[int64]interface{}) FcDataList {
// 	result := make(FcDataList, len(data))
// 	i := 0
// 	for k, v := range data {
// 		result[i] = FcData{k, v}
// 		i++
// 	}
// 	sort.Sort(sort.Reverse(result))
// 	return result
// }

// redisToFcData...  result read in redis is converted to FcDataList type
func redisToFcData(zs []redis.Z) FcDataList {
	result := make(FcDataList, len(zs))
	for k, v := range zs {
		result[k].Score = int64(v.Score)
		result[k].Detail = v.Member
	}
	return result
}
