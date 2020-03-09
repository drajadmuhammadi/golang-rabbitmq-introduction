package models

type OpenRequestAPI struct {
	TopicName  string `json:"topic_name"`
	RoutingKey string `json:"routing_key"`
	Timeout    int    `json:"timeout"`
	Type       string `json:"type"`
	Data       interface{}
}

func OpenRequestAPIDefault() OpenRequestAPI {
	var defaultRequest OpenRequestAPI

	defaultRequest.Type = "async"
	defaultRequest.Timeout = 3

	return defaultRequest
}
