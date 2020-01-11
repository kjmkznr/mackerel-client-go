package mackerel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestListChannels(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/channels" {
			t.Error("request URL should be /api/v0/channels but: ", req.URL.Path)
		}

		if req.Method != "GET" {
			t.Error("request method should be GET but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string][]map[string]interface{}{
			"channels": {
				{
					"id":      "abcdefabc",
					"name":    "email channel",
					"type":    "email",
					"emails":  []string{"test@example.com", "test2@example.com"},
					"userIds": []string{"1234", "2345"},
					"events":  []string{"alert"},
				},
				{
					"id":   "bcdefabcd",
					"name": "slack channel",
					"type": "slack",
					"url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
					"mentions": map[string]interface{}{
						"ok":      "ok message",
						"warning": "warning message",
					},
					"enabledGraphImage": true,
					"events":            []string{"alert"},
				},
				{
					"id":     "cdefabcde",
					"name":   "webhook channel",
					"type":   "webhook",
					"url":    "http://example.com/webhook",
					"events": []string{"alert"},
				},
				{
					"id":   "defabcdef",
					"name": "line channel",
					"type": "line",
				},
			},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)
	channels, err := client.ListChannels()

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if reflect.DeepEqual(channels[0].Emails, []string{"test@example.com", "test2@example.com"}) != true {
		t.Errorf("Wrong data for emails: %v", channels[0].Emails)
	}
	if reflect.DeepEqual(channels[0].UserIDs, []string{"1234", "2345"}) != true {
		t.Errorf("Wrong data for emails: %v", channels[0].UserIDs)
	}

	if channels[1].Mentions.OK != "ok message" {
		t.Error("request has mentions.ok but: ", channels[1].Mentions.OK)
	}
	if channels[1].Mentions.Warning != "warning message" {
		t.Error("request has mentions.warning but: ", channels[1].Mentions.Warning)
	}
	if channels[1].Mentions.Critical != "" {
		t.Error("request does not have mentions.critical but: ", channels[1].Mentions.Critical)
	}

	if channels[0].URL != "" {
		t.Error("request has no URL but: ", channels[0])
	}
	if channels[1].URL != "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX" {
		t.Error("request sends json including URL but: ", channels[1])
	}
	if channels[2].URL != "http://example.com/webhook" {
		t.Error("request sends json including URL but: ", channels[2])
	}
	if channels[3].URL != "" {
		t.Error("request has no URL but: ", channels[3])
	}
}

func TestCreateChannel(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/api/v0/channels" {
			t.Error("request URL should be /api/v0/channels but: ", req.URL.Path)
		}

		if req.Method != "POST" {
			t.Error("request method should be POST but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":   "abcdefabc",
			"name": "slack channel",
			"type": "slack",
			"url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
			"mentions": map[string]interface{}{
				"ok":      "ok message",
				"warning": "warning message",
			},
			"enabledGraphImage": true,
			"events":            []string{"alert"},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	channel, err := client.CreateChannel(&ChannelWithoutID{
		Name: "slack channel",
		Type: "slack",
		URL:  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
		Mentions: Mentions{
			OK:      "ok message",
			Warning: "warning message",
		},
		EnabledGraphImage: true,
		Events:            []string{"alert"},
	})

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if channel.ID != "abcdefabc" {
		t.Error("request sends json including ID but: ", channel.ID)
	}
	if channel.Name != "slack channel" {
		t.Error("request sends json including name but: ", channel.Name)
	}
	if channel.URL != "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX" {
		t.Error("request sends json including URL but: ", channel.URL)
	}
	if reflect.DeepEqual(channel.Mentions, Mentions{OK: "ok message", Warning: "warning message"}) != true {
		t.Errorf("Wrong data for mentions: %v", channel.Mentions)
	}
	if !channel.EnabledGraphImage {
		t.Error("request sends json including enabledGraphImage but: ", channel.EnabledGraphImage)
	}
	if reflect.DeepEqual(channel.Events, []string{"alert"}) != true {
		t.Errorf("Wrong data for events: %v", channel.Events)
	}
}

func TestDeleteChannel(t *testing.T) {
	channelID := "abcdefabc"
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if req.URL.Path != fmt.Sprintf("/api/v0/channels/%s", channelID) {
			t.Error("request URL should be /api/v0/channels/<ID> but: ", req.URL.Path)
		}

		if req.Method != "DELETE" {
			t.Error("request method should be DELETE but: ", req.Method)
		}

		respJSON, _ := json.Marshal(map[string]interface{}{
			"id":   channelID,
			"name": "slack channel",
			"type": "slack",
			"url":  "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX",
			"mentions": map[string]interface{}{
				"ok":      "ok message",
				"warning": "warning message",
			},
			"enabledGraphImage": true,
			"events":            []string{"alert"},
		})

		res.Header()["Content-Type"] = []string{"application/json"}
		fmt.Fprint(res, string(respJSON))
	}))
	defer ts.Close()

	client, _ := NewClientWithOptions("dummy-key", ts.URL, false)

	channel, err := client.DeleteChannel(channelID)

	if err != nil {
		t.Error("err should be nil but: ", err)
	}

	if channel.ID != "abcdefabc" {
		t.Error("request sends json including ID but: ", channel.ID)
	}
	if channel.Name != "slack channel" {
		t.Error("request sends json including name but: ", channel.Name)
	}
	if channel.URL != "https://hooks.slack.com/services/TAAAA/BBBB/XXXXX" {
		t.Error("request sends json including URL but: ", channel.URL)
	}
	if reflect.DeepEqual(channel.Mentions, Mentions{OK: "ok message", Warning: "warning message"}) != true {
		t.Errorf("Wrong data for mentions: %v", channel.Mentions)
	}
	if !channel.EnabledGraphImage {
		t.Error("request sends json including enabledGraphImage but: ", channel.EnabledGraphImage)
	}
	if reflect.DeepEqual(channel.Events, []string{"alert"}) != true {
		t.Errorf("Wrong data for events: %v", channel.Events)
	}
}
