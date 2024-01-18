package main

import (
	aw "github.com/deanishe/awgo"
	"github.com/slack-go/slack"
)

func updateChannels() {
	wf.NewItem("Update Channels").Valid(true)

	c := aw.NewCache(cache_dir)
	cfg := aw.NewConfig()
	token := cfg.Get("SLACK_TOKEN")
	api := slack.New(token)
	params := slack.GetConversationsForUserParameters{
    UserID: cfg.Get("USER_ID"),
}
	channels, _, err_channels := api.GetConversationsForUser(&params)
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
