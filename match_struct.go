package lfc

// frequency control return result
type MatchResult struct {
	Code int64
	Data interface{}
}

// frequency control rule
type MatchDetail struct {
	RuleName     string
	Period       int64
	Threshold    int64
	RecordDetail FcDataList
}
