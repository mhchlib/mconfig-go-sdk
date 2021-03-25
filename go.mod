module github.com/mhchlib/mconfig-go-sdk

go 1.14

require (
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d // indirect
	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/mregister v0.0.0-20201119163729-b999bdbd2d49
	github.com/prometheus/common v0.6.0
	github.com/spf13/viper v1.7.1
	github.com/tidwall/gjson v1.6.3
	google.golang.org/grpc v1.26.0
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../logger

replace github.com/mhchlib/mregister v0.0.0-20201119163729-b999bdbd2d49 => ../mregister

replace github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc => ../mconfig-api
