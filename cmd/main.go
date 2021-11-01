package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/llyasuoEQ/lfc"
	toml "github.com/pelletier/go-toml"
)

type Req struct {
	Uid string `json:"uid"`
}

func NewLfc() (*lfc.ProductRule, error) {
	pwd, _ := os.Getwd()
	path := pwd + "/cmd/config.toml"
	tr, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}
	return lfc.InitConfig(tr)
}

func main() {
	app := gin.New()
	r := app.Group("/")
	{
		r.POST("/frequency", frequency)
	}
	_ = app.Run(":8080")
}

func frequency(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		c.JSON(1001, fmt.Sprintf("c.GetRawData() error:%s", err.Error()))
		return
	}
	req := Req{}
	err = json.Unmarshal(data, &req)
	if err != nil {
		c.JSON(1001, fmt.Sprintf("json.Unmarshal(data, &req) error:%s", err.Error()))
		return
	}
	pr, err := NewLfc()
	if err != nil {
		c.JSON(1001, fmt.Sprintf("NewLfc() error:%s", err.Error()))
		return
	}
	input :=lfc.NewInput(map[string]interface{}{"uid": req.Uid}, lfc.SetTimeStampOption(lfc.TimeStamp()),
	lfc.SetIdOption(fmt.Sprint(lfc.TimeStamp())))
	res, err := pr.FrequencyControl(input)
	if err != nil {
		c.JSON(1001, fmt.Sprintf("pr.FrequencyControl(ctx, input) error:%s", err.Error()))
		return
	}
	c.JSON(http.StatusOK, res)
}
