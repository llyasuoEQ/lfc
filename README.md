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
         app_name = "test" # 
         # name: rule name，
         # period：frequency control time
         # threshold：how many times does it take effect
         # code：hit code
         # fields：rule field
         rules=[
            {name="『20s』limit【3】",period="20s",threshold=4,code=1001,fields=["phone","ip"]},
         ]

         # frequency control depends on the configuration of Redis
         [frequency.redis]
         address = "localhost:6379"
         password = ""
         db = 0
         poolsize = 50`
````