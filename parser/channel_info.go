package parser

import (
	"context"
	"fmt"

	"github.com/gotd/td/tg"
)

type ChannelInfo struct {
	Name           string
	TotalMembers   int
	ChannelID      int64
}

func (p *RealtimeParser) GetChannelInfo(ctx context.Context, api *tg.Client, channelID int64, channelInput tg.InputChannelClass) (*ChannelInfo, error) {
	var channelName string
	var totalMembers int
	
	fullResp, err := api.ChannelsGetFullChannel(ctx, channelInput)
	if err != nil {
		return nil, fmt.Errorf("get channel info: %w", err)
	}
	
	for _, ch := range fullResp.Chats {
		channel, ok := ch.(*tg.Channel)
		if ok {
			channelName = channel.Title
			break
		}
	}
	
	if fullResp.FullChat != nil {
		if fullChannel, ok := fullResp.FullChat.(*tg.ChannelFull); ok {
			totalMembers = fullChannel.ParticipantsCount
		}
	}
	
	if channelName == "" {
		resp, err := api.ChannelsGetChannels(ctx, []tg.InputChannelClass{channelInput})
		if err == nil {
			var chats []tg.ChatClass
			switch r := resp.(type) {
			case *tg.MessagesChats:
				chats = r.Chats
			case *tg.MessagesChatsSlice:
				chats = r.Chats
			}
			
			for _, ch := range chats {
				channel, ok := ch.(*tg.Channel)
				if ok {
					channelName = channel.Title
					break
				}
			}
		}
	}
	
	if channelName == "" {
		channelName = fmt.Sprintf("Канал %d", channelID)
	}
	
	return &ChannelInfo{
		Name:         channelName,
		TotalMembers: totalMembers,
		ChannelID:    channelID,
	}, nil
}

