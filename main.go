package main

import (
	"log"
	"time"

	b64 "encoding/base64"

	api "github.com/angelthump/cache-main/api"
	client "github.com/angelthump/cache-main/client"
	server "github.com/angelthump/cache-main/server"
	utils "github.com/angelthump/cache-main/utils"
)

func main() {
	cfgPath, err := utils.ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	err = utils.NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	client.Initalize()
	saveStreams()
	server.Initalize()
}

func saveStreams() {
	streams := api.Find()
	if streams == nil {
		time.AfterFunc(1*time.Second, func() {
			saveStreams()
		})
		return
	}

	for _, stream := range streams {
		go func(stream api.Stream) {
			base64String := b64.StdEncoding.EncodeToString([]byte(stream.Created_at + stream.User.Username))
			err := client.Rdb.Set(client.Ctx, stream.User.Username, base64String, 10*time.Second).Err()
			if err != nil {
				log.Println(err)
			}
		}(stream)
	}

	time.AfterFunc(1*time.Second, func() {
		saveStreams()
	})
}
