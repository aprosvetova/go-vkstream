package vkstream

//VkStream is the main type that allows you to interact with VK Streaming API
type VkStream struct {
	Endpoint string
	Key      string
	done     chan bool
}

//Rule is a set of keywords for filtering. Read more at https://vk.com/dev/streaming_api_docs?f=2.%20%D0%A4%D0%BE%D1%80%D0%BC%D0%B0%D1%82%20%D0%BF%D1%80%D0%B0%D0%B2%D0%B8%D0%BB
type Rule struct {
	Value string `json:"value"`
	Tag   string `json:"tag"`
}

//Event is a type that represents new action. Call Listen method to get events
type Event struct {
	Type string `json:"event_type"`
	Id   struct {
		PostOwnerId  int `json:"post_owner_id"`
		PostId       int `json:"post_id"`
		CommentId    int `json:"comment_id"`
		SharedPostId int `json:"shared_post_id"`
	} `json:"event_id"`
	Url                    string   `json:"event_url"`
	Text                   string   `json:"text"`
	Action                 string   `json:"action"`
	ActionTime             int      `json:"action_time"`
	CreationTime           int      `json:"creation_time"`
	SharedPostText         string   `json:"shared_post_text"`
	SharedPostCreationTime int      `json:"shared_post_creation_time"`
	SignerId               int      `json:"signer_id"`
	Tags                   []string `json:"tags"`
	Author                 struct {
		Id                  int    `json:"id"`
		Url                 string `json:"author_url"`
		SharedPostAuthorId  int    `json:"shared_post_author_id"`
		SharedPostAuthorUrl string `json:"shared_post_author_url"`
		Platform            int    `json:"platform"`
	} `json:"author"`
}

type vkStreamingMessage struct {
	Code           int `json:"code"`
	ServiceMessage struct {
		ServiceCode int    `json:"service_code"`
		Message     string `json:"message"`
	} `json:"service_message"`
	Event Event `json:"event"`
}

type vkStreamingError struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"error_code"`
}

type vkError struct {
	Message   string `json:"error_msg"`
	ErrorCode int    `json:"error_code"`
}

type vkAuthResponse struct {
	Response VkStream
	Error    vkError `json:"error"`
}

type vkRulesResponse struct {
	Code  int              `json:"code"`
	Rules []Rule           `json:"rules"`
	Error vkStreamingError `json:"error"`
}

type vkEmptyResponse struct {
	Code  int              `json:"code"`
	Error vkStreamingError `json:"error"`
}

type vkAddRuleRequest struct {
	Rule Rule `json:"rule"`
}

type vkDeleteRuleRequest struct {
	Tag string `json:"tag"`
}
