package parser

import (
	"context"
	"fmt"
	"strings"

	"github.com/gotd/td/tg"
)

func buildPatterns() []string {
	ru := "абвгдеёжзийклмнопрстуфхцчшщъыьэюя"
	var patterns []string
	patterns = append(patterns, "")

	for _, ch := range ru {
		patterns = append(patterns, string(ch))
		patterns = append(patterns, strings.ToUpper(string(ch)))
	}

	for _, ch := range "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" {
		patterns = append(patterns, string(ch))
	}

	return patterns
}

func channelIDFromChatID(chatID int64) int64 {
	abs := chatID
	if abs < 0 {
		abs = -abs
	}
	return abs - 1000000000000
}

func (p *RealtimeParser) findChannelInput(ctx context.Context, api *tg.Client, chatID int64) (*tg.InputChannel, error) {
	targetID := channelIDFromChatID(chatID)

	respRaw, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{},
		OffsetDate: 0,
		OffsetID:   0,
		Limit:      200,
		Hash:       0,
	})
	if err == nil {
		var chats []tg.ChatClass
		switch resp := respRaw.(type) {
		case *tg.MessagesDialogs:
			chats = resp.Chats
		case *tg.MessagesDialogsSlice:
			chats = resp.Chats
		case *tg.MessagesDialogsNotModified:
		default:
		}

		for _, ch := range chats {
			channel, ok := ch.(*tg.Channel)
			if !ok {
				continue
			}
			if channel.ID == targetID {
				return &tg.InputChannel{
					ChannelID:  channel.ID,
					AccessHash: channel.AccessHash,
				}, nil
			}
		}
	}

	resp, err := api.ChannelsGetChannels(ctx, []tg.InputChannelClass{
		&tg.InputChannel{
			ChannelID:  targetID,
			AccessHash: 0,
		},
	})
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
			if ok && channel.ID == targetID {
				return &tg.InputChannel{
					ChannelID:  channel.ID,
					AccessHash: channel.AccessHash,
				}, nil
			}
		}
	}

	fullResp, err := api.ChannelsGetFullChannel(ctx, &tg.InputChannel{
		ChannelID:  targetID,
		AccessHash: 0,
	})
	if err == nil {
		for _, ch := range fullResp.Chats {
			channel, ok := ch.(*tg.Channel)
			if ok && channel.ID == targetID {
				return &tg.InputChannel{
					ChannelID:  channel.ID,
					AccessHash: channel.AccessHash,
				}, nil
			}
		}
	}

	return nil, fmt.Errorf("channel %d not found; ensure the bot is added to the channel", chatID)
}

