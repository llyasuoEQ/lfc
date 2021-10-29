# lfc 
lfc is short for Local Frequency Control, it is a local library package developed by golang.

lfc based on redis zset data structure, you can use it for frequency control in local, enjoy it!

# Installation
## Install
````
go get github.com/llyasuoEQ/lfc
````

## Import
````
github.com/llyasuoEQ/lfc
````

# Quickstart

## Write configuration
````
configToml := `[frequency]
         app_name = "test"
         # 规则
         # name: 规则名称，period：规定的时间，也就是滑动窗口的值，threshold：限制次数的阈值
         # code：规则的返回值，fields：规则字段
         rules=[
            {name="20秒限制3次",period=20,threshold=4,code=1001,fields=["phone","ip"]},
         ]

         # 频控要依赖的redis的配置
         [frequency.redis]
         address = "localhost:6379"
         password = ""
         db = 0
         poolsize = 50`
````