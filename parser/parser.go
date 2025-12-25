package parser

import (
	"context"
	"path/filepath"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"pars/models"
)

type RealtimeParser struct {
	apiID   int
	apiHash string
}

func New(apiID int, apiHash string) *RealtimeParser {
	return &RealtimeParser{
		apiID:   apiID,
		apiHash: apiHash,
	}
}

func (p *RealtimeParser) ParseChannelRealtime(ctx context.Context, botToken string, channelID int64, outputDir string, userChan chan<- models.ParticipantInfo, channelInfoChan chan<- *models.ParseResult) (*models.ParseResult, error) {
	fw, result, err := initOutputFiles(outputDir)
	if err != nil {
		return nil, err
	}
	defer fw.Close()

	client := telegram.NewClient(p.apiID, p.apiHash, telegram.Options{
		SessionStorage: &telegram.FileSessionStorage{Path: filepath.Join(outputDir, "session.json")},
	})

	err = client.Run(ctx, func(ctx context.Context) error {
		if _, err := client.Auth().Bot(ctx, botToken); err != nil {
			return err
		}

		api := client.API()
		
		channelInput, err := p.findChannelInput(ctx, api, channelID)
		if err != nil {
			return err
		}
		
		info, err := p.GetChannelInfo(ctx, api, channelID, channelInput)
		if err == nil {
			result.ChannelName = info.Name
			result.TotalMembers = info.TotalMembers
			if channelInfoChan != nil {
				infoCopy := *result
				select {
				case channelInfoChan <- &infoCopy:
				case <-ctx.Done():
					return ctx.Err()
				}
			}
		}

		return p.parseParticipants(ctx, api, channelInput, fw, result, userChan)
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (p *RealtimeParser) parseParticipants(
	ctx context.Context,
	api *tg.Client,
	channelInput tg.InputChannelClass,
	fw *FileWriters,
	result *models.ParseResult,
	userChan chan<- models.ParticipantInfo,
) error {
	patterns := buildPatterns()
	seen := make(map[int64]struct{})

	for _, pattern := range patterns {
		if result.TotalMembers > 0 && result.TotalUsers >= result.TotalMembers {
			break
		}
		
		offset := 0
		
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if result.TotalMembers > 0 && result.TotalUsers >= result.TotalMembers {
				return nil
			}

			resp, err := fetchParticipants(ctx, api, channelInput, pattern, offset)
			if err != nil {
				return err
			}

			if resp == nil {
				continue // FLOOD_WAIT, повторная попытка
			}

			if len(resp.Participants) == 0 {
				break
			}

			userMap := buildUserMap(resp.Users)

			for _, participant := range resp.Participants {
				if err := processParticipant(participant, userMap, seen, fw, result, userChan, ctx); err != nil {
					return err
				}
				
				if result.TotalMembers > 0 && result.TotalUsers >= result.TotalMembers {
					return nil
				}
			}

			offset += len(resp.Participants)
		}
	}

	return nil
}
