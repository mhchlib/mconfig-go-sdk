package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mhchlib/mconfig-go-sdk/mconfig"
	"log"
)

func main() {
	config := mconfig.NewMconfig(
		mconfig.NameSpace("local_test"),
		mconfig.Registry(mconfig.RegisterType_Etcd, []string{"etcd.u.hcyang.top:31770"}),
		mconfig.Metadata("port", "8080"), //meta data
		mconfig.Metadata("ip", "192.0.0.1"),
		mconfig.AppKey("app_tPss5k5H_DS"),
		mconfig.ConfigKey("config_tPssCRQrGxh"),
	)
	r := gin.Default()
	r.GET("/mconfig-server/:type/name/:var", func(c *gin.Context) {
		pVar := c.Param("var")
		pType := c.Param("type")
		log.Println(pVar, pType)
		var res interface{}
		switch pType {
		case "string":
			res = config.String(pVar)
		case "list":
			res = config.List(pVar)
		case "map":
			res = config.Map(pVar)
		case "int":
			res = config.Int64(pVar)
		case "bool":
			res = config.Bool(pVar)

		}
		c.JSON(200, gin.H{
			"message": res,
		})
	})
	_ = r.Run(":8888") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
