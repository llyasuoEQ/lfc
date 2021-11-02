package lfc

import (
	"fmt"
	"testing"
	"time"

	"github.com/bmizerany/assert"
)

const configToml = `
[frequency]
 app_name = "test"
 rules=[
		{name="「20s」limit【3】",period="20s",threshold=4,code=1000,fields=["uid","ip"]},
 ]
 # redis config
 [frequency.redis]
 address = "localhost:6379"
 password = "123456"
 db = 0
 poolsize = 50
 `

func initFrequency() (*ProductRule, error) {
	productRule, err := InitByConfigStr(configToml)
	if err != nil {
		return nil, err
	}
	return productRule, nil
}

func TestFrequency(t *testing.T) {
	productRule, err := initFrequency()
	if err != nil {
		t.Fatal(err)
	}
	inputData := make(map[string]interface{})
	inputData["uid"] = "1234"
	inputData["ip"] = "127.0.0.1"

	for i := 0; i < 4; i++ {
		input := NewInput(inputData, SetTimeStampOption(TimeStamp()), SetIdOption(fmt.Sprint(TimeStamp())))
		actual, err := productRule.FrequencyControl(input)
		if err != nil {
			t.Fatal(err)
		}
		var expect int64
		if i > 3 {
			expect = 1000
		}
		assert.Equal(t, expect, actual.Code, "frequency control test failed!")
		time.Sleep(1 * time.Microsecond)
	}

}
