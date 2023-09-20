package dashboard

import (
	"github.com/gin-gonic/gin"

	"github.com/rusinikita/devex/config"
	"github.com/rusinikita/devex/datacollector"
	"github.com/rusinikita/devex/db"
)

func RunServer() datacollector.ErrStruct {
	database := db.DB()
	engine := gin.New()

	config.SetUpMiddlewares(engine, database)
	config.SetupRoutes(engine)

	err := engine.Run(":1080")

	return datacollector.ErrStruct{
		Template: "server %s \n",
		Error:    err,
	}
}
