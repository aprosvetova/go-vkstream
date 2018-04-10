//Package vkstream provides methods for working with https://vk.com/dev/streaming_api
package vkstream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"net/url"
)

//NewStream creates new instance of VkStream. Use this if you already know the endpoint and have the key
func NewStream(endpoint string, key string) *VkStream {
	return &VkStream{
		Endpoint: endpoint,
		Key:      key,
	}
}

//NewStreamWithToken creates new instance of VkStream. Use this if you don't have the endpoint and the key.
//
//It requires an accessToken. It's named "Service token" and could be obtained from the app settings page (https://vk.com/apps?act=manage)
func NewStreamWithToken(accessToken string) (*VkStream, error) {
	resp, err := http.Get("https://api.vk.com/method/streaming.getServerUrl?access_token=" + accessToken + "&v=5.74")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var v vkAuthResponse
	if err := json.Unmarshal(bodyBuf, &v); err != nil {
		return nil, err
	}
	if v.Error.ErrorCode != 0 {
		return nil, newVkError(v.Error)
	}
	if v.Response.Key == "" || v.Response.Endpoint == "" {
		return nil, errors.New("error getting credentials")
	}
	v.Response.done = nil
	return &v.Response, nil
}

//GetRules gets a list of enabled rules
//
//It will return an error if something goes wrong
func (stream *VkStream) GetRules() ([]Rule, error) {
	resp, err := http.Get(fmt.Sprintf("https://%s/rules/?key=%s", stream.Endpoint, stream.Key))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var v vkRulesResponse
	if err := json.Unmarshal(bodyBuf, &v); err != nil {
		return nil, err
	}
	if v.Code != 200 {
		return nil, newVkStreamingError(v.Error)
	}
	return v.Rules, nil
}

//AddRule adds a new rule to the rules list
//
//Read more about rules syntax at https://vk.com/dev/streaming_api_docs?f=2.%20%D0%A4%D0%BE%D1%80%D0%BC%D0%B0%D1%82%20%D0%BF%D1%80%D0%B0%D0%B2%D0%B8%D0%BB
//
//Note that the tag must be unique, so it will return an error if there is a rule with the same tag already in the list
func (stream *VkStream) AddRule(value string, tag string) error {
	request := &vkAddRuleRequest{
		Rule: Rule{
			Value: value,
			Tag:   tag,
		},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("https://%s/rules/?key=%s", stream.Endpoint, stream.Key), "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var v vkEmptyResponse
	if err := json.Unmarshal(bodyBuf, &v); err != nil {
		return err
	}
	if v.Code != 200 {
		return newVkStreamingError(v.Error)
	}
	return nil
}

//DeleteRule deletes a rule from the rules list
//
//Note that it will return an error if there is no rule with the specified tag
func (stream *VkStream) DeleteRule(tag string) error {
	request := &vkDeleteRuleRequest{
		Tag: tag,
	}
	body, err := json.Marshal(request)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://%s/rules/?key=%s", stream.Endpoint, stream.Key), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var v vkEmptyResponse
	if err := json.Unmarshal(bodyBuf, &v); err != nil {
		return err
	}
	if v.Code != 200 {
		return newVkStreamingError(v.Error)
	}
	return nil
}

//DeleteAllRules clears the rules list
//
//Note that there is no such method in VK Streaming API, so it will get all rules first and delete them one by one then. It might take time if you have many rules
func (stream *VkStream) DeleteAllRules() error {
	rules, err := stream.GetRules()
	if err != nil {
		return err
	}
	for _, rule := range rules {
		err = stream.DeleteRule(rule.Tag)
		if err != nil {
			return err
		}
	}
	return nil
}

//Listen connects to VK Streaming API endpoint and listens for new events.
//
//It returns the channel or the error if something goes wrong. Subscribe to the channel to receive events
//
//Use Stop method if you want to stop listening
func (stream *VkStream) Listen() (<-chan Event, error) {
	u := url.URL{Scheme: "wss", Host: stream.Endpoint, Path: "/stream/", RawQuery: "key=" + stream.Key}
	c, wsResp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		if err == websocket.ErrBadHandshake {
			return nil, fmt.Errorf("handshake failed with status %d", wsResp.StatusCode)
		}
		return nil, err
	}
	stream.done = make(chan bool)
	events := make(chan Event)
	go func() {
		defer func() {
			c.Close()
			close(stream.done)
			close(events)
			stream.done = nil
		}()
		for {
			select {
			case <-stream.done:
				return
			default:
			}
			_, message, err := c.ReadMessage()
			if err != nil {
				return
			}
			var v vkStreamingMessage
			err = json.Unmarshal(message, &v)
			if err == nil {
				if v.Code == 100 {
					events <- v.Event
				}
			}
		}
	}()
	return events, nil
}

//Stop just closes the connection and stops listening for events
func (stream *VkStream) Stop() {
	stream.done <- true
}
