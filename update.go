package main

import (
	aw "github.com/deanishe/awgo"
	"github.com/slack-go/slack"
)

func updateChannels() {
	wf.NewItem("Update Channels").Valid(true)

	println("START update channels")
	c := aw.NewCache(cache_dir)
	cfg := aw.NewConfig()
	token := cfg.Get("SLACK_TOKEN")
	api := slack.New(token)
	params := slack.GetConversationsForUserParameters{
		UserID: cfg.Get("USER_ID"),
	}
	var channels []slack.Channel
	var next_cursor string
	var cnt int = 0
	var err_channels error

	for {
		var newChannels []slack.Channel

		if next_cursor != "" {
			params.Cursor = next_cursor
		}

		newChannels, next_cursor, err_channels = api.GetConversationsForUser(&params)
		println(next_cursor)

		channels = append(channels, newChannels...)

		cnt++
		println(cnt)
		if next_cursor == "" || cnt > 10 {
			break
		}
	}
	team, err_team := api.GetTeamInfo()

	if err_channels != nil || err_team != nil {
		wf.Warn("Error", "Error occurred in Slack API ")
	}

	all_channels := make([]Channel, 0)
	for _, channel := range channels {
		all_channels = append(all_channels, Channel{
			Name:   channel.Name,
			ID:     channel.ID,
			TeamID: team.ID,
		})
	}

	c.StoreJSON(cache_file, all_channels)
	wf.SendFeedback()
}
