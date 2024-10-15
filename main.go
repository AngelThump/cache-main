package main

import (
	"fmt"
	"log"
	"time"

	b64 "encoding/base64"
	"encoding/json"

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
			user := api.GetUser(stream.User.Id)
			if user == nil {
				return
			}
			stream.User = *user

			base64String := b64.StdEncoding.EncodeToString([]byte(stream.Created_at + stream.User.Username))

			err := client.Rdb.Set(client.Ctx, stream.User.Username, base64String, 10*time.Second).Err()
			if err != nil {
				log.Println(err)
			}

			//set it for stream_key as well for mediamtx warmer
			err = client.Rdb.Set(client.Ctx, stream.User.StreamKey, base64String, 10*time.Second).Err()
			if err != nil {
				log.Println(err)
			}

			marshalledStream, err := json.Marshal(stream)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = client.Rdb.Set(client.Ctx, base64String+"_"+stream.User.Username, string(marshalledStream), 10*time.Second).Err()
			if err != nil {
				log.Println(err)
			}

		}(stream)
	}

	time.AfterFunc(1*time.Second, func() {
		saveStreams()
	})
}
