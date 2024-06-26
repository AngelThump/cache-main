package server

import (
	"regexp"
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
		regex := regexp.MustCompile(`_src|_medium|_low`)
		base64Channel := regex.ReplaceAllString(channel, "")
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
			client.Rdb.Set(client.Ctx, key, data, 16*time.Second)
		} else if strings.HasSuffix(endUrl, "init.mp4") {
			client.Rdb.Set(client.Ctx, key, data, 24*time.Hour)
		} else if strings.HasSuffix(endUrl, ".mp4") {
			client.Rdb.Set(client.Ctx, key, data, 20*time.Second)
		} else if strings.HasSuffix(endUrl, ".m4s") {
			client.Rdb.Set(client.Ctx, key, data, 20*time.Second)
		} else {
			c.AbortWithStatus(400)
		}
	})

	router.POST("/hls/live/:channel/:endUrl", func(c *gin.Context) {
		channel := c.Param("channel")
		regex := regexp.MustCompile(`_src|_medium|_low`)
		base64Channel := regex.ReplaceAllString(channel, "")
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
			client.Rdb.Set(client.Ctx, key, data, 16*time.Second)
		} else if strings.HasSuffix(endUrl, "init.mp4") {
			client.Rdb.Set(client.Ctx, key, data, 24*time.Hour)
		} else if strings.HasSuffix(endUrl, ".mp4") {
			client.Rdb.Set(client.Ctx, key, data, 20*time.Second)
		} else if strings.HasSuffix(endUrl, ".m4s") {
			client.Rdb.Set(client.Ctx, key, data, 20*time.Second)
		} else {
			c.AbortWithStatus(400)
		}
	})

	router.GET("/hls/:channel/:endUrl", func(c *gin.Context) {
		channel := c.Param("channel")
		endUrl := c.Param("endUrl")

		key := channel + "/" + endUrl

		data, err := client.Rdb.Get(client.Ctx, key).Result()
		if err == nil {
			c.Header("Access-Control-Allow-Origin", "*")

			if strings.HasSuffix(endUrl, ".ts") {
				c.Data(200, "video/mp2t", []byte(data))
			} else if strings.HasSuffix(endUrl, ".m3u8") {
				c.Data(200, "application/x-mpegURL", []byte(data))
			} else if strings.HasSuffix(endUrl, ".mp4") {
				c.Data(200, "video/mp4", []byte(data))
			} else if strings.HasSuffix(endUrl, ".m4s") {
				c.Data(200, "video/mp4", []byte(data))
			} else {
				c.AbortWithStatus(400)
			}
			return
		}

		c.AbortWithStatus(404)

		/*
			regex := regexp.MustCompile(`_src|_medium|_low`)
			base64Channel := regex.ReplaceAllString(c.Param("channel"), "")

			data, err = client.Rdb.Get(client.Ctx, base64Channel).Result()
			if err != nil {
				c.AbortWithStatus(400)
				return
			}

			var streamData api.Stream
			err = json.Unmarshal([]byte(data), &streamData)
			if err != nil {
				log.Printf("Unmarshal Error %v", err)
				c.AbortWithStatus(400)
				return
			}

			if !streamData.Ingest.Mediamtx {
				c.AbortWithStatus(404)
				return
			}

			restyClient := resty.New()
			resp, _ := restyClient.R().
				SetHeader("X-Api-Key", utils.Config.StreamsAPI.AuthKey).
				Get("https://" + utils.Config.IngestAPI.Username + ":" + utils.Config.IngestAPI.Password + "@" + streamData.Ingest.Server + ".angelthump.com/hls/live/" + streamData.User.Username + "/" + endUrl)

			statusCode := resp.StatusCode()
			if statusCode >= 400 {
				c.AbortWithStatus(400)
				return
			}

			c.Header("Access-Control-Allow-Origin", "*")

			if strings.HasSuffix(endUrl, ".ts") {
				client.Rdb.Set(client.Ctx, key, resp.Body(), 20*time.Second)
				c.Data(200, "video/mp2t", []byte(resp.Body()))
			} else if strings.HasSuffix(endUrl, "init.mp4") {
				client.Rdb.Set(client.Ctx, key, data, 1*time.Second)
				c.Data(200, "video/mp4", []byte(resp.Body()))
			} else if strings.HasSuffix(endUrl, ".mp4") {
				client.Rdb.Set(client.Ctx, key, resp.Body(), 20*time.Second)
				c.Data(200, "video/mp4", []byte(resp.Body()))
			} else if strings.HasSuffix(endUrl, ".m4s") {
				client.Rdb.Set(client.Ctx, key, resp.Body(), 20*time.Second)
				c.Data(200, "video/mp4", []byte(resp.Body()))
			} else if strings.HasSuffix(endUrl, ".m3u8") {
				client.Rdb.Set(client.Ctx, key, resp.Body(), 1*time.Second)
				c.Data(200, "application/x-mpegURL", []byte(resp.Body()))
			} else {
				c.AbortWithStatus(400)
			}*/
	})

	router.Run(":" + utils.Config.Port)
}
