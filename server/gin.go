package server

import (
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/angelthump/cache-main/api"
	"github.com/angelthump/cache-main/client"
	utils "github.com/angelthump/cache-main/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func Initalize() {
	if utils.Config.GinReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})

	router.POST("/hls/:channel/:endUrl", func(c *gin.Context) {
		jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		if jwtToken != utils.Config.IngestAPI.AuthKey {
			c.AbortWithStatus(500)
			return
		}

		channel := c.Param("channel")
		regex := regexp.MustCompile(`_src|_medium|_low`)
		base64Channel := regex.ReplaceAllString(channel, "")
		endUrl := c.Param("endUrl")

		base64Path, err := client.Rdb.Get(client.Ctx, base64Channel).Result()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		key := base64Path + "/" + endUrl
		data, _ := c.GetRawData()
		if strings.HasSuffix(endUrl, ".ts") {
			client.Rdb.Set(client.Ctx, key, data, 30*time.Second)
		} else if strings.HasSuffix(endUrl, ".m3u8") {
			client.Rdb.Set(client.Ctx, key, data, 30*time.Second)
		} else if strings.HasSuffix(endUrl, "init.mp4") {
			client.Rdb.Set(client.Ctx, key, data, 30*time.Second)
		} else if strings.HasSuffix(endUrl, ".mp4") {
			client.Rdb.Set(client.Ctx, key, data, 30*time.Second)
		} else if strings.HasSuffix(endUrl, ".m4s") {
			client.Rdb.Set(client.Ctx, key, data, 30*time.Second)
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

		//c.AbortWithStatus(404)

		endPathRegex := regexp.MustCompile(`.ts|.mp4|.m4s`)
		if !endPathRegex.MatchString(endUrl) {
			c.AbortWithStatus(500)
		}

		regex := regexp.MustCompile(`_src|_medium|_low`)
		base64Channel := regex.ReplaceAllString(c.Param("channel"), "")

		data, err = client.Rdb.Get(client.Ctx, base64Channel).Result()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		var streamData api.Stream
		err = json.Unmarshal([]byte(data), &streamData)
		if err != nil {
			log.Printf("Unmarshal Error %v", err)
			c.AbortWithStatus(500)
			return
		}

		if !streamData.Ingest.Mediamtx {
			c.AbortWithStatus(404)
			return
		}

		fullEndUrl := strings.Replace(c.Param("endUrl"), "index", "stream", 1)
		query := c.Request.URL.Query().Encode()
		if len(query) > 0 {
			fullEndUrl += "?" + query
		}
		mediamtxEndUrl := "https://" + utils.Config.IngestAPI.Username + ":" + utils.Config.IngestAPI.Password + "@" + streamData.Ingest.Server + ".angelthump.com/hls/live/" + streamData.User.StreamKey + "/" + fullEndUrl

		restyClient := resty.New()
		resp, _ := restyClient.R().
			SetHeader("X-Api-Key", utils.Config.StreamsAPI.AuthKey).
			Get(mediamtxEndUrl)

		statusCode := resp.StatusCode()
		if statusCode >= 400 {
			c.AbortWithStatus(statusCode)
			return
		}

		c.Header("Access-Control-Allow-Origin", "*")

		if strings.HasSuffix(endUrl, ".ts") {
			client.Rdb.Set(client.Ctx, key, resp.Body(), 30*time.Second)
			c.Data(200, "video/mp2t", []byte(resp.Body()))
		} else if strings.HasSuffix(endUrl, ".mp4") {
			client.Rdb.Set(client.Ctx, key, resp.Body(), 30*time.Second)
			c.Data(200, "video/mp4", []byte(resp.Body()))
		} else if strings.HasSuffix(endUrl, ".m4s") {
			client.Rdb.Set(client.Ctx, key, resp.Body(), 30*time.Second)
			c.Data(200, "video/mp4", []byte(resp.Body()))
		} else {
			c.AbortWithStatus(400)
		}
	})

	router.Run(":" + utils.Config.Port)
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
