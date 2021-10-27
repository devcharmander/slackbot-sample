package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/devcharmander/slack-bot/NLP"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api = slack.New(os.Getenv("SLACK_API"))
var signingSecret = os.Getenv("SLACK_SIGNING_SECRET")
var log = *logrus.New()

// You can open a dialog with a user interaction. (like pushing buttons, slash commands ...)
// https://api.slack.com/surfaces/modals
// https://api.slack.com/interactivity/entry-points
func main() {
	http.HandleFunc("/events", handler)
	http.ListenAndServe(":3000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// if err := sv.Ensure(); err != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(eventsAPIEvent.Type)
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		eventType := fmt.Sprintf("%T", eventsAPIEvent.InnerEvent.Data)
		fmt.Println(eventType)
		innerEvent := eventsAPIEvent.InnerEvent
		//fmt.Println(innerEvent.Data)
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			api.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
			client := NLP.NewClient()
			client.Analyze(ev.Text)
		case *slackevents.MessageEvent:
			log.Info("botID", ev.BotID, "message", ev.Text)
			if ev.BotID == "" {
				api.PostMessage(ev.Channel, slack.MsgOptionText(ev.Text, false))
			}
		}
	}
}
