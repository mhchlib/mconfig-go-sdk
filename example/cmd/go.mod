module github.com/mhchlib/mconfig-go-sdk/exmample/cmd

go 1.14

require (
	github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc
	github.com/mhchlib/mconfig-go-sdk v0.0.0-20201023024357-d1f3591f172b
)

replace github.com/mhchlib/logger v0.0.0-20201023050446-420de20374cc => ../../../logger

replace github.com/mhchlib/register v0.0.0-20201119163729-b999bdbd2d49 => ../../../register

replace github.com/mhchlib/mconfig-api v0.0.0-20201023050446-420de20374cc => ../../../mconfig-api

replace github.com/mhchlib/mconfig-go-sdk v0.0.0-20201023024357-d1f3591f172b => ../../../mconfig-go-sdk
