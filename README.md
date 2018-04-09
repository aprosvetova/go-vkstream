# Golang bindings for the VK Streaming API

[![GoDoc](https://godoc.org/github.com/aprosvetova/go-vkstream?status.svg)](https://godoc.org/github.com/aprosvetova/go-vkstream)

I'm very new to Go, so I'll be happy if you make some Pull Requests and help me with tests.

Read the VK Streaming API [description](https://vk.com/dev/streaming_api) and [docs](https://vk.com/dev/streaming_api_docs) to understand how it works.

**Warning! Please note that VK will send you only 1% of events that match your rules. You need to ask [the support](https://vk.com/support?act=new_api) for the full access to get all the events.**

## Example

```go
package main

import (
	"github.com/aprosvetova/go-vkstream"
	"log"
	"time"
	"go-vkstream"
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
```

## Contact me
If you have any questions about my shitty code, feel free to contact me by Telegram ([@koteeq](https://t.me/koteeq)). I speak English and Russian.