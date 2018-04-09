package main

import (
	"github.com/aprosvetova/go-vkstream"
	"log"
	"time"
)

func main() {
	stream, err := vkstream.NewStreamWithToken("your_token")
	if err != nil {
		log.Fatal(err)
		return
	}
	stream.AddRule("golang -rust", "1") //will listen for all posts/comments containing "golang" and NOT containing "rust", hehe. Read VK docs for detailed rules syntax
	events, err := stream.Listen()
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Print("Started listening")
	go func() {
		for event := range events {
			log.Print(event.Url) //will just print the url of the post/comment/etc
		}
	}()
	time.AfterFunc(time.Minute, func() {
		stream.Stop()
		log.Print("Stopped listening")
	}) //stop listening after a minute just to demonstrate how Stop works
	for {
		select {}
	}
}
