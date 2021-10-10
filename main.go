package main

import (
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/slack-go/slack"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var log = *logrus.New()

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load env", err)
	}

	channelId := os.Getenv("SLACK_CHANNEL_ID")
	authToken := os.Getenv("SLACK_AUTH_TOKEN")

	log.Info("ChannelId:", channelId)
	log.Info("authToken:", authToken)

	client := slack.New(authToken, slack.OptionDebug(true))

	attachment := slack.Attachment{
		Pretext: "Sample bot message",
		Text:    "Some text",
		Color:   "#36a64f",
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}
	_, timestamp, err := client.PostMessage(channelId, slack.MsgOptionAttachments(attachment))
	if err != nil {
		errors.Wrap(err, "Could not post the message on channel")
	}
	log.Info("Sent message", timestamp)

	channel, err := client.CreateConversation("test-channel", false)
	if err != nil {
		errors.Wrap(err, "Could not create conversation")
		return
	}
	attachment = slack.Attachment{
		Pretext: "Sample bot message",
		Text:    channel.LastRead,
		Color:   "green",
		Fields: []slack.AttachmentField{
			{
				Title: "Date",
				Value: time.Now().String(),
			},
		},
	}
	_, timestamp, err = client.PostMessage(channelId, slack.MsgOptionAttachments(attachment))
	if err != nil {
		errors.Wrap(err, "Could not post the message on channel")
	}
	log.Info("Sent message", timestamp)
}
