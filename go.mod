module github.com/mhchlib/mconfig-go-sdk

go 1.14

require (
	github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc
	github.com/micro/go-micro/v2 v2.9.1
	github.com/tidwall/gjson v1.6.3
)

replace github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc => ../mconfig-api
