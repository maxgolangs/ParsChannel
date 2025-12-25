package parser

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gotd/td/tg"

	"pars/models"
)

func extractUserID(participant tg.ChannelParticipantClass) (int64, bool) {
	switch pt := participant.(type) {
	case *tg.ChannelParticipant:
		return pt.UserID, true
	case *tg.ChannelParticipantSelf:
		return pt.UserID, true
	case *tg.ChannelParticipantCreator:
		return pt.UserID, true
	case *tg.ChannelParticipantAdmin:
		return pt.UserID, true
	case *tg.ChannelParticipantBanned:
		if peer, ok := pt.Peer.(*tg.PeerUser); ok {
			return peer.UserID, true
		}
	}
	return 0, false
}

func buildUserMap(users []tg.UserClass) map[int64]*tg.User {
	userMap := make(map[int64]*tg.User)
	for _, u := range users {
		if user, ok := u.(*tg.User); ok {
			userMap[user.ID] = user
		}
	}
	return userMap
}

func processParticipant(
	participant tg.ChannelParticipantClass,
	userMap map[int64]*tg.User,
	seen map[int64]struct{},
	fw *FileWriters,
	result *models.ParseResult,
	userChan chan<- models.ParticipantInfo,
	ctx context.Context,
) error {
	userID, ok := extractUserID(participant)
	if !ok {
		return nil
	}

	if _, exists := seen[userID]; exists {
		return nil
	}
	seen[userID] = struct{}{}
	result.TotalUsers++

	user := userMap[userID]
	if user == nil {
		return nil
	}

	fw.IDsWriter.WriteString(fmt.Sprintf("%d\n", user.ID))

	if user.Username != "" {
		fw.UsernamesWriter.WriteString(fmt.Sprintf("@%s\n", user.Username))
		result.WithUsername++
	}

	info := models.ParticipantInfo{
		ID:           user.ID,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     user.Username,
		IsBot:        user.Bot,
		IsRestricted: user.Restricted,
		IsPremium:    user.Premium,
	}

	select {
	case userChan <- info:
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	fw.CSVWriter.Write([]string{
		strconv.FormatInt(info.ID, 10),
		info.FirstName,
		info.LastName,
		info.Username,
		strconv.FormatBool(info.IsPremium),
		strconv.FormatBool(info.IsBot),
		strconv.FormatBool(info.IsRestricted),
	})

	return nil
}

func fetchParticipants(
	ctx context.Context,
	api *tg.Client,
	channelInput tg.InputChannelClass,
	pattern string,
	offset int,
) (*tg.ChannelsChannelParticipants, error) {
	respRaw, err := api.ChannelsGetParticipants(ctx, &tg.ChannelsGetParticipantsRequest{
		Channel: channelInput,
		Filter:  &tg.ChannelParticipantsSearch{Q: pattern},
		Offset:  offset,
		Limit:   200,
		Hash:    0,
	})

	if err != nil {
		if strings.Contains(err.Error(), "FLOOD_WAIT") {
			time.Sleep(30 * time.Second)
			return nil, nil // Сигнал для повторной попытки
		}
		return nil, fmt.Errorf("get participants: %w", err)
	}

	resp, ok := respRaw.(*tg.ChannelsChannelParticipants)
	if !ok {
		return nil, fmt.Errorf("unexpected response type %T", respRaw)
	}

	return resp, nil
}

