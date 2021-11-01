package lfc

type Input struct {
	Data map[string]interface{} // input data
	Ts   int64  // Millisecond timestamp
	Id   string // Identity unique
}

type InputOption func(input *Input)

func NewInput(data map[string]interface{}, options ...InputOption) *Input {
	if len(data) < 1 {
		return nil
	}
	input := &Input{Data: data}
	for _, f := range options {
		f(input)
	}
	return input
}

func SetTimeStampOption(timeStamp int64) InputOption {
	return func(i *Input) {
		i.Ts = timeStamp
	}
}

func SetIdOption(id string) InputOption {
	return func(i *Input) {
		i.Id = id
	}
}