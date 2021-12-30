package server

import (
	"strings"
	"time"

	"github.com/angelthump/cache-main/client"
	utils "github.com/angelthump/cache-main/utils"
	"github.com/gin-gonic/gin"
)

func Initalize() {
	if utils.Config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.POST("/hls/:channel/:endUrl", func(c *gin.Context) {
		channel := c.Param("channel")
		var base64Channel string
		variant := strings.Index(c.Param("channel"), "_")
		if variant == -1 {
			base64Channel = c.Param("channel")
		} else {
			base64Channel = c.Param("channel")[0:strings.Index(c.Param("channel"), "_")]
		}
		endUrl := c.Param("endUrl")

		base64String, err := client.Rdb.Get(client.Ctx, base64Channel).Result()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		key := base64String + "_" + channel + "/" + endUrl
		data, _ := c.GetRawData()
		if strings.HasSuffix(endUrl, ".ts") {
			client.Rdb.Set(client.Ctx, key, data, 20*time.Second)
		} else if strings.HasSuffix(endUrl, ".m3u8") {
			client.Rdb.Set(client.Ctx, key, data, 4*time.Second)
		} else {
			c.AbortWithStatus(400)
		}
	})

	router.GET("/hls/:channel/:endUrl", func(c *gin.Context) {
		channel := c.Param("channel")
		endUrl := c.Param("endUrl")

		key := channel + "/" + endUrl

		data, err := client.Rdb.Get(client.Ctx, key).Result()
		if err != nil {
			c.AbortWithStatus(404)
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")

		if strings.HasSuffix(endUrl, ".ts") {
			c.Data(200, "video/mp2t", []byte(data))
		} else if strings.HasSuffix(endUrl, ".m3u8") {
			c.Data(200, "application/x-mpegURL", []byte(data))
		} else {
			c.AbortWithStatus(400)
		}
	})

	router.Run(":" + utils.Config.Port)
}
