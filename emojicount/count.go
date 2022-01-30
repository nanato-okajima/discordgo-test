package count

import (
	"fmt"

	dg "github.com/bwmarrin/discordgo"
)

type EmojiCount struct {
	Emoji *dg.Emoji
	Count int64
}

func CountEmoji(s *dg.Session, chanID string, lim int, before string, after string, around string) (*map[string]EmojiCount, error) {
	fmt.Println(before)
	msgs, err := s.ChannelMessages(chanID, lim, before, after, around)
	if err != nil {
		return nil, err
	}

	emjs := make(map[string]EmojiCount)
	msgID := ""
	for _, msg := range msgs {
		for _, v := range msg.Reactions {
			if v.Emoji == nil {
				continue
			}

			if val, ok := emjs[v.Emoji.Name]; ok {
				val.Count += int64(v.Count)
				emjs[v.Emoji.Name] = val
			} else {
				emjs[v.Emoji.Name] = EmojiCount{
					Emoji: v.Emoji,
					Count: int64(v.Count),
				}
			}
		}
		msgID = msg.ID
	}
	if msgID != before {
		CountEmoji(s, chanID, lim, msgID, after, around)
	}

	return &emjs, nil
}
